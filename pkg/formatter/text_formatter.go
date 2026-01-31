package formatter

import (
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/pkg/utils"
)

// FormatAmountInt 格式化金額顯示(自動選擇單位)）
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

// FormatTimeFromTimestamp 將時間戳記格式化為 YYYY/MM 格式
func FormatTimeFromTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006/01")
}
