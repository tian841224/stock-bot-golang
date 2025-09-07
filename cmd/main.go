package main

import (
	"fmt"
	"stock-bot/config"
	"stock-bot/internal/db"
)

func main() {
	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("載入設定失敗: %v", err))
	}

	// 初始化資料庫
	db.InitDB(cfg)
}
