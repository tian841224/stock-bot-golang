// Package utils 提供工具函數
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ToString 將 any 轉換為字串
func ToString(v any) string {
	str := fmt.Sprint(v)
	str = strings.TrimSpace(str)
	return str
}

// ToInt64 將 any 轉換為 int64
func ToInt64(v any) int64 {
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

	// 如果整數部分小於 1000,直接返回
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
