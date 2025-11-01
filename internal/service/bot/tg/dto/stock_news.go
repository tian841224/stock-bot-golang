package dto

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// StockNewsMessage 股票新聞訊息結構
type StockNewsMessage struct {
	InlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
	Text                 string
}
