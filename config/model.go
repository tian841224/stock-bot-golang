package config

type Config struct {
	// Database
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	DB_LOG_MODE  bool   `mapstructure:"DB_LOG"`
	// Line Bot
	CHANNEL_ACCESS_TOKEN string `mapstructure:"CHANNEL_ACCESS_TOKEN"`
	CHANNEL_SECRET       string `mapstructure:"CHANNEL_SECRET"`
	// Telegram Bot
	TELEGRAM_ADMIN_CHAT_ID      string `mapstructure:"TELEGRAM_ADMIN_CHAT_ID"`
	TELEGRAM_BOT_TOKEN          string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TELEGRAM_BOT_WEBHOOK_DOMAIN string `mapstructure:"TELEGRAM_BOT_WEBHOOK_DOMAIN"`
	TELEGRAM_BOT_WEBHOOK_PATH   string `mapstructure:"TELEGRAM_BOT_WEBHOOK_PATH"`
	TELEGRAM_BOT_SECRET_TOKEN   string `mapstructure:"TELEGRAM_BOT_SECRET_TOKEN"`
}
