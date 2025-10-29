// Package tgbot 提供 Telegram Bot 客戶端實作
package tgbot

import (
	"github.com/tian841224/stock-bot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotClient struct {
	Client *tgbotapi.BotAPI
}

// NewBot 初始化 Telegram Bot 並設定 webhook
func NewBot(cfg config.Config) (*TgBotClient, error) {
	client, err := tgbotapi.NewBotAPI(cfg.TELEGRAM_BOT_TOKEN)
	if err != nil {
		return nil, err
	}

	// 設定 webhook
	webhookURL := cfg.TELEGRAM_BOT_WEBHOOK_DOMAIN + cfg.TELEGRAM_BOT_WEBHOOK_PATH
	if cfg.TELEGRAM_BOT_SECRET_TOKEN != "" {
		params := tgbotapi.Params{}
		params["url"] = webhookURL
		params["secret_token"] = cfg.TELEGRAM_BOT_SECRET_TOKEN
		// 捨棄待處理的舊更新
		params["drop_pending_updates"] = "true"
		if _, err := client.MakeRequest("setWebhook", params); err != nil {
			return nil, err
		}
	} else {
		webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
		if err != nil {
			return nil, err
		}
		if _, err = client.Request(webhookConfig); err != nil {
			return nil, err
		}
	}
	return &TgBotClient{Client: client}, nil
}
