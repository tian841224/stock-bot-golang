package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	appLogger.Info("應用程式啟動中...")
	appLogger.Info("載入配置成功")

	// 初始化資料庫
	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("初始化資料庫失敗", logger.Error(err))
	}
	defer db.Close()

	gormDB := db.GetDB()

	// ============================================================
	// 建立 Repository 層（Persistence）
	// ============================================================

	appLogger.Info("初始化 Repository 層...")
	userRepo := repository.NewPostgresUserRepository(gormDB)
	stockSymbolRepo := repository.NewSymbolRepository(gormDB)
	tradeDateRepo := repository.NewPostgresTradeDateRepository(gormDB)
	subscriptionRepo := repository.NewSubscriptionRepository(gormDB)
	subscriptionSymbolRepo := repository.NewSubscriptionSymbolRepository(gormDB)
	featureReader, _ := repository.NewFeatureRepository(gormDB)
	syncMetadataRepo := repository.NewSyncMetadataRepository(gormDB)
	appLogger.Info("Feature Repository 初始化成功，預設功能資料已建立")

	// ============================================================
	// 建立外部服務客戶端（External Services）
	// ============================================================
	appLogger.Info("初始化外部服務客戶端...")

	// Telegram Bot 客戶端
	appLogger.Info("初始化 Telegram Bot 客戶端...")
	tgClient, err := tgbotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("建立 Telegram Bot 客戶端失敗", logger.Error(err))
	}
	appLogger.Info("Telegram Bot 客戶端初始化成功")

	// LINE Bot 客戶端
	appLogger.Info("初始化 LINE Bot 客戶端...")
	lineClient, err := linebotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("建立 LINE Bot 客戶端失敗", logger.Error(err))
	}
	appLogger.Info("LINE Bot 客戶端初始化成功")

	// 圖片上傳服務
	appLogger.Info("初始化外部服務客戶端...")
	imgbbClient := imgbb.NewImgBBClient(cfg.IMGBB_API_KEY)

	// 股票 API 客戶端
	fugleAPI := fugle.NewFugleAPI(*cfg)
	twseAPI := twse.NewTwseAPI()
	cnyesAPI := cnyes.NewCnyesAPI()
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg)
	appLogger.Info("外部服務客戶端初始化成功")

	// ============================================================
	// 建立 Adapter 層（Gateway/Presenter）
	// ============================================================
	appLogger.Info("初始化 Adapter 層...")
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
	appLogger.Info("Adapter 層初始化成功")

	// ============================================================
	// 建立 Application 層（Use Cases）
	// ============================================================
	appLogger.Info("初始化 Use Case 層...")
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
	appLogger.Info("Use Case 層初始化成功")

	// ============================================================
	// 建立 Interfaces 層（Message Processors）
	// ============================================================
	appLogger.Info("初始化 Message Processor 層...")
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
	appLogger.Info("Message Processor 層初始化成功")

	// ============================================================
	// 啟動 Web 服務器（HTTP Handlers）
	// ============================================================
	appLogger.Info("正在啟動 Web 服務器...")

	// 建立 Gin Router
	router, err := setupRouter(cfg, tgProcessor, lineProcessor, lineClient, healthUsecaseInstance, appLogger)
	if err != nil {
		appLogger.Fatal("設定路由失敗", logger.Error(err))
	}

	appLogger.Info("Web 服務器已啟動，監聽端口: 8080")
	appLogger.Info("Telegram Webhook 路徑: " + cfg.TELEGRAM_BOT_WEBHOOK_PATH)
	appLogger.Info("LINE Webhook 路徑: " + cfg.LINE_BOT_WEBHOOK_PATH)

	// 在 goroutine 中啟動服務器
	go func() {
		if err := router.Run(":8080"); err != nil {
			appLogger.Fatal("HTTP 服務器啟動失敗", logger.Error(err))
		}
	}()

	// ============================================================
	// 優雅關閉
	// ============================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("正在關閉應用程式...")

	// 清理資源
	ctx := context.Background()
	_ = ctx // 用於優雅關閉的 context

	appLogger.Info("應用程式已關閉")
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
		return nil, fmt.Errorf("配置不能為空")
	}
	if tgProcessor == nil {
		return nil, fmt.Errorf("Telegram 處理器不能為空")
	}
	if lineProcessor == nil {
		return nil, fmt.Errorf("LINE 處理器不能為空")
	}
	if lineClient == nil {
		return nil, fmt.Errorf("LINE 客戶端不能為空")
	}
	if log == nil {
		return nil, fmt.Errorf("Logger 不能為空")
	}

	router := gin.Default()

	// 健康檢查端點
	healthHandlerInstance := healthHandler.NewHealthHandler(healthUsecase, log)
	router.GET("/health", healthHandlerInstance.HealthCheck)

	// Telegram Webhook
	tgHandler := telegram.NewTgHandler(cfg, tgProcessor, log)
	telegram.RegisterRoutes(router, tgHandler, cfg.TELEGRAM_BOT_WEBHOOK_PATH)
	log.Info("Telegram Webhook 路由註冊成功")

	// LINE Webhook
	lineHandler := linebot.NewLineBotHandler(lineClient, lineProcessor, log)
	linebot.RegisterRoutes(router, lineHandler, cfg.LINE_BOT_WEBHOOK_PATH)
	log.Info("LINE Webhook 路由註冊成功")

	return router, nil
}
