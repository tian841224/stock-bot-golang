package dto

// TaiwanStockSplitPriceResponseDto 台股分割資訊
type TaiwanStockSplitPriceResponseDto struct {
	Msg    string                      `json:"msg"`
	Status int                         `json:"status"`
	Data   []TaiwanStockSplitPriceData `json:"data"`
}

type TaiwanStockSplitPriceData struct {
	Date        string  `json:"date"`
	StockID     string  `json:"stock_id"`
	Type        string  `json:"type"`
	BeforePrice float64 `json:"before_price"`
	AfterPrice  float64 `json:"after_price"`
	MaxPrice    float64 `json:"max_price"`
	MinPrice    float64 `json:"min_price"`
	OpenPrice   float64 `json:"open_price"`
}
