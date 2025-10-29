package mappers

import (
	"time"

	"github.com/tian841224/stock-bot/internal/domain/stock"
	cnyesDto "github.com/tian841224/stock-bot/internal/infrastructure/cnyes/dto"
	finmindDto "github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade/dto"
	stockDto "github.com/tian841224/stock-bot/internal/service/twstock/dto"
)

// StockMapper 股票資料轉換器
type StockMapper struct{}

// NewStockMapper 建立股票轉換器
func NewStockMapper() *StockMapper {
	return &StockMapper{}
}

// FromCnyesDto 從鉅亨網 DTO 轉換為領域模型
func (m *StockMapper) FromCnyesDto(dto cnyesDto.CnyesStockQuoteDataDto) (*stock.Stock, error) {
	// 建立股票代號
	stockID, err := stock.NewStockID(dto.StockID)
	if err != nil {
		return nil, err
	}

	// 建立價格資訊
	priceInfo := m.buildStockPrice(dto)

	// 建立財務指標
	financials := m.buildFinancialMetrics(dto)

	// 建立市場指標
	marketData := m.buildMarketMetrics(dto)

	return &stock.Stock{
		ID:          stockID.String(),
		Name:        dto.StockName,
		Symbol:      dto.StockID,
		Industry:    dto.Industry,
		Market:      dto.Market,
		CurrentInfo: priceInfo,
		Financials:  financials,
		MarketData:  marketData,
	}, nil
}

// FromFinmindDto 從 FinMind DTO 轉換為領域模型
func (m *StockMapper) FromFinmindDto(dto finmindDto.TaiwanStockPriceData) (*stock.Stock, error) {
	// 建立股票代號
	stockID, err := stock.NewStockID(dto.StockID)
	if err != nil {
		return nil, err
	}

	// 建立價格資訊
	priceInfo, err := m.buildStockPriceFromFinmind(dto)
	if err != nil {
		return nil, err
	}

	return &stock.Stock{
		ID:          stockID.String(),
		Name:        "", // FinMind 資料中沒有股票名稱
		Symbol:      dto.StockID,
		CurrentInfo: priceInfo,
	}, nil
}

// ToStockQuoteDto 轉換為股票報價 DTO
func (m *StockMapper) ToStockQuoteDto(stock *stock.Stock) *stockDto.StockQuoteInfo {
	if stock == nil {
		return nil
	}

	dto := &stockDto.StockQuoteInfo{
		StockID:   stock.ID,
		StockName: stock.Name,
		Industry:  stock.Industry,
		Market:    stock.Market,
	}

	// 轉換價格資訊
	if stock.CurrentInfo != nil {
		dto.CurrentPrice = stock.CurrentInfo.CurrentPrice
		dto.Change = stock.CurrentInfo.Change
		dto.ChangeRate = stock.CurrentInfo.ChangeRate
		dto.OpenPrice = stock.CurrentInfo.OpenPrice
		dto.HighPrice = stock.CurrentInfo.HighPrice
		dto.LowPrice = stock.CurrentInfo.LowPrice
		dto.PrevClose = stock.CurrentInfo.PrevClose
		dto.Volume = stock.CurrentInfo.Volume
		dto.Turnover = stock.CurrentInfo.Turnover
		dto.VolumeRatio = stock.CurrentInfo.VolumeRatio
		dto.Amplitude = stock.CurrentInfo.Amplitude
	}

	// 轉換財務指標
	if stock.Financials != nil {
		dto.PE = stock.Financials.PE
		dto.PB = stock.Financials.PB
		dto.MarketCap = stock.Financials.MarketCap
		dto.BookValue = stock.Financials.BookValue
		dto.EPS = stock.Financials.EPS
		dto.QuarterEPS = stock.Financials.QuarterEPS
		dto.Dividend = stock.Financials.Dividend
		dto.DividendRate = stock.Financials.DividendRate
		dto.GrossMargin = stock.Financials.GrossMargin
		dto.OperMargin = stock.Financials.OperMargin
		dto.NetMargin = stock.Financials.NetMargin
	}

	// 轉換市場指標
	if stock.MarketData != nil {
		dto.UpperLimit = stock.MarketData.UpperLimit
		dto.LowerLimit = stock.MarketData.LowerLimit
		dto.High52W = stock.MarketData.High52W
		dto.Low52W = stock.MarketData.Low52W
		dto.High52WDate = stock.MarketData.High52WDate.Format("2006-01-02")
		dto.Low52WDate = stock.MarketData.Low52WDate.Format("2006-01-02")
		dto.BidPrices = stock.MarketData.BidPrices
		dto.AskPrices = stock.MarketData.AskPrices
		dto.OutVolume = stock.MarketData.OutVolume
		dto.InVolume = stock.MarketData.InVolume
		dto.OutRatio = stock.MarketData.OutRatio
		dto.InRatio = stock.MarketData.InRatio
	}

	return dto
}

// buildStockPrice 建立股票價格資訊
func (m *StockMapper) buildStockPrice(dto cnyesDto.CnyesStockQuoteDataDto) *stock.StockPrice {
	updateTime := time.Unix(dto.UpdateTime, 0)

	return &stock.StockPrice{
		CurrentPrice: dto.CurrentPrice,
		Change:       dto.Change,
		ChangeRate:   dto.ChangeRate,
		OpenPrice:    dto.OpenPrice,
		HighPrice:    dto.HighPrice,
		LowPrice:     dto.LowPrice,
		PrevClose:    dto.PrevClose,
		Volume:       int64(dto.Volume),
		Turnover:     dto.Turnover / 1e8,    // 轉換為億元
		VolumeRatio:  dto.VolumeRatio * 100, // 轉換為百分比
		Amplitude:    dto.Amplitude,
		UpdateTime:   updateTime,
	}
}

// buildStockPriceFromFinmind 從 FinMind 資料建立股票價格資訊
func (m *StockMapper) buildStockPriceFromFinmind(dto finmindDto.TaiwanStockPriceData) (*stock.StockPrice, error) {
	date, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return nil, err
	}

	// 計算漲跌和漲跌幅
	change := dto.Close - dto.Open
	changeRate := 0.0
	if dto.Open != 0 {
		changeRate = (change / dto.Open) * 100
	}

	return &stock.StockPrice{
		CurrentPrice: dto.Close,
		Change:       change,
		ChangeRate:   changeRate,
		OpenPrice:    dto.Open,
		HighPrice:    dto.Max,
		LowPrice:     dto.Min,
		PrevClose:    dto.Open, // FinMind 資料中沒有昨收價，使用開盤價
		Volume:       dto.TradingVolume,
		Turnover:     float64(dto.TradingMoney) / 1e8, // 轉換為億元
		UpdateTime:   date,
	}, nil
}

// buildFinancialMetrics 建立財務指標
func (m *StockMapper) buildFinancialMetrics(dto cnyesDto.CnyesStockQuoteDataDto) *stock.FinancialMetrics {
	return &stock.FinancialMetrics{
		PE:           dto.PE,
		PB:           dto.PB,
		MarketCap:    dto.MarketCap / 1e12, // 轉換為兆元
		BookValue:    dto.BookValue,
		EPS:          dto.EPS,
		QuarterEPS:   dto.QuarterEPS,
		Dividend:     dto.Dividend,
		DividendRate: dto.DividendRate,
		GrossMargin:  dto.GrossMargin,
		OperMargin:   dto.OperMargin,
		NetMargin:    dto.NetMargin,
	}
}

// buildMarketMetrics 建立市場指標
func (m *StockMapper) buildMarketMetrics(dto cnyesDto.CnyesStockQuoteDataDto) *stock.MarketMetrics {
	// 解析 52 週高低點日期
	high52WDate, _ := time.Parse("2006-01-02", dto.High52WDate)
	low52WDate, _ := time.Parse("2006-01-02", dto.Low52WDate)

	return &stock.MarketMetrics{
		UpperLimit:  dto.UpperLimit,
		LowerLimit:  dto.LowerLimit,
		High52W:     dto.High52W,
		Low52W:      dto.Low52W,
		High52WDate: high52WDate,
		Low52WDate:  low52WDate,
		BidPrices: []float64{
			dto.BidPrice1, dto.BidPrice2, dto.BidPrice3, dto.BidPrice4, dto.BidPrice5,
		},
		AskPrices: []float64{
			dto.AskPrice1, dto.AskPrice2, dto.AskPrice3, dto.AskPrice4, dto.AskPrice5,
		},
		OutVolume: int64(dto.OutVolume),
		InVolume:  int64(dto.InVolume),
		OutRatio:  dto.OutRatio,
		InRatio:   dto.InRatio,
	}
}
