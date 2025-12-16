package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/pkg/utils"
)

// FormatCurrency 格式化貨幣顯示
func FormatCurrency(amount float64, currency string) string {
	formatted := utils.FormatFloatWithCommas(amount, 2)
	return fmt.Sprintf("%s %s", formatted, currency)
}

// FormatVolume 格式化成交量顯示（自動選擇單位）
func FormatVolume(volume int64) string {
	if volume >= 1000000 {
		return fmt.Sprintf("%.1f百萬", float64(volume)/1000000)
	}
	if volume >= 1000 {
		return fmt.Sprintf("%.1f千", float64(volume)/1000)
	}
	return utils.FormatNumberWithCommas(volume)
}

// FormatAmount 格式化金額顯示（自動選擇單位）
func FormatAmount(amount float64) string {
	if amount >= 1000000000000 { // 兆
		return fmt.Sprintf("%.2f 兆", amount/1000000000000)
	}
	if amount >= 100000000 { // 億
		return fmt.Sprintf("%.2f 億", amount/100000000)
	}
	if amount >= 10000 { // 萬
		return fmt.Sprintf("%.2f 萬", amount/10000)
	}
	return utils.FormatFloatWithCommas(amount, 2)
}

// FormatAmount 格式化金額顯示（自動選擇單位）
func FormatAmountInt(amount int64) string {
	if amount >= 1000000000000 { // 兆
		return fmt.Sprintf("%.2f 兆", float64(amount)/1000000000000)
	}
	if amount >= 100000000 { // 億
		return fmt.Sprintf("%.2f 億", float64(amount)/100000000)
	}
	if amount >= 10000 { // 萬
		return fmt.Sprintf("%.2f 萬", float64(amount)/10000)
	}
	return utils.FormatNumberWithCommas(amount)
}

// 格式化字串成int64
func FormatStringToInt64(str string) int64 {
	return utils.ToInt64(str)
}

// TruncateString 截斷字串並加上省略號
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// PadString 填充字串到指定長度
func PadString(str string, length int, padChar rune) string {
	if len(str) >= length {
		return str
	}
	padding := strings.Repeat(string(padChar), length-len(str))
	return str + padding
}

// CleanString 清理字串（移除多餘空白和特殊字符）
func CleanString(str string) string {
	// 移除前後空白
	str = strings.TrimSpace(str)

	// 移除多餘的空白字符
	str = strings.ReplaceAll(str, "\t", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")

	// 移除連續的空格
	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}

	return str
}

// FormatPrice 格式化價格
func FormatPrice(price float64) string {
	if price >= 1000 {
		return fmt.Sprintf("%.0f", price)
	} else if price >= 100 {
		return fmt.Sprintf("%.1f", price)
	} else if price >= 10 {
		return fmt.Sprintf("%.2f", price)
	} else {
		return fmt.Sprintf("%.3f", price)
	}
}

// FormatPercentage 格式化百分比
func FormatPercentage(value float64) string {
	if value > 0 {
		return fmt.Sprintf("+%.2f%%", value)
	} else if value < 0 {
		return fmt.Sprintf("%.2f%%", value)
	}
	return "0.00%"
}

// FormatTimeFromTimestamp 將時間戳記格式化為 YYYY/MM 格式
func FormatTimeFromTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006/01")
}
