package twse

import (
	twstockService "stock-bot/internal/service/twstock"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TwseHandler struct {
	twstockService *twstockService.StockService
}

func NewTwseHandler(twstockService *twstockService.StockService) *TwseHandler {
	return &TwseHandler{twstockService: twstockService}
}

func (h *TwseHandler) GetDailyMarketInfo(count *int, c *gin.Context) {
	// 參數處理：如果 count 為 nil 或無效值，則使用預設值 1
	actualCount := 1
	if count != nil && *count > 0 {
		actualCount = *count
	}

	// 也可以從查詢參數獲取 count
	if countParam := c.Query("count"); countParam != "" {
		if parsedCount, err := strconv.Atoi(countParam); err == nil && parsedCount > 0 {
			actualCount = parsedCount
		}
	}

	dailyMarketData, err := h.twstockService.GetDailyMarketInfo(actualCount)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": dailyMarketData,
	})
}

// GetAfterTradingVolume 盤後資訊
func (h *TwseHandler) GetAfterTradingVolume(c *gin.Context) {
	symbol := c.Query("symbol")
	date := c.Query("date")

	result, err := h.twstockService.GetAfterTradingVolume(symbol, date)
	if err != nil {
		// 根據錯誤類型回傳不同的 HTTP 狀態碼
		if err.Error() == "symbol 為必填參數" {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "查無資料" || err.Error() == "查無資料或資料表結構異常" {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// // GetTopVolumeItems 成交量前 20 股票
// func (h *TwseHandler) GetTopVolumeItems(c *gin.Context) {
// 	result, err := h.twstockService.GetTopVolumeItems()
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, result)
// }
