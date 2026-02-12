package formatter

import (
	"fmt"
	"strings"

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

// FormatStockInfo 格式化股票資訊為 Telegram 訊息
func (tf *telegramFormatter) FormatStockInfo(stockInfo interface{}) string {
	return tf.buildStockMessage()
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

// buildStockMessage 建構股票訊息
func (tf *telegramFormatter) buildStockMessage() string {
	var message strings.Builder

	message.WriteString(tf.formatHeader("📊 股票資訊"))
	message.WriteString("\n")

	return message.String()
}

// formatHeader 格式化標題
func (tf *telegramFormatter) formatHeader(title string) string {
	return fmt.Sprintf("╭─ %s ─╮\n", title)
}

// EscapeMarkdown 跳脫 Markdown 特殊字符
func (tf *telegramFormatter) EscapeMarkdown(text string) string {
	// Telegram MarkdownV2 需要跳脫的字符
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}

// FormatBold 格式化粗體文字
func (tf *telegramFormatter) FormatBold(text string) string {
	return fmt.Sprintf("*%s*", tf.EscapeMarkdown(text))
}

// FormatItalic 格式化斜體文字
func (tf *telegramFormatter) FormatItalic(text string) string {
	return fmt.Sprintf("_%s_", tf.EscapeMarkdown(text))
}

// FormatCode 格式化程式碼文字
func (tf *telegramFormatter) FormatCode(text string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(text, "`", "\\`"))
}

// FormatCodeBlock 格式化程式碼區塊
func (tf *telegramFormatter) FormatCodeBlock(text string) string {
	return fmt.Sprintf("```\n%s\n```", text)
}
