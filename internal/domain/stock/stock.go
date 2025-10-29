package stock

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
	UpdateTime   time.Time
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

// 業務方法

// IsValid 檢查股票是否有效
func (s *Stock) IsValid() bool {
	return s.ID != "" && s.Name != "" && s.Symbol != ""
}

// GetPriceChangeStatus 取得價格變動狀態
func (s *Stock) GetPriceChangeStatus() string {
	if s.CurrentInfo == nil {
		return "未知"
	}

	if s.CurrentInfo.ChangeRate > 0 {
		return "上漲"
	} else if s.CurrentInfo.ChangeRate < 0 {
		return "下跌"
	}
	return "持平"
}

// IsTradingDay 檢查是否為交易日
func (s *Stock) IsTradingDay() bool {
	if s.CurrentInfo == nil {
		return false
	}

	// 檢查更新時間是否為今天
	now := time.Now()
	return s.CurrentInfo.UpdateTime.Year() == now.Year() &&
		s.CurrentInfo.UpdateTime.Month() == now.Month() &&
		s.CurrentInfo.UpdateTime.Day() == now.Day()
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
