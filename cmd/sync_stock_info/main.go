package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	healthUsecase "github.com/tian841224/stock-bot/internal/application/usecase/health"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock_sync"
	healthAdapter "github.com/tian841224/stock-bot/internal/infrastructure/adapter/health"
	stock "github.com/tian841224/stock-bot/internal/infrastructure/adapter/stock"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	database "github.com/tian841224/stock-bot/internal/infrastructure/persistence"
	repository "github.com/tian841224/stock-bot/internal/infrastructure/persistence/postgres"
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

	appLogger.Info("starting application...")
	appLogger.Info("config loaded")

	// 初始化資料庫
	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("failed to init database", logger.Error(err))
	}
	defer db.Close()

	gormDB := db.GetDB()

	stockSymbolRepo := repository.NewSymbolRepository(gormDB, appLogger)
	syncMetadataRepo := repository.NewSyncMetadataRepository(gormDB, appLogger)
	tradeDateRepo := repository.NewPostgresTradeDateRepository(gormDB, appLogger)
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg, appLogger)
	fugleAPI := fugle.NewFugleAPI(*cfg, appLogger)
	stockInfoProvider := stock.NewFinmindStockInfoAdapter(finmindAPI)
	stockSyncUsecase := stock_sync.NewStockSyncUsecase(stockSymbolRepo, stockInfoProvider, syncMetadataRepo, tradeDateRepo, appLogger)

	healthChecker := healthAdapter.NewHealthChecker(gormDB, finmindAPI, fugleAPI, syncMetadataRepo)
	healthUsecaseInstance := healthUsecase.NewHealthCheckUsecase(healthChecker, "stock-sync", "1.0.0", appLogger)

	appLogger.Info("all services initialized")

	// 建立 context 用於優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 監聽系統中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 啟動背景同步服務
	go runBackgroundSync(ctx, stockSyncUsecase, appLogger)

	// 啟動健康檢查 HTTP 服務器
	go func() {
		router := gin.Default()
		healthHandlerInstance := healthHandler.NewHealthHandler(healthUsecaseInstance, appLogger)
		router.GET("/health", healthHandlerInstance.HealthCheck)

		port := cfg.PORT
		if port == 0 {
			port = 8081
		}
		appLogger.Info("health check server started", logger.Int("port", port))
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			appLogger.Error("health check server start failed", logger.Error(err))
		}
	}()

	// 等待中斷信號
	<-quit
	appLogger.Info("received shutdown signal, shutting down gracefully...")

	// cancel context to stop background tasks
	cancel()

	appLogger.Info("=== sync service stopped ===")
}

func runBackgroundSync(ctx context.Context, stockSyncUsecase stock_sync.StockSyncUsecase, appLogger logger.Logger) {
	defer func() {
		appLogger.Info("background sync task fully stopped")
	}()

	// 初始執行：產生唯一 job_id，讓整次同步鏈路完整串聯
	initJobID := fmt.Sprintf("sync-init-%d", time.Now().UnixNano())
	initCtx := logger.WithLogger(ctx, appLogger.With(logger.String("job_id", initJobID)))

	appLogger.Info("running initial sync...", logger.String("job_id", initJobID))
	if err := stockSyncUsecase.SyncTaiwanStockInfo(initCtx); err != nil {
		appLogger.Error("TW stock sync failed", logger.String("job_id", initJobID), logger.Error(err))
	}
	appLogger.Info("TW stock sync completed", logger.String("job_id", initJobID))

	if err := stockSyncUsecase.SyncUSStockInfo(initCtx); err != nil {
		appLogger.Error("US stock sync failed", logger.String("job_id", initJobID), logger.Error(err))
	}
	appLogger.Info("US stock sync completed", logger.String("job_id", initJobID))

	if stats, err := stockSyncUsecase.GetSyncStats(initCtx); err == nil {
		appLogger.Info("initial sync stats", logger.String("job_id", initJobID), logger.Any("stats", stats))
	}

	err := stockSyncUsecase.SyncTaiwanStockTradingDate(initCtx)
	if err != nil {
		appLogger.Error("failed to sync TW trade dates", logger.String("job_id", initJobID), logger.Error(err))
	}
	appLogger.Info("TW trade date sync completed", logger.String("job_id", initJobID))

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			appLogger.Info("received stop signal, stopping background sync...")
			return
		case <-ticker.C:
			// 每次定時觸發產生獨立 job_id，方便追蹤單次排程的完整鏈路
			jobID := fmt.Sprintf("sync-%d", time.Now().UnixNano())
			jobCtx := logger.WithLogger(ctx, appLogger.With(logger.String("job_id", jobID)))

			appLogger.Info("running scheduled sync...", logger.String("job_id", jobID))
			if err := stockSyncUsecase.SyncTaiwanStockInfo(jobCtx); err != nil {
				appLogger.Error("TW stock sync failed", logger.String("job_id", jobID), logger.Error(err))
			}
			appLogger.Info("TW stock sync completed", logger.String("job_id", jobID))

			if err := stockSyncUsecase.SyncUSStockInfo(jobCtx); err != nil {
				appLogger.Error("US stock sync failed", logger.String("job_id", jobID), logger.Error(err))
			}
			appLogger.Info("US stock sync completed", logger.String("job_id", jobID))

			if stats, err := stockSyncUsecase.GetSyncStats(jobCtx); err == nil {
				appLogger.Info("scheduled sync stats", logger.String("job_id", jobID), logger.Any("stats", stats))
			}
		}
	}
}
