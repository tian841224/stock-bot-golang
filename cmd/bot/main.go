package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tian841224/stock-bot/internal/application/usecase/bot"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	formatterAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/formatter"
	marketAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/market"
	validationAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/presenter"
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
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("載入配置失敗: %v", err)
	}

	// 初始化 Logger
	appLogger, _ := logger.NewLogger()
	appLogger.Info("應用程式啟動中...")

	// 初始化資料庫
	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("初始化資料庫失敗", logger.Error(err))
	}
	defer db.Close()

	gormDB := db.GetDB()

	// ============================================================
	// 2. 建立 Repository 層（Persistence）
	// ============================================================
	userRepo := repository.NewPostgresUserRepository(gormDB)
	stockSymbolRepo := repository.NewSymbolRepository(gormDB)
	tradeDateRepo := repository.NewPostgresTradeDateRepository(gormDB)

	// ============================================================
	// 3. 建立外部服務客戶端（External Services）
	// ============================================================
	// Bot 客戶端
	tgClient, err := tgbotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("建立 Telegram Bot 客戶端失敗", logger.Error(err))
	}

	lineClient, err := linebotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("建立 LINE Bot 客戶端失敗", logger.Error(err))
	}

	// 圖片上傳服務
	imgbbClient := imgbb.NewImgBBClient(cfg.IMGBB_API_KEY)

	// 股票 API 客戶端
	fugleAPI := fugle.NewFugleAPI(*cfg)
	twseAPI := twse.NewTwseAPI()
	cnyesAPI := cnyes.NewCnyesAPI()
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg)

	// ============================================================
	// 4. 建立 Adapter 層（Gateway/Presenter）
	// ============================================================
	// Validation Gateway
	validationGateway := validationAdapter.NewValidationGateway(nil, stockSymbolRepo)

	// Market Data Gateway
	marketDataGateway := marketAdapter.NewMarketDataGateway(
		nil,
		twseAPI,
		cnyesAPI,
		fugleAPI,
		finmindAPI,
		validationGateway,
		tradeDateRepo,
	)

	// Market Chart Gateway
	marketChartGateway := marketAdapter.NewMarketChartGateway(
		nil,
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
	// ============================================================
	// 5. 建立 Application 層（Use Cases）
	// ============================================================
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

	// Bot Command Use Case
	botCommandUsecase := bot.NewBotCommandUsecase(
		nil, // botCommandPort - 這是介面層，不需要在這裡注入
		formatterGateway,
		marketDataUsecase,
		marketChartUsecase,
	)

	// Bot Platform Use Cases
	tgCommandUsecase := bot.NewTgBotCommandUsecase(
		formatterGateway,
		botCommandUsecase,
		marketDataUsecase,
		tgClient,
	)

	lineCommandUsecase := bot.NewLineBotCommandUsecase(
		botCommandUsecase,
		lineClient,
		imgbbClient,
	)

	// ============================================================
	// 6. 建立 Interfaces 層（Message Processors）
	// ============================================================
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

	// ============================================================
	// 7. 啟動 Web 服務器（HTTP Handlers）
	// ============================================================
	appLogger.Info("正在啟動 Web 服務器...")

	// 建立 Gin Router
	router := setupRouter(cfg, tgProcessor, lineProcessor, lineClient, appLogger)

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
	// 8. 優雅關閉
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
	log logger.Logger,
) *gin.Engine {
	router := gin.Default()

	// 健康檢查端點
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Telegram Webhook
	tgHandler := telegram.NewTgHandler(cfg, tgProcessor, log)
	telegram.RegisterRoutes(router, tgHandler, cfg.TELEGRAM_BOT_WEBHOOK_PATH)

	// LINE Webhook
	lineHandler := linebot.NewLineBotHandler(lineClient, lineProcessor, log)
	linebot.RegisterRoutes(router, lineHandler, cfg.LINE_BOT_WEBHOOK_PATH)

	return router
}
