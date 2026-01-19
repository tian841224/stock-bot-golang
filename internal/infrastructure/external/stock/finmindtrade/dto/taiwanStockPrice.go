package dto

// TaiwanStockPriceResponseDto 台股盤後資訊
type TaiwanStockPriceResponseDto struct {
	Msg    string                 `json:"msg"`
	Status int                    `json:"status"`
	Data   []TaiwanStockPriceData `json:"data"`
}
type TaiwanStockPriceData struct {
	Date            string  `json:"date"`
	StockID         string  `json:"stock_id"`
	TradingVolume   int64   `json:"Trading_Volume"`
	TradingMoney    int64   `json:"Trading_money"`
	Open            float64 `json:"open"`
	Max             float64 `json:"max"`
	Min             float64 `json:"min"`
	Close           float64 `json:"close"`
	Spread          float64 `json:"spread"`
	TradingTurnover float32 `json:"Trading_turnover"`
}
