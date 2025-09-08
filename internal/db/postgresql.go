package db

import (
	"database/sql"
	"fmt"
	"stock-bot/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(cfg *config.Config) {
	createDatabaseIfNotExists(cfg)
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

	// 設定資料庫連接字串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_PORT)
	fmt.Println(dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
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

// createDatabaseIfNotExists 檢查並建立資料庫
func createDatabaseIfNotExists(cfg *config.Config) error {
	sqlDB, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_PORT))
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = sqlDB.QueryRow(query, cfg.DB_NAME).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = sqlDB.Exec("CREATE DATABASE " + cfg.DB_NAME)
		if err != nil {
			return err
		}
		fmt.Printf("資料庫 %s 建立成功\n", cfg.DB_NAME)
	}

	return nil
}
