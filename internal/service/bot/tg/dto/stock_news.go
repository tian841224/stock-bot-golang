// Package dto 提供 Telegram Bot 的資料傳輸物件
package dto

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type StockNewsMessage struct {
	Text                 string
	InlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
}
