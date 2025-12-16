package dto

// StockPerformance 聚合股票績效資料。
type StockPerformance struct {
	Symbol string
	Name   string
	Data   []StockPerformanceData
}

type StockPerformanceData struct {
	Period      string
	PeriodName  string
	Performance string
}
