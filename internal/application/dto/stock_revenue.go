package dto

// StockRevenue 股票營收資料
type StockRevenue struct {
	// 股票代號
	StockSymbol string
	// 股票名稱
	StockName string
	// 時間戳記
	Time []int64
	// 股票價格
	StockPrice []float64
	// 月營收
	SaleMonth []int64
	// 累積營收
	SaleAccumulated []int64
	// 月營收年增率
	YoY []float64
	// 累積營收年增率
	YoYAccumulated []float64
}
