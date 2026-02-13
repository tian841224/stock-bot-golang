package entity

import (
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