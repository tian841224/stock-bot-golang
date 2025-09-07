package db

import (
	"fmt"
	"stock-bot/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(cfg *config.Config) {
	connectDB(cfg)
	createOrUpdateTable()
}

func connectDB(cfg *config.Config) {
	var err error

	// 設定日誌模式，從環境變數讀取
	var logLevel logger.LogLevel

	// 設定DB日誌模式
	if cfg.DB_LOG_MODE {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err = gorm.Open(postgres.Open(cfg.DATABASE_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic(fmt.Sprintf("連接資料庫失敗: %v", err))
	}
}

func createOrUpdateTable() {
	err := db.AutoMigrate()
	if err != nil {
		panic(fmt.Sprintf("資料表遷移失敗: %v", err))
	}
}
