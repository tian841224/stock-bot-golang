package stock

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	finmindtradeDto "github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade/dto"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/twse"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	"github.com/tian841224/stock-bot/pkg/utils"
)

type marketDataGateway struct {
	marketDataPort  port.MarketDataPort
	twseAPI         *twse.TwseAPI
	cnyesAPI        *cnyes.CnyesAPI
	fugleAPI        *fugle.FugleAPI
	finmindAPI      *finmindtrade.FinmindTradeAPI
	validationPort  port.ValidationPort
	logger          logger.Logger
	tradeDateReader port.TradeDateReader
}

func NewMarketDataGateway(marketDataPort port.MarketDataPort, twseAPI *twse.TwseAPI, cnyesAPI *cnyes.CnyesAPI, fugleAPI *fugle.FugleAPI, finmindAPI *finmindtrade.FinmindTradeAPI, validationPort port.ValidationPort, tradeDateReader port.TradeDateReader) *marketDataGateway {
	return &marketDataGateway{
		marketDataPort:  marketDataPort,
		twseAPI:         twseAPI,
		cnyesAPI:        cnyesAPI,
		fugleAPI:        fugleAPI,
		finmindAPI:      finmindAPI,
		validationPort:  validationPort,
		tradeDateReader: tradeDateReader,
	}
}

var _ port.MarketDataPort = (*marketDataGateway)(nil)

// 取得大盤每日成交資訊
func (m *marketDataGateway) GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {

	response, err := m.twseAPI.GetDailyMarketInfo()
	if err != nil {
		m.logger.Error("呼叫 TWSE API 失敗", logger.Error(err))
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無資料，請確認後再試。")
	}

	// 如果指定了筆數且小於總資料數，則從最後開始取指定筆數
	if count > 0 && count < len(response.Data) {
		// 取最後的 count 筆資料（從陣列末尾開始）
		startIndex := len(response.Data) - count
		response.Data = response.Data[startIndex:]
	}

	responseDto := make([]dto.DailyMarketInfo, len(response.Data))
	for i, data := range response.Data {
		responseDto[i] = dto.DailyMarketInfo{
			Date:        data[0],
			Volume:      data[1],
			Amount:      data[2],
			Transaction: data[3],
			Index:       data[4],
			Change:      data[5],
		}
	}

	return &responseDto, nil
}

func (m *marketDataGateway) GetTopVolumeStock(ctx context.Context) ([]*dto.TopVolume, error) {
	// 呼叫 TWSE API
	response, err := m.twseAPI.GetTopVolumeItems()
	if err != nil {
		m.logger.Error("呼叫 TWSE API 失敗", logger.Error(err))
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無交易量資料")
	}

	var result []*dto.TopVolume
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

		stockInfo := &dto.TopVolume{
			StockSymbol:      stockID,
			StockName:        stockName,
			Date:             response.Date,
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

func (m *marketDataGateway) GetStockPrice(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error) {
	stock, err := m.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		m.logger.Error("驗證股票代號失敗", logger.Error(err))
		return nil, err
	}
	if stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 建立請求參數
	requestDto := finmindtradeDto.FinmindtradeRequestDto{
		DataID: symbol,
	}

	// 如果只有一個日期，則取該日期、若有兩個日期則取區間日期、若沒有日期則取今天
	if len(dates) == 1 {
		requestDto.StartDate = dates[0].Format("2006-01-02")
		requestDto.EndDate = dates[0].Format("2006-01-02")
	} else if len(dates) == 2 {
		requestDto.StartDate = dates[0].Format("2006-01-02")
		requestDto.EndDate = dates[1].Format("2006-01-02")
	} else {
		requestDto.StartDate = time.Now().Format("2006-01-02")
		requestDto.EndDate = time.Now().Format("2006-01-02")
	}

	// 呼叫 FinMind API
	response, err := m.finmindAPI.GetTaiwanStockPrice(requestDto)
	if err != nil {
		m.logger.Error("呼叫 FinMind API 失敗", logger.Error(err))
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無資料")
	}

	result := make([]dto.StockPrice, len(response.Data))
	for i, data := range response.Data {
		parsedDate, err := time.Parse("2006-01-02", data.Date)
		if err != nil {
			m.logger.Error("解析日期失敗", logger.String("date", data.Date), logger.Error(err))
			return nil, fmt.Errorf("解析日期失敗: %w", err)
		}
		result[i] = dto.StockPrice{
			Symbol:       stock.Symbol,
			Name:         stock.Name,
			Date:         parsedDate,
			OpenPrice:    data.Open,
			ClosePrice:   data.Close,
			HighPrice:    data.Max,
			LowPrice:     data.Min,
			Volume:       int64(data.TradingVolume),
			Transactions: int64(data.TradingTurnover),
			Amount:       float64(data.TradingTurnover),
		}
	}
	return &result, nil
}

// 取得股票近五年價格歷史資料
func (m *marketDataGateway) GetStockPerformance(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error) {
	result := make([]dto.StockPerformanceData, 0)
	stock, err := m.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		m.logger.Error("驗證股票代號失敗", logger.Error(err))
		return nil, err
	}
	if stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 建立請求參數
	requestDto := finmindtradeDto.FinmindtradeRequestDto{
		DataID:    symbol,
		StartDate: time.Now().AddDate(-5, 0, 0).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
	}

	// 呼叫 FinMind API
	priceResponse, err := m.finmindAPI.GetTaiwanStockPrice(requestDto)
	if err != nil {
		m.logger.Error("呼叫 FinMind API 失敗", logger.Error(err))
		return nil, err
	}

	// 檢查是否有資料
	if priceResponse.Status != 200 || len(priceResponse.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料")
	}

	// 取得分割資料
	splitRequestDto := finmindtradeDto.FinmindtradeRequestDto{
		DataID:    symbol,
		StartDate: "1900-01-01",
	}

	splitResponse, err := m.finmindAPI.GetTaiwanStockSplitPrice(splitRequestDto)
	if err != nil {
		m.logger.Error("取得分割資料失敗", logger.Error(err))
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

		periodData := dto.StockPerformanceData{
			Period:      priceData.Date,
			PeriodName:  dateStr,
			Performance: fmt.Sprintf("%.2f%%", percentageChange),
		}

		result = append(result, periodData)
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

			periodData := dto.StockPerformanceData{
				Period:      priceData.Date,
				PeriodName:  dateStr,
				Performance: fmt.Sprintf("%.2f%%", percentageChange),
			}

			result = append(result, periodData)
		}
	}

	return result, nil
}

func (m *marketDataGateway) GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error) {
	stock, err := m.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		m.logger.Error("驗證股票代號失敗", logger.Error(err))
		return nil, err
	}
	if stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	response, err := m.cnyesAPI.GetStockQuote(stock.Symbol)
	if err != nil {
		return nil, err
	}

	// 檢查是否有資料
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料: %s", stock.Symbol)
	}

	// 格式化資料（取第一筆）
	stockInfo := dto.StockCompanyInfo{
		Symbol:       response.Data[0].Symbol,
		Name:         response.Data[0].StockName,
		Industry:     response.Data[0].Industry,
		Market:       response.Data[0].Market,
		PE:           response.Data[0].PE,
		PB:           response.Data[0].PB,
		MarketCap:    response.Data[0].MarketCap,
		BookValue:    response.Data[0].BookValue,
		EPS:          response.Data[0].EPS,
		QuarterEPS:   response.Data[0].QuarterEPS,
		Dividend:     response.Data[0].Dividend,
		DividendRate: response.Data[0].DividendRate,
		GrossMargin:  response.Data[0].GrossMargin,
		OperMargin:   response.Data[0].OperMargin,
		NetMargin:    response.Data[0].NetMargin,
	}
	return &stockInfo, nil
}

func (m *marketDataGateway) GetStockRevenue(ctx context.Context, symbol string) (*dto.StockRevenue, error) {
	stock, err := m.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		m.logger.Error("驗證股票代號失敗", logger.Error(err))
		return nil, err
	}
	if stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 取得近12個月財報
	response, err := m.cnyesAPI.GetRevenue(stock.Symbol, 12)
	if err != nil {
		return nil, err
	}

	return &dto.StockRevenue{
		StockSymbol:     stock.Symbol,
		StockName:       stock.Name,
		Time:            response.Data.Time,
		StockPrice:      response.Data.Datasets.C,
		SaleMonth:       response.Data.Datasets.SaleMonth,
		SaleAccumulated: response.Data.Datasets.SaleAccumulated,
		YoY:             response.Data.Datasets.YoY,
		YoYAccumulated:  response.Data.Datasets.YoYAccumulated,
	}, nil
}

func (m *marketDataGateway) GetLatestTradeDate(ctx context.Context) (time.Time, error) {
	tradeDate, err := m.tradeDateReader.GetByDateRange(ctx, time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return time.Time{}, err
	}

	if len(tradeDate) == 0 {
		return time.Time{}, fmt.Errorf("找不到交易日資料")
	}

	now := time.Now()
	// 從最新的日期開始往回找
	for i := len(tradeDate) - 1; i >= 0; i-- {
		date := tradeDate[i]
		// 只比較日期部分（忽略時間）
		if date.Date.Year() == now.Year() &&
			date.Date.Month() == now.Month() &&
			date.Date.Day() == now.Day() {
			return date.Date, nil
		}
		// 找到最近一個過去的交易日
		if date.Date.Before(now) || date.Date.Equal(now.Truncate(24*time.Hour)) {
			return date.Date, nil
		}
	}
	return time.Time{}, fmt.Errorf("找不到交易日資料")
}

func (m *marketDataGateway) GetLatestTradeDateByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error) {
	tradeDate, err := m.tradeDateReader.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return []time.Time{}, err
	}

	if len(tradeDate) == 0 {
		return []time.Time{}, fmt.Errorf("找不到交易日資料")
	}

	tradeDates := make([]time.Time, len(tradeDate))
	for i, date := range tradeDate {
		tradeDates[i] = date.Date
	}
	return tradeDates, nil
}

func (m *marketDataGateway) GetStockNews(ctx context.Context, symbol string) ([]dto.StockNews, error) {

	requestDto := finmindtradeDto.FinmindtradeRequestDto{
		DataID:    symbol,
		StartDate: time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
	}

	response, err := m.finmindAPI.GetTaiwanStockNews(requestDto)
	if err != nil {
		return nil, err
	}

	stockNews := make([]dto.StockNews, 0)
	for _, news := range response.Data {
		stockNews = append(stockNews, dto.StockNews{
			Date:        news.Date,
			StockSymbol: news.StockID,
			Link:        news.Link,
			Source:      news.Source,
			Title:       news.Title,
		})
	}
	return stockNews, nil
}
