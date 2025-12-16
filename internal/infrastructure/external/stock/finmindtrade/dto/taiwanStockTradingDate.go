package dto

// TaiwanStockTradingDateResponseDto 台股交易日
type TaiwanStockTradingDateResponseDto struct {
	Msg    string                       `json:"msg"`
	Status int                          `json:"status"`
	Data   []TaiwanStockTradingDateData `json:"data"`
}
type TaiwanStockTradingDateData struct {
	Date string `json:"date"`
}
