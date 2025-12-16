// Package config 提供應用程式設定管理
package config

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
