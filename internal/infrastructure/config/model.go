// Package config 提供應用程式設定管理
package config

import (
	"fmt"
	"strings"
)

type Config struct {
	LINE_BOT_WEBHOOK_PATH       string `mapstructure:"LINE_BOT_WEBHOOK_PATH"`
	TELEGRAM_BOT_SECRET_TOKEN   string `mapstructure:"TELEGRAM_BOT_SECRET_TOKEN"`
	DB_USER                     string `mapstructure:"DB_USER"`
	TELEGRAM_ADMIN_CHAT_ID      string `mapstructure:"TELEGRAM_ADMIN_CHAT_ID"`
	DB_NAME                     string `mapstructure:"DB_NAME"`
	SCHEDULER_STOCK_SPEC        string `mapstructure:"SCHEDULER_STOCK_SPEC"`
	CHANNEL_ACCESS_TOKEN        string `mapstructure:"CHANNEL_ACCESS_TOKEN"`
	CHANNEL_SECRET              string `mapstructure:"CHANNEL_SECRET"`
	SCHEDULER_TIMEZONE          string `mapstructure:"SCHEDULER_TIMEZONE"`
	TELEGRAM_BOT_TOKEN          string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	DB_PASSWORD                 string `mapstructure:"DB_PASSWORD"`
	TELEGRAM_BOT_WEBHOOK_DOMAIN string `mapstructure:"TELEGRAM_BOT_WEBHOOK_DOMAIN"`
	TELEGRAM_BOT_WEBHOOK_PATH   string `mapstructure:"TELEGRAM_BOT_WEBHOOK_PATH"`
	DB_HOST                     string `mapstructure:"DB_HOST"`
	FINMIND_TOKEN               string `mapstructure:"FINMIND_TOKEN"`
	FUGLE_API_KEY               string `mapstructure:"FUGLE_API_KEY"`
	IMGBB_API_KEY               string `mapstructure:"IMGBB_API_KEY"`
	DB_PORT                     int    `mapstructure:"DB_PORT"`
	DB_LOG_MODE                 bool   `mapstructure:"DB_LOG"`
}

// Validate 驗證配置的必要欄位
func (c *Config) Validate() error {
	var missingFields []string

	// 資料庫配置驗證
	if c.DB_HOST == "" {
		missingFields = append(missingFields, "DB_HOST")
	}
	if c.DB_USER == "" {
		missingFields = append(missingFields, "DB_USER")
	}
	if c.DB_NAME == "" {
		missingFields = append(missingFields, "DB_NAME")
	}
	if c.DB_PORT == 0 {
		missingFields = append(missingFields, "DB_PORT")
	}

	// Telegram Bot 配置驗證
	if c.TELEGRAM_BOT_TOKEN == "" {
		missingFields = append(missingFields, "TELEGRAM_BOT_TOKEN")
	}
	if c.TELEGRAM_BOT_WEBHOOK_PATH == "" {
		missingFields = append(missingFields, "TELEGRAM_BOT_WEBHOOK_PATH")
	}

	// LINE Bot 配置驗證
	if c.CHANNEL_ACCESS_TOKEN == "" {
		missingFields = append(missingFields, "CHANNEL_ACCESS_TOKEN")
	}
	if c.CHANNEL_SECRET == "" {
		missingFields = append(missingFields, "CHANNEL_SECRET")
	}
	if c.LINE_BOT_WEBHOOK_PATH == "" {
		missingFields = append(missingFields, "LINE_BOT_WEBHOOK_PATH")
	}

	// API Key 驗證
	if c.FUGLE_API_KEY == "" {
		missingFields = append(missingFields, "FUGLE_API_KEY")
	}
	if c.IMGBB_API_KEY == "" {
		missingFields = append(missingFields, "IMGBB_API_KEY")
	}
	if c.FINMIND_TOKEN == "" {
		missingFields = append(missingFields, "FINMIND_TOKEN")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("缺少必要的配置項目: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
