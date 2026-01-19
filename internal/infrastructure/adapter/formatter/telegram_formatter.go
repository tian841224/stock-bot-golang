package formatter

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/pkg/utils"
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

// formatSection 格式化區塊
func (tf *telegramFormatter) formatSection(title string) string {
	return fmt.Sprintf("\n├─ %s\n", title)
}

// formatField 格式化欄位
func (tf *telegramFormatter) formatField(label, value string, emoji string) string {
	if emoji != "" {
		return fmt.Sprintf("│ %s %s: %s\n", emoji, label, value)
	}
	return fmt.Sprintf("│ %s: %s\n", label, value)
}

// formatFieldWithChange 格式化帶漲跌的欄位
func (tf *telegramFormatter) formatFieldWithChange(label, value string, change float64, emoji string) string {
	var changeEmoji string
	var changeText string

	if change > 0 {
		changeEmoji = "📈"
		changeText = fmt.Sprintf("+%.2f", change)
	} else if change < 0 {
		changeEmoji = "📉"
		changeText = fmt.Sprintf("%.2f", change)
	} else {
		changeEmoji = "➖"
		changeText = "0.00"
	}

	return fmt.Sprintf("│ %s %s: %s (%s %s)\n", emoji, label, value, changeEmoji, changeText)
}

// formatPercentage 格式化百分比
func (tf *telegramFormatter) formatPercentage(value float64) string {
	if value > 0 {
		return fmt.Sprintf("📈 +%.2f%%", value)
	} else if value < 0 {
		return fmt.Sprintf("📉 %.2f%%", value)
	}
	return "➖ 0.00%"
}

// formatPriceRange 格式化價格區間
func (tf *telegramFormatter) formatPriceRange(label, high, low string) string {
	return fmt.Sprintf("│ %s: %s ~ %s\n", label, low, high)
}

// formatBidAskPrices 格式化五檔報價
func (tf *telegramFormatter) formatBidAskPrices(bidPrices, askPrices []float64) string {
	var result strings.Builder

	result.WriteString("├─ 📋 五檔報價\n")

	// 賣盤 (由高到低)
	for i := 4; i >= 0; i-- {
		if i < len(askPrices) {
			result.WriteString(fmt.Sprintf("│ 賣%d: %.2f\n", i+1, askPrices[i]))
		}
	}

	result.WriteString("│ ────────\n")

	// 買盤 (由高到低)
	for i := 0; i < 5 && i < len(bidPrices); i++ {
		result.WriteString(fmt.Sprintf("│ 買%d: %.2f\n", i+1, bidPrices[i]))
	}

	return result.String()
}

// formatFooter 格式化頁尾
func (tf *telegramFormatter) formatFooter() string {
	return "╰─────────────────╯"
}

// formatVolume 格式化成交量
func (tf *telegramFormatter) formatVolume(volume int64) string {
	if volume >= 1000000 {
		return fmt.Sprintf("%.1f百萬張", float64(volume)/1000000)
	}
	if volume >= 1000 {
		return fmt.Sprintf("%.1f千張", float64(volume)/1000)
	}
	return fmt.Sprintf("%s張", utils.FormatNumberWithCommas(volume))
}

// formatAmount 格式化金額
func (tf *telegramFormatter) formatAmount(amount float64) string {
	if amount >= 1000000000000 { // 兆
		return fmt.Sprintf("%.2f兆", amount/1000000000000)
	}
	if amount >= 100000000 { // 億
		return fmt.Sprintf("%.2f億", amount/100000000)
	}
	if amount >= 10000 { // 萬
		return fmt.Sprintf("%.2f萬", amount/10000)
	}
	return utils.FormatFloatWithCommas(amount, 2)
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
