package dto

// TaiwanExchangeRateResponseDto 兌台幣匯率
type TaiwanExchangeRateResponseDto struct {
	Msg    string                   `json:"msg"`
	Status int                      `json:"status"`
	Data   []TaiwanExchangeRateData `json:"data"`
}
type TaiwanExchangeRateData struct {
	Date     string  `json:"date"`
	Currency string  `json:"currency"`
	CashBuy  float64 `json:"cash_buy"`
	CashSell float64 `json:"cash_sell"`
	SpotBuy  float64 `json:"spot_buy"`
	SpotSell float64 `json:"spot_sell"`
}
