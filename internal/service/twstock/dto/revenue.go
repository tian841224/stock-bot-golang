package dto

type RevenueDto struct {
	// 時間戳記
	Time []int64 `json:"time"`
	// 股票代號
	Code string `json:"code"`
	// 股票名稱
	Name string `json:"name"`
	// 股價
	StockPrice []float64 `json:"stockPrice"`
	// 月營收
	SaleMonth []int64 `json:"saleMonth"`
	// 累積營收
	SaleAccumulated []int64 `json:"saleAccumulated"`
	// 月營收年增率
	YoY []float64 `json:"yoy"`
	// 累積營收年增率
	YoYAccumulated []float64 `json:"yoyAccumulated"`
}
