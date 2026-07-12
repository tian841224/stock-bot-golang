package config

import "testing"

func TestConfigNormalize(t *testing.T) {
	cfg := Config{
		DB_HOST:                     " postgres ",
		CHANNEL_ACCESS_TOKEN:        "'line-token'",
		TELEGRAM_BOT_TOKEN:          "'123456:abcdef'",
		TELEGRAM_BOT_WEBHOOK_DOMAIN: "\"https://example.com\"",
		TELEGRAM_BOT_WEBHOOK_PATH:   "'/telegram/webhook'",
		FINMIND_TOKEN:               " plain-token ",
	}

	cfg.normalize()

	if cfg.DB_HOST != "postgres" {
		t.Fatalf("expected DB_HOST to be trimmed, got %q", cfg.DB_HOST)
	}
	if cfg.CHANNEL_ACCESS_TOKEN != "line-token" {
		t.Fatalf("expected CHANNEL_ACCESS_TOKEN quotes removed, got %q", cfg.CHANNEL_ACCESS_TOKEN)
	}
	if cfg.TELEGRAM_BOT_TOKEN != "123456:abcdef" {
		t.Fatalf("expected TELEGRAM_BOT_TOKEN quotes removed, got %q", cfg.TELEGRAM_BOT_TOKEN)
	}
	if cfg.TELEGRAM_BOT_WEBHOOK_DOMAIN != "https://example.com" {
		t.Fatalf("expected TELEGRAM_BOT_WEBHOOK_DOMAIN quotes removed, got %q", cfg.TELEGRAM_BOT_WEBHOOK_DOMAIN)
	}
	if cfg.TELEGRAM_BOT_WEBHOOK_PATH != "/telegram/webhook" {
		t.Fatalf("expected TELEGRAM_BOT_WEBHOOK_PATH quotes removed, got %q", cfg.TELEGRAM_BOT_WEBHOOK_PATH)
	}
	if cfg.FINMIND_TOKEN != "plain-token" {
		t.Fatalf("expected FINMIND_TOKEN to be trimmed, got %q", cfg.FINMIND_TOKEN)
	}
}
