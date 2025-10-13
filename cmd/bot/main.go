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
	"stock-bot/internal/api/tgbot"
	"stock-bot/internal/db"
	cnyesInfra "stock-bot/internal/infrastructure/cnyes"
	"stock-bot/internal/infrastructure/finmindtrade"
	fugleInfra "stock-bot/internal/infrastructure/fugle"
	"stock-bot/internal/infrastructure/imgbb"
	linebotInfra "stock-bot/internal/infrastructure/linebot"
	tgbotInfra "stock-bot/internal/infrastructure/tgbot"
	twseInfra "stock-bot/internal/infrastructure/twse"
	"stock-bot/internal/repository"
	lineService "stock-bot/internal/service/bot/line"
	tgService "stock-bot/internal/service/bot/tg"
	twstockService "stock-bot/internal/service/twstock"
	"stock-bot/internal/service/user"
	"stock-bot/pkg/logger"

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

	// 初始化 Repository
	userRepo := repository.NewUserRepository(db.GetDB())
	symbolsRepo := repository.NewSymbolRepository(db.GetDB())
	userSubscriptionRepo := repository.NewUserSubscriptionRepository(db.GetDB())

	// 初始化外部 API 客戶端
	fugleAPI := fugleInfra.NewFugleAPI(*cfg)
	finmindClient := finmindtrade.NewFinmindTradeAPI(*cfg)
	twseAPI := twseInfra.NewTwseAPI()
	cnyesAPI := cnyesInfra.NewCnyesAPI()

	// 初始化 ImgBB 客戶端
	var imgbbClient *imgbb.ImgBBClient
	if cfg.IMGBB_API_KEY != "" {
		imgbbClient = imgbb.NewImgBBClient(cfg.IMGBB_API_KEY)
		logger.Log.Info("ImgBB 客戶端初始化成功")
	} else {
		logger.Log.Warn("IMGBB_API_KEY 未設定，圖片上傳功能將不可用")
	}

	// 初始化服務
	userService := user.NewUserService(userRepo)
	stockService := twstockService.NewStockService(finmindClient, twseAPI, cnyesAPI, fugleAPI, symbolsRepo)

	// 設定 Gin 模式（根據環境變數自動設定）
	// 在 Docker 環境中，GIN_MODE 環境變數會自動設定為 release
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // 預設為 debug 模式（開發環境）
	}
	gin.SetMode(ginMode)
	logger.Log.Info("Gin 模式設定", zap.String("mode", ginMode))

	// 建立 Gin 引擎與註冊路由
	router := gin.Default()

	// 健康檢查
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 初始化 LINE Bot 並註冊路由
	botClient, err := linebotInfra.NewBot(*cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化 LINE Bot 失敗: %v", err))
	}

	// 建立 LINE Bot 服務層
	lineSvc := lineService.NewLineService(stockService, userSubscriptionRepo)
	lineCommandHandler := lineService.NewLineCommandHandler(botClient.Client, lineSvc, userService, userSubscriptionRepo, imgbbClient)
	service := lineService.NewBotService(botClient, lineCommandHandler, userService)
	handler := linebot.NewLineBotHandler(service, botClient)
	linebot.RegisterRoutes(router, handler)

	// 初始化 Telegram Bot 並註冊路由
	tgClient, err := tgbotInfra.NewBot(*cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化 Telegram Bot 失敗: %v", err))
	}
	tgSvc := tgService.NewTgService(stockService, userSubscriptionRepo)
	tgCommandHandler := tgService.NewTgCommandHandler(tgClient.Client, tgSvc, userService, userSubscriptionRepo)
	tgServiceHandler := tgService.NewTgHandler(tgCommandHandler, userService)
	tgHandler := tgbot.NewTgHandler(cfg, tgServiceHandler)
	tgbot.RegisterRoutes(router, tgHandler, cfg.TELEGRAM_BOT_WEBHOOK_PATH)

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
