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

	appLogger.Info("應用程式啟動中...")
	appLogger.Info("載入配置成功")

	// 初始化資料庫
	db := database.NewDatabase()
	if err := db.Init(cfg); err != nil {
		appLogger.Fatal("初始化資料庫失敗", logger.Error(err))
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

	appLogger.Info("服務初始化成功")

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

		appLogger.Info("健康檢查服務器啟動，監聽端口: 8081")
		if err := router.Run(":8081"); err != nil {
			appLogger.Error("健康檢查服務器啟動失敗", logger.Error(err))
		}
	}()

	// 等待中斷信號
	<-quit
	appLogger.Info("收到關閉信號，正在優雅關閉...")

	// 取消 context，停止背景任務
	cancel()

	appLogger.Info("=== 程式已關閉 ===")
}

func runBackgroundSync(ctx context.Context, stockSyncUsecase stock_sync.StockSyncUsecase, appLogger logger.Logger) {
	defer func() {
		appLogger.Info("背景同步任務已完全停止")
	}()

	appLogger.Info("執行初始同步...")
	if err := stockSyncUsecase.SyncTaiwanStockInfo(ctx); err != nil {
		appLogger.Error("台股同步失敗", logger.Error(err))
	}
	appLogger.Info("台股同步完成")

	if err := stockSyncUsecase.SyncUSStockInfo(ctx); err != nil {
		appLogger.Error("美股同步失敗", logger.Error(err))
	}
	appLogger.Info("美股同步完成")

	if stats, err := stockSyncUsecase.GetSyncStats(ctx); err == nil {
		appLogger.Info("初始同步統計", logger.Any("stats", stats))
	}

	err := stockSyncUsecase.SyncTaiwanStockTradingDate(ctx)
	if err != nil {
		appLogger.Error("取得台股交易日失敗", logger.Error(err))
	}
	appLogger.Info("台股交易日同步完成")

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			appLogger.Info("收到停止信號，背景同步任務正在關閉...")
			return
		case <-ticker.C:
			appLogger.Info("開始定時同步...")
			if err := stockSyncUsecase.SyncTaiwanStockInfo(ctx); err != nil {
				appLogger.Error("台股同步失敗", logger.Error(err))
			}
			appLogger.Info("台股同步完成")

			if err := stockSyncUsecase.SyncUSStockInfo(ctx); err != nil {
				appLogger.Error("美股同步失敗", logger.Error(err))
			}
			appLogger.Info("美股同步完成")

			if stats, err := stockSyncUsecase.GetSyncStats(ctx); err == nil {
				appLogger.Info("定時同步統計", logger.Any("stats", stats))
			}
		}
	}
}
