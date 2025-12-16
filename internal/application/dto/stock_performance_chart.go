package dto

type StockPerformanceChart struct {
	Symbol    string
	StockName string
	Data      []StockPerformanceData
	ChartData []byte
}
