package dto

// USStockInfoResponseDto 美股股票清單
type USStockInfoResponseDto struct {
	Msg    string        `json:"msg"`
	Status int           `json:"status"`
	Data   []USStockInfo `json:"data"`
}
type USStockInfo struct {
	Date      string  `json:"date"`
	StockID   string  `json:"stock_id"`
	Country   string  `json:"Country"`
	IPOYear   float64 `json:"IPOYear"`
	MarketCap float64 `json:"MarketCap"`
	Subsector string  `json:"Subsector"`
	StockName string  `json:"stock_name"`
}
