package twstock

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade/dto"
	fugleDto "github.com/tian841224/stock-bot/internal/infrastructure/fugle/dto"
	twseDto "github.com/tian841224/stock-bot/internal/infrastructure/twse/dto"
	stockDto "github.com/tian841224/stock-bot/internal/service/twstock/dto"
	"github.com/tian841224/stock-bot/pkg/logger"
	"github.com/tian841224/stock-bot/pkg/utils"

	"go.uber.org/zap"
)

// ========== 股價相關方法 ==========

// GetStockPrice 取得股票價格資訊
func (s *StockService) GetStockPrice(stockID string, date ...string) (*stockDto.StockPriceInfo, error) {
	logger.Log.Info("取得股票價格", zap.String("stockID", stockID))

	// 建立請求參數
	requestDto := dto.FinmindtradeRequestDto{
		DataID: stockID,
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

	latestData := &dto.TaiwanStockPriceData{}
	// 取得最新一筆資料
	if len(response.Data) == 2 {
		latestData = &response.Data[0]
	} else {
		latestData = &response.Data[len(response.Data)-1]
	}

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

	return &stockDto.StockPriceInfo{
		StockID:          stockID,
		StockName:        stockName,
		Date:             latestData.Date,
		OpenPrice:        latestData.Open,
		ClosePrice:       latestData.Close,
		HighPrice:        latestData.Max,
		LowPrice:         latestData.Min,
		Volume:           utils.FormatNumberWithCommas(int64(latestData.TradingVolume)),
		Transaction:      utils.FormatNumberWithCommas(int64(latestData.TradingTurnover)),
		ChangeAmount:     changeAmount,
		PercentageChange: percentageChange,
		UpDownSign:       upDownSign,
	}, nil
}

// GetStockPerformance 取得股票績效
func (s *StockService) GetStockPerformance(stockID string) (*stockDto.StockPerformanceResponseDto, error) {
	logger.Log.Info("取得股票績效", zap.String("stockID", stockID))

	// 取得股票名稱
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	if err != nil || symbol == nil {
		logger.Log.Error("取得股票名稱失敗", zap.Error(err))
		return nil, fmt.Errorf("查無股票名稱")
	}

	// 定義要查詢的期間
	periods := []struct {
		period     string
		periodName string
		startDate  func(now time.Time) time.Time
	}{
		{"YTD", "今年至今", func(now time.Time) time.Time {
			return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		}},
		{"1M", "一個月", func(now time.Time) time.Time {
			return now.AddDate(0, -1, 0)
		}},
		{"6M", "半年", func(now time.Time) time.Time {
			return now.AddDate(0, -6, 0)
		}},
		{"1Y", "一年", func(now time.Time) time.Time {
			return now.AddDate(-1, 0, 0)
		}},
		{"3Y", "三年", func(now time.Time) time.Time {
			return now.AddDate(-3, 0, 0)
		}},
		{"5Y", "五年", func(now time.Time) time.Time {
			return now.AddDate(-5, 0, 0)
		}},
		{"10Y", "十年", func(now time.Time) time.Time {
			return now.AddDate(-10, 0, 0)
		}},
	}

	var performancePeriods []stockDto.StockPerformanceData
	now := time.Now()

	// 取得分割資料
	splitRequestDto := dto.FinmindtradeRequestDto{
		DataID:    stockID,
		StartDate: "1900-01-01",
	}

	splitResponse, err := s.finmindClient.GetTaiwanStockSplitPrice(splitRequestDto)
	if err != nil {
		logger.Log.Error("取得分割資料失敗", zap.Error(err))
		return nil, err
	}

	for _, p := range periods {
		// 計算起始日期
		startDate := p.startDate(now)

		// 判斷是否有分割
		hasSplit := len(splitResponse.Data) > 0

		// 取得起始期間的股價
		startRequestDto := dto.FinmindtradeRequestDto{
			DataID:    stockID,
			StartDate: startDate.Format("2006-01-02"),
		}

		startResponse, err := s.finmindClient.GetTaiwanStockPrice(startRequestDto)
		if err != nil {
			logger.Log.Error("取得起始股價失敗", zap.Error(err))
			continue
		}

		// 檢查是否有資料
		if startResponse.Status != 200 || len(startResponse.Data) == 0 {
			continue
		}

		// 取得第一天和最後一天價格
		startPrice := startResponse.Data[0].Close
		endPrice := startResponse.Data[len(startResponse.Data)-1].Close

		// 計算分割後的股價
		if hasSplit {
			for _, split := range splitResponse.Data {
				// 解析分割日期
				splitDate, err := time.Parse("2006-01-02", split.Date)
				if err != nil {
					logger.Log.Error("解析分割日期失敗", zap.String("date", split.Date), zap.Error(err))
					continue
				}

				// 判斷分割日期是否在起始日期之前
				if splitDate.After(startDate) {
					// 計算分割比例
					splitRatio := split.AfterPrice / split.BeforePrice
					startPrice = startPrice * splitRatio
				}
			}
		}

		// 計算績效
		changeAmount := endPrice - startPrice
		percentageChange := fmt.Sprintf("%.2f%%", (changeAmount/startPrice)*100)

		periodData := stockDto.StockPerformanceData{
			Period:      p.period,
			PeriodName:  p.periodName,
			Performance: percentageChange,
		}

		performancePeriods = append(performancePeriods, periodData)
	}

	if len(performancePeriods) == 0 {
		return nil, fmt.Errorf("查無績效資料")
	}

	return &stockDto.StockPerformanceResponseDto{
		Data: performancePeriods,
	}, nil
}

// GetStockPriceHistory 取得股票每日價格歷史（近5年）
func (s *StockService) GetStockPriceHistory(stockID string) ([]stockDto.StockPerformanceData, error) {
	now := time.Now()
	var performancePeriods []stockDto.StockPerformanceData

	// 取得最近5年的股價資料
	startDate := now.AddDate(-5, 0, 0)
	startRequestDto := dto.FinmindtradeRequestDto{
		DataID:    stockID,
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   now.Format("2006-01-02"),
	}

	priceResponse, err := s.finmindClient.GetTaiwanStockPrice(startRequestDto)
	if err != nil {
		logger.Log.Error("取得股價資料失敗", zap.Error(err))
		return nil, err
	}

	// 檢查是否有資料
	if priceResponse.Status != 200 || len(priceResponse.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料")
	}

	// 取得分割資料
	splitRequestDto := dto.FinmindtradeRequestDto{
		DataID:    stockID,
		StartDate: "1900-01-01",
	}

	splitResponse, err := s.finmindClient.GetTaiwanStockSplitPrice(splitRequestDto)
	if err != nil {
		logger.Log.Error("取得分割資料失敗", zap.Error(err))
		return nil, err
	}

	// 取得基準價格（第一天的收盤價）
	basePrice := priceResponse.Data[0].Close

	// 處理股票分割對基準價格的影響
	if len(splitResponse.Data) > 0 {
		for _, split := range splitResponse.Data {
			splitDate, err := time.Parse("2006-01-02", split.Date)
			if err != nil {
				continue
			}

			// 如果分割日期在基準日期之後，需要調整基準價格
			baseDateParsed, _ := time.Parse("2006-01-02", priceResponse.Data[0].Date)
			if splitDate.After(baseDateParsed) {
				splitRatio := split.AfterPrice / split.BeforePrice
				basePrice = basePrice * splitRatio
			}
		}
	}

	// 每隔幾天取一個點，避免資料點過多
	step := len(priceResponse.Data) / 50 // 最多50個點，適合5年資料
	if step < 1 {
		step = 1
	}
	if step > 50 {
		step = 50 // 限制最大間隔
	}

	// 計算每日相對於基準日的累積漲跌幅
	for i := 0; i < len(priceResponse.Data); i += step {
		priceData := priceResponse.Data[i]
		currentPrice := priceData.Close

		// 計算相對於基準價格的漲跌幅
		changeAmount := currentPrice - basePrice
		percentageChange := (changeAmount / basePrice) * 100

		// 格式化日期
		date, _ := time.Parse("2006-01-02", priceData.Date)
		dateStr := date.Format("2006/01/02")

		periodData := stockDto.StockPerformanceData{
			Period:      priceData.Date,
			PeriodName:  dateStr,
			Performance: fmt.Sprintf("%.2f%%", percentageChange),
		}

		performancePeriods = append(performancePeriods, periodData)
	}

	// 確保包含最後一天的資料
	if len(priceResponse.Data) > 0 {
		lastIndex := len(priceResponse.Data) - 1
		if (lastIndex % step) != 0 {
			priceData := priceResponse.Data[lastIndex]
			currentPrice := priceData.Close

			changeAmount := currentPrice - basePrice
			percentageChange := (changeAmount / basePrice) * 100

			date, _ := time.Parse("2006-01-02", priceData.Date)
			dateStr := date.Format("01/02")

			periodData := stockDto.StockPerformanceData{
				Period:      priceData.Date,
				PeriodName:  dateStr,
				Performance: fmt.Sprintf("%.2f%%", percentageChange),
			}

			performancePeriods = append(performancePeriods, periodData)
		}
	}

	return performancePeriods, nil
}

// GetAfterTradingVolume 取得盤後資訊
func (s *StockService) GetAfterTradingVolume(symbol, date string) (*twseDto.AfterTradingVolumeResponseDto, error) {
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
		if strings.TrimSpace(utils.ToString(row[0])) != strings.TrimSpace(symbol) {
			continue
		}

		openPrice := utils.ToFloat(row[5])
		changeAmount := utils.ToFloat(row[10])
		percentage := utils.PercentageChange(changeAmount, openPrice)

		result := &twseDto.AfterTradingVolumeResponseDto{
			StockId:          utils.ToString(row[0]),
			StockName:        utils.ToString(row[1]),
			Volume:           utils.ToString(row[2]),
			Transaction:      utils.ToString(row[3]),
			Amount:           utils.ToString(row[4]),
			OpenPrice:        openPrice,
			ClosePrice:       utils.ToFloat(row[8]),
			HighPrice:        utils.ToFloat(row[6]),
			LowPrice:         utils.ToFloat(row[7]),
			UpDownSign:       utils.ExtractUpDownSign(utils.ToString(row[9])),
			ChangeAmount:     changeAmount,
			PercentageChange: percentage,
		}
		return result, nil
	}

	return nil, fmt.Errorf("找不到指定股票: %s", symbol)
}

// GetStockNews 取得股票新聞
func (s *StockService) GetStockNews(stockID string) ([]dto.TaiwanNewsResponseData, error) {
	requestDto := dto.FinmindtradeRequestDto{
		DataID:    stockID,
		StartDate: time.Now().Format("2006-01-02"),
	}
	response, err := s.finmindClient.GetTaiwanStockNews(requestDto)
	if err != nil {
		logger.Log.Error("呼叫 FinMind API 失敗", zap.Error(err))
		return nil, err
	}
	if response.Status != 200 {
		return nil, fmt.Errorf("API 回應錯誤: %s", response.Msg)
	}
	return response.Data, nil
}

// GetStockIntradayQuote 取得股票盤中即時資料
func (s *StockService) GetStockIntradayQuote(dto fugleDto.FugleStockQuoteRequestDto) (*fugleDto.FugleStockQuoteResponseDto, error) {
	response, err := s.fugleClient.GetStockIntradayQuote(dto)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetStockHistoricalCandles 取得股票歷史 K 線
func (s *StockService) GetStockHistoricalCandles(dto fugleDto.FugleCandlesRequestDto) (*fugleDto.FugleCandlesResponseDto, error) {
	response, err := s.fugleClient.GetStockHistoricalCandles(dto)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetStockInfo 取得股票詳細資訊
func (s *StockService) GetStockInfo(stockID string) (*stockDto.StockQuoteInfo, error) {
	logger.Log.Info("取得股票詳細資訊", zap.String("stockID", stockID))

	response, err := s.cnyesAPI.GetStockQuote(stockID)
	if err != nil {
		return nil, err
	}

	// 檢查回應狀態
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("API回應錯誤: %s", response.Message)
	}

	// 檢查是否有資料
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料: %s", stockID)
	}

	// 格式化資料（取第一筆）
	stockInfo := s.formatStockInfo(response.Data[0])
	return stockInfo, nil
}

// GetStockQuote 取得股票報價資訊
func (s *StockService) GetStockQuote(stockID string) (*stockDto.StockQuoteInfo, error) {
	// 建構股票符號 (格式: TWS:2330:STOCK)
	symbol := fmt.Sprintf("TWS:%s:STOCK", stockID)

	// 呼叫API
	response, err := s.cnyesAPI.GetStockQuote(symbol)
	if err != nil {
		return nil, fmt.Errorf("取得股票報價失敗: %v", err)
	}

	// 檢查回應
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("API回應錯誤: %s", response.Message)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料: %s", stockID)
	}

	// 格式化資料（取第一筆）
	stockInfo := s.formatStockQuote(response.Data[0])
	return stockInfo, nil
}

// GetTopVolumeItems 取得交易量前20名
func (s *StockService) GetTopVolumeItems() ([]*stockDto.StockPriceInfo, error) {
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

	var result []*stockDto.StockPriceInfo
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

		// 轉換成交量和筆數（格式化為有千分位的字串）
		volumeStr := strings.ReplaceAll(fmt.Sprintf("%v", item[3]), ",", "")
		transactionStr := strings.ReplaceAll(fmt.Sprintf("%v", item[4]), ",", "")

		// 轉換為數字後再格式化為千分位字串
		volumeInt, _ := strconv.ParseInt(volumeStr, 10, 64)
		transactionInt, _ := strconv.ParseInt(transactionStr, 10, 64)

		volume := utils.FormatNumberWithCommas(volumeInt)
		transaction := utils.FormatNumberWithCommas(transactionInt)

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

		stockInfo := &stockDto.StockPriceInfo{
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
