package twse

import (
	"fmt"
	twseInfra "stock-bot/internal/infrastructure/twse"
	"stock-bot/internal/infrastructure/twse/dto"
	"strings"
)

type TwseService struct {
	twseAPI *twseInfra.TwseAPI
}

func NewTwseService(twseAPI *twseInfra.TwseAPI) *TwseService {
	return &TwseService{twseAPI: twseAPI}
}

// GetDailyMarketInfo 取得每日市場資訊
func (s *TwseService) GetDailyMarketInfo(count int) ([]dto.DailyMarketInfoData, error) {
	// 參數處理：如果 count 無效值，則使用預設值 1
	if count <= 0 {
		count = 1
	}

	response, err := s.twseAPI.GetDailyMarketInfo()
	if err != nil {
		return nil, err
	}

	// 轉換 Data [][]interface{} 為 DailyMarketInfoData 格式
	var dailyMarketData []dto.DailyMarketInfoData

	// 限制回傳筆數，取最後的 count 筆資料
	data := response.Data
	if count < len(data) {
		data = data[len(data)-count:]
	}

	// 遍歷篩選後的資料
	for _, row := range data {
		if len(row) >= 6 { // 確保有足夠的欄位
			marketInfo := dto.DailyMarketInfoData{
				Date:        s.toString(row[0]), // 日期
				Volume:      s.toString(row[1]), // 成交股數
				Amount:      s.toString(row[2]), // 成交金額
				Transaction: s.toString(row[3]), // 成交筆數
				Index:       s.toString(row[4]), // 發行量加權股價指數
				Change:      s.toString(row[5]), // 漲跌點數
			}
			dailyMarketData = append(dailyMarketData, marketInfo)
		}
	}

	return dailyMarketData, nil
}

// GetAfterTradingVolume 取得盤後資訊
func (s *TwseService) GetAfterTradingVolume(symbol, date string) (*dto.AfterTradingVolumeResponseDto, error) {
	if strings.TrimSpace(symbol) == "" {
		return nil, fmt.Errorf("symbol 為必填參數")
	}

	response, err := s.twseAPI.GetAfterTradingVolume(symbol, date)
	if err != nil {
		return nil, err
	}

	// 檢查資料結構
	if len(response.Tables) <= 8 {
		return nil, fmt.Errorf("查無資料或資料表結構異常")
	}

	stockList := response.Tables[8]
	if len(stockList.Data) == 0 {
		return nil, fmt.Errorf("查無資料")
	}

	// 第 9 個 table 為個股清單，篩選指定股票
	for _, row := range stockList.Data {
		if len(row) < 13 {
			continue
		}
		if strings.TrimSpace(s.toString(row[0])) != strings.TrimSpace(symbol) {
			continue
		}

		openPrice := s.toFloat(row[5])
		changeAmount := s.toFloat(row[10])
		percentage := s.percentageChange(changeAmount, openPrice)

		result := &dto.AfterTradingVolumeResponseDto{
			StockId:          s.toString(row[0]),
			StockName:        s.toString(row[1]),
			Volume:           s.toString(row[2]),
			Transaction:      s.toString(row[3]),
			Amount:           s.toString(row[4]),
			OpenPrice:        openPrice,
			ClosePrice:       s.toFloat(row[8]),
			HighPrice:        s.toFloat(row[6]),
			LowPrice:         s.toFloat(row[7]),
			UpDownSign:       s.extractUpDownSign(s.toString(row[9])),
			ChangeAmount:     changeAmount,
			PercentageChange: percentage,
		}
		return result, nil
	}

	return nil, fmt.Errorf("找不到指定股票: %s", symbol)
}

// GetTopVolumeItems 取得成交量前 20 股票
func (s *TwseService) GetTopVolumeItems() ([]dto.TopVolumeItemsData, error) {
	response, err := s.twseAPI.GetTopVolumeItems()
	if err != nil {
		return nil, err
	}

	// 檢查是否有資料
	if len(response.Data) == 0 {
		return []dto.TopVolumeItemsData{}, nil
	}

	// 將資料轉換為 TopVolumeItemsData 格式
	result := make([]dto.TopVolumeItemsData, 0, len(response.Data))
	for index, item := range response.Data {
		if len(item) < 13 {
			continue
		}

		// 處理數值轉換
		openPrice := s.toFloat(item[5])
		changeAmount := s.toFloat(item[10])

		// 計算漲跌幅
		percentageChange := s.percentageChange(changeAmount, openPrice)

		data := dto.TopVolumeItemsData{
			Rank:             fmt.Sprintf("%d", index+1),               // 排名
			StockId:          s.toString(item[1]),                      // 證券代號
			StockName:        s.toString(item[2]),                      // 證券名稱
			Volume:           s.toString(item[3]),                      // 成交股數
			Transaction:      s.toString(item[4]),                      // 成交筆數
			OpenPrice:        openPrice,                                // 開盤價
			HighPrice:        s.toFloat(item[6]),                       // 最高價
			LowPrice:         s.toFloat(item[7]),                       // 最低價
			ClosePrice:       s.toFloat(item[8]),                       // 收盤價
			UpDownSign:       s.extractUpDownSign(s.toString(item[9])), // 漲跌(+/-)
			ChangeAmount:     changeAmount,                             // 漲跌價差
			PercentageChange: percentageChange,                         // 漲跌幅
			BuyPrice:         s.toFloat(item[11]),                      // 最後揭示買價
			SellPrice:        s.toFloat(item[12]),                      // 最後揭示賣價
		}
		result = append(result, data)
	}
	return result, nil
}

// 輔助函數：將 interface{} 轉換為字串
func (s *TwseService) toString(v interface{}) string {
	str := fmt.Sprint(v)
	str = strings.TrimSpace(str)
	return str
}

// 輔助函數：將 interface{} 轉換為浮點數
func (s *TwseService) toFloat(v interface{}) float64 {
	str := s.toString(v)
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

// 輔助函數：提取漲跌符號
func (s *TwseService) extractUpDownSign(str string) string {
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

// 輔助函數：計算漲跌幅
func (s *TwseService) percentageChange(changeAmount, openPrice float64) string {
	if openPrice == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", (changeAmount/openPrice)*100)
}
