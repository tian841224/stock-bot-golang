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
