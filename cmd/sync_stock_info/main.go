package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stock-bot/config"
	"stock-bot/internal/db"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/repository"
	"stock-bot/internal/service/stock_sync"
	"stock-bot/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 初始化日誌
	logger.InitLogger()
	defer logger.Log.Sync()

	logger.Log.Info("=== 股票資料同步程式啟動 ===")

	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Panic("載入設定失敗", zap.Error(err))
	}
	logger.Log.Info("設定載入成功")

	// 初始化資料庫
	if err := db.InitDB(cfg); err != nil {
		logger.Log.Panic("資料庫初始化失敗", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Log.Error("資料庫關閉失敗", zap.Error(err))
		}
	}()
	logger.Log.Info("資料庫初始化成功")

	// 初始化 Repository 和 Service
	symbolsRepo := repository.NewSymbolsRepository(db.GetDB())
	finmindClient := finmindtrade.NewFinmindTradeAPI(*cfg)
	stockSyncService := stock_sync.NewStockSyncService(symbolsRepo, finmindClient)
	logger.Log.Info("服務初始化成功")

	// 建立 context 用於優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 監聽系統中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 啟動背景同步服務
	go runBackgroundSync(ctx, stockSyncService)

	// 等待中斷信號
	<-quit
	logger.Log.Info("收到關閉信號，正在優雅關閉...")

	// 取消 context，停止背景任務
	cancel()

	logger.Log.Info("=== 程式已關閉 ===")
}

// runBackgroundSync 執行背景同步任務
func runBackgroundSync(ctx context.Context, stockSyncService *stock_sync.StockSyncService) {
	defer func() {
		logger.Log.Info("背景同步任務已完全停止")
	}()

	// 程式啟動時立即執行一次同步
	logger.Log.Info("執行初始同步...")
	if err := stockSyncService.SyncTaiwanStockInfo(); err != nil {
		logger.Log.Error("初始同步失敗", zap.Error(err))
	} else {
		// 顯示同步統計
		if stats, err := stockSyncService.GetSyncStats(); err == nil {
			logger.Log.Info("初始同步統計", zap.Any("stats", stats))
		}
	}

	// 建立定時器
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("收到停止信號，背景同步任務正在關閉...")
			return
		case <-ticker.C:
			logger.Log.Info("開始定時同步...")
			if err := stockSyncService.SyncTaiwanStockInfo(); err != nil {
				logger.Log.Error("定時同步失敗", zap.Error(err))
			} else {
				// 顯示同步統計
				if stats, err := stockSyncService.GetSyncStats(); err == nil {
					logger.Log.Info("定時同步統計", zap.Any("stats", stats))
				}
			}
		}
	}
}
