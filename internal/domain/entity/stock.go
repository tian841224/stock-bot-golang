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
