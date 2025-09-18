package utils

import (
	"fmt"
	"strings"
)

// FormatCurrency 格式化貨幣顯示
func FormatCurrency(amount float64, currency string) string {
	formatted := FormatFloatWithCommas(amount, 2)
	return fmt.Sprintf("%s %s", formatted, currency)
}

// FormatPercentage 格式化百分比顯示
func FormatPercentage(value float64, precision int) string {
	return fmt.Sprintf("%."+fmt.Sprintf("%d", precision)+"f%%", value)
}

// FormatVolume 格式化成交量顯示（自動選擇單位）
func FormatVolume(volume int64) string {
	if volume >= 1000000 {
		return fmt.Sprintf("%.1f百萬", float64(volume)/1000000)
	}
	if volume >= 1000 {
		return fmt.Sprintf("%.1f千", float64(volume)/1000)
	}
	return FormatNumberWithCommas(volume)
}

// FormatAmount 格式化金額顯示（自動選擇單位）
func FormatAmount(amount float64) string {
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
