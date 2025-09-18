package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ToString 將 interface{} 轉換為字串
func ToString(v interface{}) string {
	str := fmt.Sprint(v)
	str = strings.TrimSpace(str)
	return str
}

// ToFloat 將 interface{} 轉換為浮點數
func ToFloat(v interface{}) float64 {
	str := ToString(v)
	if str == "--" || str == "" {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "％", "")
	if str == "+" || str == "-" {
		return 0
	}
	var f float64
	_, err := fmt.Sscan(str, &f)
	if err != nil {
		return 0
	}
	return f
}

// ToInt 將 interface{} 轉換為整數
func ToInt(v interface{}) int {
	str := ToString(v)
	if str == "--" || str == "" {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// ToInt64 將 interface{} 轉換為 int64
func ToInt64(v interface{}) int64 {
	str := ToString(v)
	if str == "--" || str == "" {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// ExtractUpDownSign 提取漲跌符號
func ExtractUpDownSign(str string) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}
	if strings.Contains(str, "+") || strings.Contains(str, "＋") {
		return "+"
	}
	if strings.Contains(str, "-") || strings.Contains(str, "－") {
		return "-"
	}
	return ""
}

// PercentageChange 計算漲跌幅
func PercentageChange(changeAmount, openPrice float64) string {
	if openPrice == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", (changeAmount/openPrice)*100)
}

// FormatNumberWithCommas 將數字格式化為千分位字串
func FormatNumberWithCommas(num int64) string {
	str := strconv.FormatInt(num, 10)

	// 如果數字小於 1000，直接返回
	if len(str) <= 3 {
		return str
	}

	// 從右邊開始，每三位加一個逗號
	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}

	return result
}

// FormatFloatWithCommas 將浮點數格式化為千分位字串
func FormatFloatWithCommas(num float64, precision int) string {
	str := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)

	// 分離整數和小數部分
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	// 如果整數部分小於 1000，直接返回
	if len(intPart) <= 3 {
		return str
	}

	// 格式化整數部分
	result := ""
	for i, char := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}

	return result + decPart
}
