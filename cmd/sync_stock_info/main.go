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
	finmindAPI := finmindtrade.NewFinmindTradeAPI(*cfg)
	fugleAPI := fugle.NewFugleAPI(*cfg)
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

		appLogger.Info("health check server started", logger.String("port", "8081"))
		if err := router.Run(":8081"); err != nil {
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

	appLogger.Info("running initial sync...")
	if err := stockSyncUsecase.SyncTaiwanStockInfo(ctx); err != nil {
		appLogger.Error("TW stock sync failed", logger.Error(err))
	}
	appLogger.Info("TW stock sync completed")

	if err := stockSyncUsecase.SyncUSStockInfo(ctx); err != nil {
		appLogger.Error("US stock sync failed", logger.Error(err))
	}
	appLogger.Info("US stock sync completed")

	if stats, err := stockSyncUsecase.GetSyncStats(ctx); err == nil {
		appLogger.Info("initial sync stats", logger.Any("stats", stats))
	}

	err := stockSyncUsecase.SyncTaiwanStockTradingDate(ctx)
	if err != nil {
		appLogger.Error("failed to sync TW trade dates", logger.Error(err))
	}
	appLogger.Info("TW trade date sync completed")

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			appLogger.Info("received stop signal, stopping background sync...")
			return
		case <-ticker.C:
			appLogger.Info("running scheduled sync...")
			if err := stockSyncUsecase.SyncTaiwanStockInfo(ctx); err != nil {
				appLogger.Error("TW stock sync failed", logger.Error(err))
			}
			appLogger.Info("TW stock sync completed")

			if err := stockSyncUsecase.SyncUSStockInfo(ctx); err != nil {
				appLogger.Error("US stock sync failed", logger.Error(err))
			}
			appLogger.Info("US stock sync completed")

			if stats, err := stockSyncUsecase.GetSyncStats(ctx); err == nil {
				appLogger.Info("scheduled sync stats", logger.Any("stats", stats))
			}
		}
	}
}
