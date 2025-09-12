package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stock-bot/config"
	"stock-bot/internal/api/linebot"
	"stock-bot/internal/db"
	"stock-bot/pkg/logger"
	linebotInfra "stock-bot/internal/infrastructure/linebot"
	lineService "stock-bot/internal/service/bot/line"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// 初始化日誌
	logger.InitLogger()
	defer logger.Log.Sync()

	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("載入設定失敗: %v", err))
	}
	logger.Log.Info("設定載入成功")

	// 初始化資料庫
	if err := db.InitDB(cfg); err != nil {
		logger.Log.Panic("資料庫初始化失敗", zap.Error(err))
	}
	logger.Log.Info("資料庫初始化成功")

	// 建立 Gin 引擎與註冊路由
	router := gin.Default()

	// 健康檢查
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 初始化 LINE Bot 依賴並註冊路由
	botClient, err := linebotInfra.NewBot(*cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化 LINE Bot 失敗: %v", err))
	}
	service := lineService.NewBotService(botClient)
	handler := linebot.NewLineBotHandler(service, botClient)
	linebot.RegisterRoutes(router, handler)

	// 從環境變數讀取埠號，預設 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 啟動伺服器（背景）
	serverErr := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	logger.Log.Info("HTTP 伺服器啟動成功")
	logger.Log.Info("程式執行中...")

	// 等待終止訊號或啟動錯誤
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		// 繼續往下優雅關閉
	case err := <-serverErr:
		logger.Log.Error("啟動 HTTP 伺服器失敗", zap.Error(err))
		// 立刻以非 0 退出，先同步日誌
		_ = logger.Log.Sync()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	if shutdownErr != nil {
		logger.Log.Error("伺服器關閉失敗", zap.Error(shutdownErr))
	}

	dbErr := db.Close()
	if dbErr != nil {
		logger.Log.Error("資料庫關閉失敗", zap.Error(dbErr))
	}

	if shutdownErr != nil || dbErr != nil {
		// 有關閉錯誤，用非 0 退出；先同步日誌
		_ = logger.Log.Sync()
		os.Exit(1)
	}
}
