package twstock

type StockPriceInfo struct {
	StockID          string  `json:"stock_id"`
	StockName        string  `json:"stock_name"`
	Date             string  `json:"date"`
	OpenPrice        float64 `json:"open_price"`
	ClosePrice       float64 `json:"close_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	Volume           int64   `json:"volume"`
	Transaction      int64   `json:"transaction"`
	Amount           int64   `json:"amount"`
	ChangeAmount     float64 `json:"change_amount"`
	PercentageChange string  `json:"percentage_change"`
	UpDownSign       string  `json:"up_down_sign"`
}
