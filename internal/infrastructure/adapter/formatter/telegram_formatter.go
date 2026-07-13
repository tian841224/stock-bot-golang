package formatter

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tian841224/stock-bot/internal/application/dto"
)

type TelegramFormatter interface {
	FormatStockNews(news []dto.StockNews, stockName, symbol string) *dto.TgStockNewsMessage
}

// TelegramFormatter Telegram 訊息格式化器
type telegramFormatter struct {
}

// NewTelegramFormatter 建立新的 Telegram 格式化器
func NewTelegramFormatter() *telegramFormatter {
	return &telegramFormatter{}
}

// FormatStockNews 格式化 Telegram 股票新聞訊息（包含按鈕）
func (tf *telegramFormatter) FormatStockNews(news []dto.StockNews, stockName, symbol string) *dto.TgStockNewsMessage {
	if len(news) == 0 {
		return &dto.TgStockNewsMessage{
			Text: fmt.Sprintf("⚡️%s(%s)-即時新聞\n\n暫無新聞資料", stockName, symbol),
		}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, n := range news {
		btn := tgbotapi.NewInlineKeyboardButtonURL(n.Title, n.Link)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &dto.TgStockNewsMessage{
		Text:                 fmt.Sprintf("⚡️%s(%s)-即時新聞", stockName, symbol),
		InlineKeyboardMarkup: &keyboard,
	}
}
