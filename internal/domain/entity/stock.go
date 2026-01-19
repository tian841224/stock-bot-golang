package entity

import (
	"fmt"
	"time"
)

// Stock 股票領域模型
type Stock struct {
	ID          string
	Name        string
	Symbol      string
	Industry    string
	Market      string
	CurrentInfo *StockPrice
	Financials  *FinancialMetrics
	MarketData  *MarketMetrics
}

// StockPrice 股票價格資訊
type StockPrice struct {
	CurrentPrice float64
	Change       float64
	ChangeRate   float64
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	PrevClose    float64
	Volume       int64
	Turnover     float64
	VolumeRatio  float64
	Amplitude    float64
}

// StockPerformance 股票績效
type StockPerformance struct {
	StockID     string
	Period      string
	PeriodName  string
	Performance string
}

// StockNews 股票新聞資料
type StockNews struct {
	StockID     string
	Title       string
	Summary     string
	Link        string
	Source      string
	PublishedAt time.Time
}

// FinancialMetrics 財務指標
type FinancialMetrics struct {
	PE           float64
	PB           float64
	MarketCap    float64
	BookValue    float64
	EPS          float64
	QuarterEPS   float64
	Dividend     float64
	DividendRate float64
	GrossMargin  float64
	OperMargin   float64
	NetMargin    float64
}

// MarketMetrics 市場指標
type MarketMetrics struct {
	UpperLimit  float64
	LowerLimit  float64
	High52W     float64
	Low52W      float64
	High52WDate time.Time
	Low52WDate  time.Time
	BidPrices   []float64
	AskPrices   []float64
	OutVolume   int64
	InVolume    int64
	OutRatio    float64
	InRatio     float64
}

// Revenue 營收資料
type Revenue struct {
	StockID         string
	Time            []int64
	StockPrice      []float64
	SaleMonth       []int64
	SaleAccumulated []int64
	YoY             []float64
	YoYAccumulated  []float64
}

// AfterTradingVolume 盤後交易量資料
type AfterTradingVolume struct {
	StockID          string
	Volume           int64
	Transaction      int64
	Amount           float64
	OpenPrice        float64
	ClosePrice       float64
	HighPrice        float64
	LowPrice         float64
	ChangeAmount     float64
	PercentageChange float64
	UpDownSign       string
}

// StockQuote 股票報價資料
type StockQuote struct {
	// 日期
	Date string
	// 類型
	Type string
	// 交易所
	Exchange string
	// 市場
	Market string
	// 股票代碼
	Symbol string
	// 股票簡稱
	Name string
	// 今日參考價
	ReferencePrice float64
	// 昨日收盤價
	PreviousClose float64
	// 開盤價
	OpenPrice float64
	// 開盤價成交時間
	OpenTime float64
	// 最高價
	HighPrice float64
	// 最高價成交時間
	HighTime float64
	// 最低價
	LowPrice float64
	// 最低價成交時間
	LowTime float64
	// 收盤價（最後成交價）
	ClosePrice float64
	// 收盤價（最後成交價）成交時間
	CloseTime float64
	// 當日成交均價
	AvgPrice float64
	// 最後一筆成交漲跌（含試撮）
	Change float64
	// 最後一筆成交漲跌幅（含試撮）
	ChangePercent float64
	// 當日振幅
	Amplitude float64
	LastPrice float64
	// 最後一筆成交數量（含試撮）
	LastSize float64
	// 最佳五檔委買
	Bids []struct {
		Price float64
		Size  float64
	}
	// 最佳五檔委賣
	Asks []struct {
		Price float64
		Size  float64
	}
	// 統計資訊
	Total struct {
		TradeValue       float64
		TradeVolume      float64
		TradeVolumeAtBid float64
		TradeVolumeAtAsk float64
		Transaction      float64
		Time             float64
	}
	// 最後一筆成交資訊
	LastTrade struct {
		Bid    float64
		Ask    float64
		Price  float64
		Size   float64
		Time   float64
		Serial float64
	}
	// 最後一筆試撮資訊
	LastTrial struct {
		Bid    float64
		Ask    float64
		Price  float64
		Size   float64
		Time   float64
		Serial float64
	}
	// 最後成交為逐筆交易
	IsContinuous bool
	// 最後更新時間
	LastUpdated float64
	// 流水號
	Serial float64
}

// GetMarketCapInTrillions 取得市值（兆元）
func (s *Stock) GetMarketCapInTrillions() float64 {
	if s.Financials == nil {
		return 0
	}
	return s.Financials.MarketCap / 1e12
}

// GetTurnoverInBillions 取得成交額（億元）
func (s *Stock) GetTurnoverInBillions() float64 {
	if s.CurrentInfo == nil {
		return 0
	}
	return s.CurrentInfo.Turnover / 1e8
}

// CalculatePriceChange 計算價格變動
func (s *Stock) CalculatePriceChange(current, previous float64) (change float64, changeRate float64) {
	change = current - previous
	if previous != 0 {
		changeRate = (change / previous) * 100
	}
	return change, changeRate
}

// CalculateAmplitude 計算振幅
func (s *Stock) CalculateAmplitude(high, low, prevClose float64) float64 {
	if prevClose == 0 {
		return 0
	}

	amplitude := ((high - low) / prevClose) * 100
	return amplitude
}

// CalculateVolumeRatio 計算週轉率
func (s *Stock) CalculateVolumeRatio(volume, avgVolume float64) float64 {
	if avgVolume == 0 {
		return 0
	}
	return (volume / avgVolume) * 100
}

// CalculatePerformance 計算績效
func (s *Stock) CalculatePerformance(stock *Stock, period string) (*StockPerformance, error) {
	if stock == nil {
		return nil, fmt.Errorf("股票資料不能為空")
	}

	performance := &StockPerformance{
		StockID:     stock.ID,
		Period:      period,
		PeriodName:  s.getPeriodName(period),
		Performance: s.calculatePerformanceValue(stock),
	}

	return performance, nil
}

// getPeriodName 取得週期名稱
func (s *Stock) getPeriodName(period string) string {
	periodMap := map[string]string{
		"1d":  "1日",
		"1w":  "1週",
		"1m":  "1月",
		"3m":  "3月",
		"6m":  "6月",
		"1y":  "1年",
		"ytd": "年初至今",
	}

	if name, exists := periodMap[period]; exists {
		return name
	}
	return period
}

// calculatePerformanceValue 計算績效值
func (s *Stock) calculatePerformanceValue(stock *Stock) string {
	if stock.CurrentInfo == nil {
		return "無資料"
	}

	// 根據期間和價格變動計算績效
	changeRate := stock.CurrentInfo.ChangeRate

	if changeRate > 0 {
		return fmt.Sprintf("上漲 %.2f%%", changeRate)
	} else if changeRate < 0 {
		return fmt.Sprintf("下跌 %.2f%%", changeRate)
	}

	return "持平"
}
