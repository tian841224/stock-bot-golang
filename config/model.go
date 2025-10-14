package config

type Config struct {
	// Database
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     int    `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
	DB_LOG_MODE bool   `mapstructure:"DB_LOG"`
	// Line Bot
	CHANNEL_ACCESS_TOKEN  string `mapstructure:"CHANNEL_ACCESS_TOKEN"`
	CHANNEL_SECRET        string `mapstructure:"CHANNEL_SECRET"`
	LINE_BOT_WEBHOOK_PATH string `mapstructure:"LINE_BOT_WEBHOOK_PATH"`
	// Telegram Bot
	TELEGRAM_ADMIN_CHAT_ID      string `mapstructure:"TELEGRAM_ADMIN_CHAT_ID"`
	TELEGRAM_BOT_TOKEN          string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TELEGRAM_BOT_WEBHOOK_DOMAIN string `mapstructure:"TELEGRAM_BOT_WEBHOOK_DOMAIN"`
	TELEGRAM_BOT_WEBHOOK_PATH   string `mapstructure:"TELEGRAM_BOT_WEBHOOK_PATH"`
	TELEGRAM_BOT_SECRET_TOKEN   string `mapstructure:"TELEGRAM_BOT_SECRET_TOKEN"`
	// Finmind Trade
	FINMIND_TOKEN string `mapstructure:"FINMIND_TOKEN"`
	// Fugle
	FUGLE_API_KEY string `mapstructure:"FUGLE_API_KEY"`
	// Imgbb
	IMGBB_API_KEY string `mapstructure:"IMGBB_API_KEY"`
}
