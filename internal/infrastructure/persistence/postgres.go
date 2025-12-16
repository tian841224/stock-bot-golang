// Package db 提供資料庫連線與操作功能
package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	model "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 定義資料庫操作介面
type Database interface {
	Init(cfg *config.Config) error
	GetDB() *gorm.DB
	Close() error
}

// postgresDatabase 實作 Database 介面
type postgresDatabase struct {
	db *gorm.DB
}

// NewDatabase 建立新的 Database 實例
func NewDatabase() Database {
	return &postgresDatabase{}
}

// Init 初始化資料庫連線
func (d *postgresDatabase) Init(cfg *config.Config) error {
	if err := d.createDatabaseIfNotExists(cfg); err != nil {
		return err
	}
	if err := d.connectDB(cfg); err != nil {
		return err
	}
	if err := d.createOrUpdateTable(); err != nil {
		return err
	}
	return nil
}

// GetDB 回傳資料庫連線實例
func (d *postgresDatabase) GetDB() *gorm.DB {
	return d.db
}

// Close 關閉資料庫連線
func (d *postgresDatabase) Close() error {
	if d.db == nil {
		return nil
	}
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *postgresDatabase) connectDB(cfg *config.Config) error {
	var err error

	var logLevel logger.LogLevel
	if cfg.DB_LOG_MODE {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return err
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return nil
}

func (d *postgresDatabase) createOrUpdateTable() error {
	allModels := model.AllModels()

	if err := d.db.AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("資料庫遷移失敗: %w", err)
	}

	return nil
}

func (d *postgresDatabase) createDatabaseIfNotExists(cfg *config.Config) error {
	sqlDB, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			// 忽略關閉錯誤，因為這只是用於檢查資料庫的臨時連線
		}
	}()

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

// 向後相容的全域變數和函數
var db *gorm.DB

// GetDB 回傳資料庫連線實例（向後相容）
func GetDB() *gorm.DB {
	return db
}

// InitDB 初始化資料庫（向後相容）
func InitDB(cfg *config.Config) error {
	database := NewDatabase()
	if err := database.Init(cfg); err != nil {
		return err
	}
	db = database.GetDB()
	return nil
}

// Close 關閉資料庫連線（向後相容）
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
