// Package db 提供資料庫連線與操作功能
package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/config"
	"github.com/tian841224/stock-bot/internal/db/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// GetDB 回傳資料庫連線實例
func GetDB() *gorm.DB {
	return db
}

func InitDB(cfg *config.Config) error {
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return err
	}
	if err := connectDB(cfg); err != nil {
		return err
	}
	if err := createOrUpdateTable(); err != nil {
		return err
	}
	return nil
}

func connectDB(cfg *config.Config) error {
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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return err
	}

	// 設定連線池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return nil
}

func createOrUpdateTable() error {
	allModels := models.AllModels()

	if err := db.AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("資料庫遷移失敗: %w", err)
	}

	return nil
}

// createDatabaseIfNotExists 檢查並建立資料庫
func createDatabaseIfNotExists(cfg *config.Config) error {
	sqlDB, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD))
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
	}

	return nil
}

// Close 關閉資料庫連線
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
