package stock

import (
	"fmt"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/infrastructure/finmindtrade/dto"
	"stock-bot/internal/infrastructure/twse"
	"stock-bot/internal/repository"
	"stock-bot/pkg/logger"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type StockService struct {
	finmindClient finmindtrade.FinmindTradeAPIInterface
	twseAPI       *twse.TwseAPI
	symbolsRepo   repository.SymbolRepository
}

func NewStockService(
	finmindClient finmindtrade.FinmindTradeAPIInterface,
	twseAPI *twse.TwseAPI,
	symbolsRepo repository.SymbolRepository,
) *StockService {
	return &StockService{
		finmindClient: finmindClient,
		twseAPI:       twseAPI,
		symbolsRepo:   symbolsRepo,
	}
}

// StockPriceInfo 股票價格資訊
type StockPriceInfo struct {
	StockID          string  `json:"stock_id"`
	StockName        string  `json:"stock_name"`
	Date             string  `json:"date"`
	OpenPrice        float64 `json:"open_price"`
	ClosePrice       float64 `json:"close_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	Volume           int64   `json:"volume"`
	Transaction      int64   `json:"transaction"`
	Amount           int64   `json:"amount"`
	ChangeAmount     float64 `json:"change_amount"`
	PercentageChange string  `json:"percentage_change"`
	UpDownSign       string  `json:"up_down_sign"`
}

// MarketInfo 市場資訊
type MarketInfo struct {
	Date        string `json:"date"`
	Volume      string `json:"volume"`
	Amount      string `json:"amount"`
	Transaction string `json:"transaction"`
	Index       string `json:"index"`
	Change      string `json:"change"`
}

// NewsInfo 新聞資訊
type NewsInfo struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}

// GetStockPrice 取得股票價格資訊
func (s *StockService) GetStockPrice(stockID string, date ...string) (*StockPriceInfo, error) {
	logger.Log.Info("取得股票價格", zap.String("stockID", stockID))

	// 建立請求參數
	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	// 如果有指定日期，設定日期範圍
	if len(date) > 0 && date[0] != "" {
		requestDto.StartDate = date[0]
		requestDto.EndDate = date[0]
	} else {
		// 預設取得最近一天的資料
		today := time.Now().Format("2006-01-02")
		requestDto.StartDate = today
		requestDto.EndDate = today
	}

	// 呼叫 FinMind API
	response, err := s.finmindClient.GetTaiwanStockPrice(requestDto)
	if err != nil {
		logger.Log.Error("呼叫 FinMind API 失敗", zap.Error(err))
		return nil, err
	}

	if response.Status != 200 || len(response.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料")
	}

	// 取得最新一筆資料
	latestData := response.Data[len(response.Data)-1]

	// 取得股票名稱
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	stockName := stockID
	if err == nil && symbol != nil {
		stockName = symbol.Name
	}

	// 計算漲跌幅
	changeAmount := latestData.Close - latestData.Open
	percentageChange := fmt.Sprintf("%.2f%%", (changeAmount/latestData.Open)*100)
	upDownSign := ""
	if changeAmount > 0 {
		upDownSign = "+"
	} else if changeAmount < 0 {
		upDownSign = "-"
		changeAmount = -changeAmount
	}

	return &StockPriceInfo{
		StockID:          stockID,
		StockName:        stockName,
		Date:             latestData.Date,
		OpenPrice:        latestData.Open,
		ClosePrice:       latestData.Close,
		HighPrice:        latestData.Max,
		LowPrice:         latestData.Min,
		Volume:           int64(latestData.TradingVolume),
		Transaction:      int64(latestData.TradingTurnover),
		ChangeAmount:     changeAmount,
		PercentageChange: percentageChange,
		UpDownSign:       upDownSign,
	}, nil
}

// GetDailyMarketInfo 取得大盤資訊
func (s *StockService) GetDailyMarketInfo(count int) ([]*MarketInfo, error) {
	logger.Log.Info("取得大盤資訊", zap.Int("count", count))

	// 呼叫 TWSE API
	response, err := s.twseAPI.GetDailyMarketInfo()
	if err != nil {
		logger.Log.Error("呼叫 TWSE API 失敗", zap.Error(err))
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無市場資料")
	}

	// 取得指定數量的資料
	var result []*MarketInfo
	dataLen := len(response.Data)
	startIdx := 0
	if count < dataLen {
		startIdx = dataLen - count
	}

	for i := startIdx; i < dataLen; i++ {
		data := response.Data[i]
		if len(data) < 6 {
			continue
		}

		marketInfo := &MarketInfo{
			Date:        fmt.Sprintf("%v", data[0]),
			Volume:      fmt.Sprintf("%v", data[1]),
			Amount:      fmt.Sprintf("%v", data[2]),
			Transaction: fmt.Sprintf("%v", data[3]),
			Index:       fmt.Sprintf("%v", data[4]),
			Change:      fmt.Sprintf("%v", data[5]),
		}
		result = append(result, marketInfo)
	}

	return result, nil
}

// GetTopVolumeItems 取得交易量前20名
func (s *StockService) GetTopVolumeItems() ([]*StockPriceInfo, error) {
	logger.Log.Info("取得交易量前20名")

	// 呼叫 TWSE API
	response, err := s.twseAPI.GetTopVolumeItems()
	if err != nil {
		logger.Log.Error("呼叫 TWSE API 失敗", zap.Error(err))
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無交易量資料")
	}

	var result []*StockPriceInfo
	for _, item := range response.Data {
		if len(item) < 12 {
			continue
		}

		// 解析資料 (根據 TWSE API 格式)
		stockID := fmt.Sprintf("%v", item[1])
		stockName := fmt.Sprintf("%v", item[2])

		// 轉換價格
		openPrice, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[5]), 64)
		highPrice, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[6]), 64)
		lowPrice, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[7]), 64)
		closePrice, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[8]), 64)

		// 轉換成交量和筆數
		volume, _ := strconv.ParseInt(fmt.Sprintf("%v", item[3]), 10, 64)
		transaction, _ := strconv.ParseInt(fmt.Sprintf("%v", item[4]), 10, 64)

		// 計算漲跌幅
		changeAmount := closePrice - openPrice
		percentageChange := "0.00%"
		if openPrice != 0 {
			percentageChange = fmt.Sprintf("%.2f%%", (changeAmount/openPrice)*100)
		}

		upDownSign := ""
		if changeAmount > 0 {
			upDownSign = "+"
		} else if changeAmount < 0 {
			upDownSign = "-"
			changeAmount = -changeAmount
		}

		stockInfo := &StockPriceInfo{
			StockID:          stockID,
			StockName:        stockName,
			OpenPrice:        openPrice,
			ClosePrice:       closePrice,
			HighPrice:        highPrice,
			LowPrice:         lowPrice,
			Volume:           volume,
			Transaction:      transaction,
			ChangeAmount:     changeAmount,
			PercentageChange: percentageChange,
			UpDownSign:       upDownSign,
		}
		result = append(result, stockInfo)
	}

	return result, nil
}

// GetStockAnalysis 取得股票分析圖表
func (s *StockService) GetStockAnalysis(stockID string) ([]byte, string, error) {
	logger.Log.Info("取得股票分析", zap.String("stockID", stockID))

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	// 呼叫 FinMind API
	response, err := s.finmindClient.GetTaiwanStockAnalysisPlot(requestDto)
	if err != nil {
		logger.Log.Error("呼叫 FinMind API 失敗", zap.Error(err))
		return nil, "", err
	}

	if response.Status != 200 {
		return nil, "", fmt.Errorf("API 回應錯誤: %s", response.Msg)
	}

	// 取得股票名稱
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	stockName := stockID
	if err == nil && symbol != nil {
		stockName = symbol.Name
	}

	// 由於 FinMind API 的分析圖表回應不包含圖片數據，暫時返回空數據
	// 實際使用時需要根據 API 文檔調整
	return []byte{}, stockName, nil
}

// ValidateStockID 驗證股票代號是否存在
func (s *StockService) ValidateStockID(stockID string) (bool, string, error) {
	// 先從本地資料庫查詢
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	if err == nil && symbol != nil {
		return true, symbol.Name, nil
	}

	// 如果本地找不到，嘗試從 API 查詢
	_, err = s.GetStockPrice(stockID)
	if err != nil {
		return false, "", err
	}

	return true, stockID, nil
}
