package dto

// USStockPriceResponseDto 美股盤後股價回應
type USStockPriceResponseDto struct {
	Msg    string           `json:"msg"`
	Status int              `json:"status"`
	Data   USStockPriceData `json:"data"`
}
type USStockPriceData struct {
	Date    string  `json:"date"`
	StockID string  `json:"stock_id"`
	Close   float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Open    float64 `json:"open"`
	Volume  int64   `json:"volume"`
}
