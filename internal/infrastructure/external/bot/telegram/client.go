// Package tgbot 提供 Telegram Bot 客戶端實作
package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type TgBotClient struct {
	Client *tgbotapi.BotAPI
	logger logger.Logger
}

// NewBot 初始化 Telegram Bot 並設定 webhook
func NewBot(cfg config.Config, log logger.Logger) (*TgBotClient, error) {
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
	return &TgBotClient{Client: client, logger: log}, nil
}

// SendMessage 發送訊息
func (c *TgBotClient) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.Client.Send(msg)
	if err != nil {
		c.logger.Error("發送訊息失敗", logger.Error(err))
	}
	return err
}

// SendMessageWithKeyboard 發送帶有鍵盤的訊息
func (c *TgBotClient) SendMessageWithKeyboard(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, err := c.Client.Send(msg)
	if err != nil {
		c.logger.Error("發送帶有鍵盤的訊息失敗", logger.Error(err))
	}
	return err
}

// SendMessageHTML 發送 HTML 訊息
func (c *TgBotClient) SendMessageHTML(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.Client.Send(msg)
	if err != nil {
		c.logger.Error("發送 HTML 訊息失敗", logger.Error(err))
	}
	return err
}

// SendPhoto 發送圖片
func (c *TgBotClient) SendPhoto(chatID int64, data []byte, caption string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  "chart.png",
		Bytes: data,
	})
	photo.Caption = caption
	photo.ParseMode = tgbotapi.ModeHTML
	_, err := c.Client.Send(photo)
	if err != nil {
		c.logger.Error("發送圖片失敗", logger.Error(err))
	}
	return err
}
