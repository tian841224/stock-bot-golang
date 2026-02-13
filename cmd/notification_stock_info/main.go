package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	healthUsecase "github.com/tian841224/stock-bot/internal/application/usecase/health"
	notificationUseCase "github.com/tian841224/stock-bot/internal/application/usecase/notification"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	formatterAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/formatter"
	healthAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/health"
	marketAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/market"
	presenterAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/presenter"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	tgbotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/twse"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	database "github.com/tian841224/stock-bot/internal/infrastructure/persistence"
	repository "github.com/tian841224/stock-bot/internal/infrastructure/persistence/postgres"
	healthHandler "github.com/tian841224/stock-bot/internal/interfaces/health"
)

func main() {
	// ============================================================
	// 基礎設施初始化
	// ============================================================
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("載入配置失敗: %v", err)
	}

	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("初始化 Logger 失敗: %v", err)
	}

	appLogger.Info("=== 通知服務啟動中 ===")

	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("初始化資料庫失敗", logger.Error(err))
	}
	defer db.Close()

	gormDB := db.GetDB()

	// ============================================================
	// 外部服務客戶端
	// ============================================================
	tgClient, err := tgbotInfra.NewBot(*cfg, appLogger)
	if err != nil {
		appLogger.Fatal("建立 Telegram Bot 客戶端失敗", logger.Error(err))
	}

	fugleAPI := fugle.NewFugleAPI(*cfg)
	twseAPI := twse.NewTwseAPI()
	cnyesAPI := cnyes.NewCnyesAPI()
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg)

	// ============================================================
	// Repository
	// ============================================================
	stockSymbolRepo := repository.NewSymbolRepository(gormDB, appLogger)
	tradeDateRepo := repository.NewPostgresTradeDateRepository(gormDB, appLogger)
	subscriptionSymbolRepo := repository.NewSubscriptionSymbolRepository(gormDB, appLogger)
	syncMetadataRepo := repository.NewSyncMetadataRepository(gormDB, appLogger)

	// ============================================================
	// Health Check
	// ============================================================
	healthChecker := healthAdapter.NewHealthChecker(gormDB, finmindAPI, fugleAPI, syncMetadataRepo)
	healthUsecaseInstance := healthUsecase.NewHealthCheckUsecase(healthChecker, "stock-scheduler", "1.0.0", appLogger)

	// ============================================================
	// Adapter / Gateway
	// ============================================================
	validationGateway := presenterAdapter.NewValidationGateway(nil, stockSymbolRepo)

	marketDataGateway := marketAdapter.NewMarketDataGateway(
		twseAPI,
		cnyesAPI,
		fugleAPI,
		finmindAPI,
		validationGateway,
		tradeDateRepo,
	)

	marketChartGateway := marketAdapter.NewMarketChartGateway(
		marketDataGateway,
		validationGateway,
		fugleAPI,
	)

	telegramFormatter := formatterAdapter.NewTelegramFormatter()
	lineFormatter := formatterAdapter.NewLineFormatter()
	formatterGateway := formatterAdapter.NewFormatterAdapter(
		marketChartGateway,
		validationGateway,
		telegramFormatter,
		lineFormatter,
	)

	// ============================================================
	// Use Case
	// ============================================================
	marketDataUsecase := stock.NewMarketDataUsecase(
		marketDataGateway,
		validationGateway,
		tradeDateRepo,
		appLogger,
	)

	sendNotificationUsecase := notificationUseCase.NewSendNotificationUsecase(
		subscriptionSymbolRepo,
		marketDataUsecase,
		formatterGateway,
		tgClient,
		appLogger,
	)

	scheduleHandlerUsecase := notificationUseCase.NewScheduleHandlerUsecase(sendNotificationUsecase, appLogger)
	appLogger.Info("所有服務初始化完成")

	// ============================================================
	// 啟動服務
	// ============================================================
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 啟動健康檢查 HTTP 服務器
	go func() {
		router := gin.Default()
		healthHandlerInstance := healthHandler.NewHealthHandler(healthUsecaseInstance, appLogger)
		router.GET("/health", healthHandlerInstance.HealthCheck)

		appLogger.Info("健康檢查服務器啟動，監聽端口: 8081")
		if err := router.Run(":8081"); err != nil {
			appLogger.Error("健康檢查服務器啟動失敗", logger.Error(err))
		}
	}()

	// 啟動排程通知任務
	go runScheduledNotifications(ctx, scheduleHandlerUsecase, appLogger)

	<-quit
	appLogger.Info("收到關閉信號，正在優雅關閉...")
	cancel()
	appLogger.Info("=== 通知服務已關閉 ===")
}

func runScheduledNotifications(ctx context.Context, scheduler notificationUseCase.ScheduleHandlerUsecase, log logger.Logger) {
	log.Info("排程通知服務已啟動")

	// 啟動時立刻執行一次
	log.Info("正在執行啟動時的通知任務...")
	if err := scheduler.RunScheduledTasks(ctx); err != nil {
		log.Error("啟動時通知任務執行失敗", logger.Error(err))
	} else {
		log.Info("啟動時通知任務執行完成")
	}

	// 設定每天下午三點 (台北時間) 執行
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Error("無法載入 Asia/Taipei 時區，將使用 Local 時區", logger.Error(err))
		loc = time.Local
	}

	for {
		now := time.Now().In(loc)
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, loc)

		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		duration := nextRun.Sub(now)
		log.Info("下次排程任務將在 " + duration.String() + " 後執行 (" + nextRun.Format("2006-01-02 15:04:05") + ")")

		select {
		case <-ctx.Done():
			log.Info("排程通知服務已停止")
			return
		case <-time.After(duration):
			log.Info("執行下午三點的排程通知任務...")
			if err := scheduler.RunScheduledTasks(ctx); err != nil {
				log.Error("排程通知任務執行失敗", logger.Error(err))
			} else {
				log.Info("排程通知任務執行完成")
			}
		}
	}
}
