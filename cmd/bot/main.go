package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tian841224/stock-bot/internal/application/usecase/bot"
	healthUsecase "github.com/tian841224/stock-bot/internal/application/usecase/health"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	"github.com/tian841224/stock-bot/internal/application/usecase/user"
	formatterAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/formatter"
	healthAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/health"
	marketAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/market"
	presenterAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/presenter"
	userSubscriptionAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/user"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/line"
	tgbotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/imgbb"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/twse"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	database "github.com/tian841224/stock-bot/internal/infrastructure/persistence"
	repository "github.com/tian841224/stock-bot/internal/infrastructure/persistence/postgres"
	linebot "github.com/tian841224/stock-bot/internal/interfaces/bot/line"
	telegram "github.com/tian841224/stock-bot/internal/interfaces/bot/telegram"
	healthHandler "github.com/tian841224/stock-bot/internal/interfaces/health"
	"github.com/tian841224/stock-bot/internal/interfaces/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("載入配置失敗: %v", err)
	}

	// 初始化 Logger
	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("初始化 Logger 失敗: %v", err)
	}

	appLogger.Info("starting application...")
	appLogger.Info("config loaded")

	// 初始化資料庫
	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("failed to init database", logger.Error(err))
	}
	defer db.Close()

	gormDB := db.GetDB()

	// ============================================================
	// 建立 Repository 層（Persistence）
	// ============================================================

	appLogger.Info("initializing repository layer...")
	userRepo := repository.NewPostgresUserRepository(gormDB, appLogger)
	stockSymbolRepo := repository.NewSymbolRepository(gormDB, appLogger)
	tradeDateRepo := repository.NewPostgresTradeDateRepository(gormDB, appLogger)
	subscriptionRepo := repository.NewSubscriptionRepository(gormDB, appLogger)
	subscriptionSymbolRepo := repository.NewSubscriptionSymbolRepository(gormDB, appLogger)
	featureReader, _ := repository.NewFeatureRepository(gormDB, appLogger)
	syncMetadataRepo := repository.NewSyncMetadataRepository(gormDB, appLogger)
	appLogger.Info("feature repository initialized")

	// ============================================================
	// 建立外部服務客戶端（External Services）
	// ============================================================
	appLogger.Info("initializing external service clients...")

	// Telegram Bot client
	appLogger.Info("initializing Telegram Bot client...")
	tgClient, err := tgbotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create Telegram Bot client", logger.Error(err))
	}
	appLogger.Info("Telegram Bot client initialized")

	// LINE Bot client
	appLogger.Info("initializing LINE Bot client...")
	lineClient, err := linebotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create LINE Bot client", logger.Error(err))
	}
	appLogger.Info("LINE Bot client initialized")

	// image upload service
	appLogger.Info("initializing external service clients...")
	imgbbClient := imgbb.NewImgBBClient(cfg.IMGBB_API_KEY)

	// 股票 API 客戶端
	fugleAPI := fugle.NewFugleAPI(*cfg, appLogger)
	twseAPI := twse.NewTwseAPI(appLogger)
	cnyesAPI := cnyes.NewCnyesAPI(appLogger)
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg, appLogger)
	appLogger.Info("external service clients initialized")

	// ============================================================
	// 建立 Adapter 層（Gateway/Presenter）
	// ============================================================
	appLogger.Info("initializing adapter layer...")
	// Validation Gateway
	validationGateway := presenterAdapter.NewValidationGateway(nil, stockSymbolRepo)

	// User Subscription Gateway
	userSubscriptionGateway := userSubscriptionAdapter.NewUserSubscriptionGateway(
		subscriptionRepo,
		subscriptionSymbolRepo,
		stockSymbolRepo,
		subscriptionRepo,
		subscriptionSymbolRepo,
		featureReader,
		userRepo,
	)

	// Market Data Gateway
	marketDataGateway := marketAdapter.NewMarketDataGateway(
		twseAPI,
		cnyesAPI,
		fugleAPI,
		finmindAPI,
		validationGateway,
		tradeDateRepo,
	)

	// Market Chart Gateway
	marketChartGateway := marketAdapter.NewMarketChartGateway(
		marketDataGateway,
		validationGateway,
		fugleAPI,
	)

	// Formatter Adapter
	telegramFormatter := formatterAdapter.NewTelegramFormatter()
	lineFormatter := formatterAdapter.NewLineFormatter()
	formatterGateway := formatterAdapter.NewFormatterAdapter(
		marketChartGateway,
		validationGateway,
		telegramFormatter,
		lineFormatter,
	)
	appLogger.Info("adapter layer initialized")

	// ============================================================
	// 建立 Application 層（Use Cases）
	// ============================================================
	appLogger.Info("initializing use case layer...")
	// Stock Use Cases
	marketDataUsecase := stock.NewMarketDataUsecase(
		marketDataGateway,
		validationGateway,
		tradeDateRepo,
		appLogger,
	)

	marketChartUsecase := stock.NewMarketDataChartUsecase(
		marketChartGateway,
		validationGateway,
		appLogger,
	)

	// User Subscription Use Case
	userSubscriptionUsecase := user.NewUserSubscriptionUsecase(
		userRepo,                // UserAccountPort
		userSubscriptionGateway, // UserSubscriptionPort
		validationGateway,       // ValidationPort
	)

	// Bot Command Use Case
	botCommandUsecase := bot.NewBotCommandUsecase(
		formatterGateway,
		marketDataUsecase,
		marketChartUsecase,
		userSubscriptionUsecase,
	)

	// Health Check Use Case
	healthChecker := healthAdapter.NewHealthChecker(gormDB, finmindAPI, fugleAPI, syncMetadataRepo)
	healthUsecaseInstance := healthUsecase.NewHealthCheckUsecase(healthChecker, "stock-bot", "1.0.0", appLogger)

	// Bot Platform Use Cases
	tgCommandUsecase := bot.NewTgBotCommandUsecase(
		formatterGateway,
		botCommandUsecase,
		marketDataUsecase,
		userRepo,
		tgClient,
		appLogger,
	)

	lineCommandUsecase := bot.NewLineBotCommandUsecase(
		botCommandUsecase,
		lineClient,
		imgbbClient,
	)
	appLogger.Info("use case layer initialized")

	// ============================================================
	// 建立 Interfaces 層（Message Processors）
	// ============================================================
	appLogger.Info("initializing message processor layer...")
	tgProcessor := bot.NewTelegramMessageProcessor(
		tgCommandUsecase,
		userRepo,
		tgClient,
		appLogger,
	)

	lineProcessor := bot.NewLineMessageProcessor(
		lineCommandUsecase,
		userRepo,
		lineClient,
		appLogger,
	)
	appLogger.Info("message processor layer initialized")

	// ============================================================
	// 啟動 Web 服務器（HTTP Handlers）
	// ============================================================
	appLogger.Info("starting web server...")

	// 建立 Gin Router
	router, err := setupRouter(cfg, tgProcessor, lineProcessor, lineClient, healthUsecaseInstance, appLogger)
	if err != nil {
		appLogger.Fatal("failed to setup router", logger.Error(err))
	}

	appLogger.Info("application started",
		logger.String("port", "8080"),
		logger.String("tg_webhook", cfg.TELEGRAM_BOT_WEBHOOK_PATH),
		logger.String("line_webhook", cfg.LINE_BOT_WEBHOOK_PATH))

	// 在 goroutine 中啟動服務器
	go func() {
		if err := router.Run(":8080"); err != nil {
			appLogger.Fatal("HTTP server start failed", logger.Error(err))
		}
	}()

	// ============================================================
	// 優雅關閉
	// ============================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down application...")

	// cleanup resources
	ctx := context.Background()
	_ = ctx // context for graceful shutdown

	appLogger.Info("application shutdown complete")
}

// setupRouter 設定 HTTP 路由
func setupRouter(
	cfg *config.Config,
	tgProcessor *bot.TelegramMessageProcessor,
	lineProcessor *bot.LineMessageProcessor,
	lineClient *linebotInfra.LineBotClient,
	healthUsecase healthUsecase.HealthCheckUsecase,
	log logger.Logger,
) (*gin.Engine, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}
	if tgProcessor == nil {
		return nil, fmt.Errorf("telegram processor must not be nil")
	}
	if lineProcessor == nil {
		return nil, fmt.Errorf("line processor must not be nil")
	}
	if lineClient == nil {
		return nil, fmt.Errorf("line client must not be nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger must not be nil")
	}

	router := gin.New()
	if zapL, ok := logger.ExtractZapLogger(log); ok {
		router.Use(ginzap.Ginzap(zapL, time.RFC3339, true))
		router.Use(ginzap.RecoveryWithZap(zapL, true))
	} else {
		router.Use(gin.Recovery())
	}
	// 為每個 HTTP request 產生唯一 request_id 並注入 context
	router.Use(middleware.RequestID(log))

	// 健康檢查端點
	healthHandlerInstance := healthHandler.NewHealthHandler(healthUsecase, log)
	router.GET("/health", healthHandlerInstance.HealthCheck)

	// Telegram Webhook
	tgHandler := telegram.NewTgHandler(cfg, tgProcessor, log)
	telegram.RegisterRoutes(router, tgHandler, cfg.TELEGRAM_BOT_WEBHOOK_PATH)
	log.Info("Telegram Webhook registered")

	// LINE Webhook
	lineHandler := linebot.NewLineBotHandler(lineClient, lineProcessor, log)
	linebot.RegisterRoutes(router, lineHandler, cfg.LINE_BOT_WEBHOOK_PATH)
	log.Info("LINE Webhook registered")

	return router, nil
}
