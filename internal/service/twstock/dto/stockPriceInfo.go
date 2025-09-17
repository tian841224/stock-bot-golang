package twstock

type StockPriceInfo struct {
	StockID          string  `json:"stock_id"`
	StockName        string  `json:"stock_name"`
	Date             string  `json:"date"`
	OpenPrice        float64 `json:"open_price"`
	ClosePrice       float64 `json:"close_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	Volume           string  `json:"volume"`
	Transaction      string  `json:"transaction"`
	Amount           string  `json:"amount"`
	ChangeAmount     float64 `json:"change_amount"`
	PercentageChange string  `json:"percentage_change"`
	UpDownSign       string  `json:"up_down_sign"`
}
