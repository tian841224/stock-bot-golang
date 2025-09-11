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
	"stock-bot/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("載入設定失敗: %v", err))
	}

	// 初始化資料庫
	db.InitDB(cfg)

	// 建立 Gin 引擎與註冊路由
	router := gin.Default()

	// 健康檢查
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

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
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("啟動 HTTP 伺服器失敗: %v", err))
		}
	}()

	// 等待終止訊號，優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("伺服器關閉失敗: %v\n", err)
	}

	if err := db.Close(); err != nil {
		fmt.Printf("資料庫關閉失敗: %v\n", err)
	}
}
