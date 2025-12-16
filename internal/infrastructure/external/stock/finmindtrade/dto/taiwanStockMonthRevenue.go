package dto

// TaiwanStockMonthRevenueResponseDto 月營收表
type TaiwanStockMonthRevenueResponseDto struct {
	Msg    string                        `json:"msg"`
	Status int                           `json:"status"`
	Data   []TaiwanStockMonthRevenueData `json:"data"`
}
type TaiwanStockMonthRevenueData struct {
	Date         string `json:"date"`
	StockID      string `json:"stock_id"`
	Country      string `json:"country"`
	Revenue      int64  `json:"revenue"`
	RevenueMonth int64  `json:"revenue_month"`
	RevenueYear  int64  `json:"revenue_year"`
}
