package config

import "strings"

func (c *Config) normalize() {
	c.LINE_BOT_WEBHOOK_PATH = normalizeString(c.LINE_BOT_WEBHOOK_PATH)
	c.TELEGRAM_BOT_SECRET_TOKEN = normalizeString(c.TELEGRAM_BOT_SECRET_TOKEN)
	c.DB_USER = normalizeString(c.DB_USER)
	c.DB_NAME = normalizeString(c.DB_NAME)
	c.CHANNEL_ACCESS_TOKEN = normalizeString(c.CHANNEL_ACCESS_TOKEN)
	c.CHANNEL_SECRET = normalizeString(c.CHANNEL_SECRET)
	c.TELEGRAM_BOT_TOKEN = normalizeString(c.TELEGRAM_BOT_TOKEN)
	c.DB_PASSWORD = normalizeString(c.DB_PASSWORD)
	c.TELEGRAM_BOT_WEBHOOK_DOMAIN = normalizeString(c.TELEGRAM_BOT_WEBHOOK_DOMAIN)
	c.TELEGRAM_BOT_WEBHOOK_PATH = normalizeString(c.TELEGRAM_BOT_WEBHOOK_PATH)
	c.DB_HOST = normalizeString(c.DB_HOST)
	c.FINMIND_TOKEN = normalizeString(c.FINMIND_TOKEN)
	c.FUGLE_API_KEY = normalizeString(c.FUGLE_API_KEY)
	c.IMGBB_API_KEY = normalizeString(c.IMGBB_API_KEY)
}

func normalizeString(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 2 {
		return trimmed
	}

	if (trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'') ||
		(trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"') {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	return trimmed
}
