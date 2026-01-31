package entity

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
