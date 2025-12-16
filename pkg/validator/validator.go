package utils

import (
	"regexp"
	"strings"

	utils "github.com/tian841224/stock-bot/pkg/utils"
)

// IsValidStockID 驗證股票代碼是否有效
func IsValidStockID(stockID string) bool {
	if stockID == "" {
		return false
	}

	// 台股代碼通常是4位數字，或是特殊ETF代碼
	matched, _ := regexp.MatchString(`^[0-9]{4}$`, stockID)
	return matched
}

// IsValidDate 驗證日期格式是否正確 (YYYY-MM-DD)
func IsValidDate(date string) bool {
	if date == "" {
		return false
	}

	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	return matched
}

// IsEmpty 檢查字串是否為空或只包含空白字符
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNumeric 檢查字串是否為數字
func IsNumeric(str string) bool {
	if str == "" {
		return false
	}

	// 移除千分位逗號
	str = strings.ReplaceAll(str, ",", "")

	// 檢查是否為數字（包含小數點）
	matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, str)
	return matched
}

// IsPositiveNumber 檢查字串是否為正數
func IsPositiveNumber(str string) bool {
	if !IsNumeric(str) {
		return false
	}

	value := ToFloat(str)
	return value > 0
}

// SanitizeString 清理字串中的危險字符
func SanitizeString(str string) string {
	// 移除HTML標籤
	re := regexp.MustCompile(`<[^>]*>`)
	str = re.ReplaceAllString(str, "")

	// 移除SQL注入相關字符
	dangerousChars := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_"}
	for _, char := range dangerousChars {
		str = strings.ReplaceAll(str, char, "")
	}

	return CleanString(str)
}


func CleanString(str string) string {
	// 移除前後空白
	str = strings.TrimSpace(str)

	// 移除多餘的空白字符
	str = strings.ReplaceAll(str, "\t", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")

	return str
}

func ToFloat(str string) float64 {
	return utils.ToFloat(str)
}