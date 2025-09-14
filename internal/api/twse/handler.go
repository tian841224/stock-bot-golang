package twse

import (
	"fmt"
	twseInfra "stock-bot/internal/infrastructure/twse"
	"stock-bot/internal/infrastructure/twse/dto"
	"strings"

	"github.com/gin-gonic/gin"
)

type TwseHandler struct {
	twseService *twseInfra.TwseAPI
}

func NewTwseHandler(twseService *twseInfra.TwseAPI) *TwseHandler {
	return &TwseHandler{twseService: twseService}
}

func (h *TwseHandler) GetDailyMarketInfo(count *int, c *gin.Context) {
	// 參數處理：如果 count 為 nil 或無效值，則使用預設值 1
	actualCount := 1
	if count != nil && *count > 0 {
		actualCount = *count
	}

	response, err := h.twseService.GetDailyMarketInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 轉換 Data [][]interface{} 為 DailyMarketInfoData 格式
	var dailyMarketData []dto.DailyMarketInfoData

	// 限制回傳筆數，取最後的 actualCount 筆資料
	data := response.Data
	if actualCount < len(data) {
		data = data[len(data)-actualCount:]
	}

	// 遍歷篩選後的資料
	for _, row := range data {
		if len(row) >= 6 { // 確保有足夠的欄位
			marketInfo := dto.DailyMarketInfoData{
				Date:        h.toString(row[0]), // 日期
				Volume:      h.toString(row[1]), // 成交股數
				Amount:      h.toString(row[2]), // 成交金額
				Transaction: h.toString(row[3]), // 成交筆數
				Index:       h.toString(row[4]), // 發行量加權股價指數
				Change:      h.toString(row[5]), // 漲跌點數
			}
			dailyMarketData = append(dailyMarketData, marketInfo)
		}
	}

	c.JSON(200, gin.H{
		"data": dailyMarketData,
	})
}

// GetAfterTradingVolume 盤後資訊
func (h *TwseHandler) GetAfterTradingVolume(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(400, gin.H{"error": "symbol 為必填參數"})
		return
	}
	date := c.Query("date")

	response, err := h.twseService.GetAfterTradingVolume(symbol, date)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 檢查資料結構
	if len(response.Tables) <= 8 {
		c.JSON(404, gin.H{"error": "查無資料或資料表結構異常"})
		return
	}

	stockList := response.Tables[8]
	if len(stockList.Data) == 0 {
		c.JSON(404, gin.H{"error": "查無資料"})
		return
	}

	// 第 9 個 table 為個股清單，篩選指定股票
	for _, row := range stockList.Data {
		if len(row) < 13 {
			continue
		}
		if strings.TrimSpace(h.toString(row[0])) != strings.TrimSpace(symbol) {
			continue
		}

		openPrice := h.toFloat(row[5])
		changeAmount := h.toFloat(row[10])
		percentage := h.percentageChange(changeAmount, openPrice)

		result := dto.AfterTradingVolumeResponseDto{
			StockId:          h.toString(row[0]),
			StockName:        h.toString(row[1]),
			Volume:           h.toString(row[2]),
			Transaction:      h.toString(row[3]),
			Amount:           h.toString(row[4]),
			OpenPrice:        openPrice,
			ClosePrice:       h.toFloat(row[8]),
			HighPrice:        h.toFloat(row[6]),
			LowPrice:         h.toFloat(row[7]),
			UpDownSign:       h.extractUpDownSign(h.toString(row[9])),
			ChangeAmount:     changeAmount,
			PercentageChange: percentage,
		}
		c.JSON(200, result)
		return
	}

	c.JSON(404, gin.H{"error": fmt.Sprintf("找不到指定股票: %s", symbol)})
}

// GetTopVolumeItems 成交量前 20 股票
func (h *TwseHandler) GetTopVolumeItems(c *gin.Context) {
	response, err := h.twseService.GetTopVolumeItems()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 檢查是否有資料
	if len(response.Data) == 0 {
		c.JSON(200, []dto.TopVolumeItemsData{})
		return
	}

	// 將資料轉換為 TopVolumeItemsData 格式
	result := make([]dto.TopVolumeItemsData, 0, len(response.Data))
	for index, item := range response.Data {
		if len(item) < 13 {
			continue
		}

		// 處理數值轉換
		openPrice := h.toFloat(item[5])
		changeAmount := h.toFloat(item[10])

		// 計算漲跌幅
		percentageChange := h.percentageChange(changeAmount, openPrice)

		data := dto.TopVolumeItemsData{
			Rank:             fmt.Sprintf("%d", index+1),               // 排名
			StockId:          h.toString(item[1]),                      // 證券代號
			StockName:        h.toString(item[2]),                      // 證券名稱
			Volume:           h.toString(item[3]),                      // 成交股數
			Transaction:      h.toString(item[4]),                      // 成交筆數
			OpenPrice:        openPrice,                                // 開盤價
			HighPrice:        h.toFloat(item[6]),                       // 最高價
			LowPrice:         h.toFloat(item[7]),                       // 最低價
			ClosePrice:       h.toFloat(item[8]),                       // 收盤價
			UpDownSign:       h.extractUpDownSign(h.toString(item[9])), // 漲跌(+/-)
			ChangeAmount:     changeAmount,                             // 漲跌價差
			PercentageChange: percentageChange,                         // 漲跌幅
			BuyPrice:         h.toFloat(item[11]),                      // 最後揭示買價
			SellPrice:        h.toFloat(item[12]),                      // 最後揭示賣價
		}
		result = append(result, data)
	}
	c.JSON(200, result)
}

// 輔助函數：將 interface{} 轉換為字串
func (h *TwseHandler) toString(v interface{}) string {
	s := fmt.Sprint(v)
	s = strings.TrimSpace(s)
	return s
}

// 輔助函數：將 interface{} 轉換為浮點數
func (h *TwseHandler) toFloat(v interface{}) float64 {
	s := h.toString(v)
	if s == "--" || s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "％", "")
	if s == "+" || s == "-" {
		return 0
	}
	var f float64
	_, err := fmt.Sscan(s, &f)
	if err != nil {
		return 0
	}
	return f
}

// 輔助函數：提取漲跌符號
func (h *TwseHandler) extractUpDownSign(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "+") || strings.Contains(s, "＋") {
		return "+"
	}
	if strings.Contains(s, "-") || strings.Contains(s, "－") {
		return "-"
	}
	return ""
}

// 輔助函數：計算漲跌幅
func (h *TwseHandler) percentageChange(changeAmount, openPrice float64) string {
	if openPrice == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", (changeAmount/openPrice)*100)
}
