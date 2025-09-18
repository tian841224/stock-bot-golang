package utils

import (
	"fmt"
	"strings"
)

// TelegramFormatter Telegram 訊息格式化器
type TelegramFormatter struct{}

// NewTelegramFormatter 建立新的 Telegram 格式化器
func NewTelegramFormatter() *TelegramFormatter {
	return &TelegramFormatter{}
}

// FormatStockInfo 格式化股票資訊為 Telegram 訊息
func (tf *TelegramFormatter) FormatStockInfo(stockInfo interface{}) string {
	// 這裡使用 interface{} 是為了避免循環依賴
	// 在實際使用時會傳入 StockQuoteInfo 結構
	return tf.buildStockMessage(stockInfo)
}

// buildStockMessage 建構股票訊息
func (tf *TelegramFormatter) buildStockMessage(data interface{}) string {
	var message strings.Builder

	// 這裡需要使用反射或類型斷言來處理
	// 暫時先建立一個通用的格式化函數
	message.WriteString(tf.formatHeader("📊 股票資訊"))
	message.WriteString("\n")

	return message.String()
}

// formatHeader 格式化標題
func (tf *TelegramFormatter) formatHeader(title string) string {
	return fmt.Sprintf("╭─ %s ─╮\n", title)
}

// formatSection 格式化區塊
func (tf *TelegramFormatter) formatSection(title string) string {
	return fmt.Sprintf("\n├─ %s\n", title)
}

// formatField 格式化欄位
func (tf *TelegramFormatter) formatField(label, value string, emoji string) string {
	if emoji != "" {
		return fmt.Sprintf("│ %s %s: %s\n", emoji, label, value)
	}
	return fmt.Sprintf("│ %s: %s\n", label, value)
}

// formatFieldWithChange 格式化帶漲跌的欄位
func (tf *TelegramFormatter) formatFieldWithChange(label, value string, change float64, emoji string) string {
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
func (tf *TelegramFormatter) formatPercentage(value float64) string {
	if value > 0 {
		return fmt.Sprintf("📈 +%.2f%%", value)
	} else if value < 0 {
		return fmt.Sprintf("📉 %.2f%%", value)
	}
	return "➖ 0.00%"
}

// formatPriceRange 格式化價格區間
func (tf *TelegramFormatter) formatPriceRange(label, high, low string) string {
	return fmt.Sprintf("│ %s: %s ~ %s\n", label, low, high)
}

// formatBidAskPrices 格式化五檔報價
func (tf *TelegramFormatter) formatBidAskPrices(bidPrices, askPrices []float64) string {
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
func (tf *TelegramFormatter) formatFooter() string {
	return "╰─────────────────╯"
}

// formatVolume 格式化成交量
func (tf *TelegramFormatter) formatVolume(volume int64) string {
	if volume >= 1000000 {
		return fmt.Sprintf("%.1f百萬張", float64(volume)/1000000)
	}
	if volume >= 1000 {
		return fmt.Sprintf("%.1f千張", float64(volume)/1000)
	}
	return fmt.Sprintf("%s張", FormatNumberWithCommas(volume))
}

// formatAmount 格式化金額
func (tf *TelegramFormatter) formatAmount(amount float64) string {
	if amount >= 1000000000000 { // 兆
		return fmt.Sprintf("%.2f兆", amount/1000000000000)
	}
	if amount >= 100000000 { // 億
		return fmt.Sprintf("%.2f億", amount/100000000)
	}
	if amount >= 10000 { // 萬
		return fmt.Sprintf("%.2f萬", amount/10000)
	}
	return FormatFloatWithCommas(amount, 2)
}

// EscapeMarkdown 跳脫 Markdown 特殊字符
func (tf *TelegramFormatter) EscapeMarkdown(text string) string {
	// Telegram MarkdownV2 需要跳脫的字符
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}

// FormatBold 格式化粗體文字
func (tf *TelegramFormatter) FormatBold(text string) string {
	return fmt.Sprintf("*%s*", tf.EscapeMarkdown(text))
}

// FormatItalic 格式化斜體文字
func (tf *TelegramFormatter) FormatItalic(text string) string {
	return fmt.Sprintf("_%s_", tf.EscapeMarkdown(text))
}

// FormatCode 格式化程式碼文字
func (tf *TelegramFormatter) FormatCode(text string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(text, "`", "\\`"))
}

// FormatCodeBlock 格式化程式碼區塊
func (tf *TelegramFormatter) FormatCodeBlock(text string) string {
	return fmt.Sprintf("```\n%s\n```", text)
}
