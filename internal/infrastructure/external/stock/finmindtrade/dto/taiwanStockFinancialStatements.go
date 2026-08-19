package dto

// TaiwanStockFinancialStatementsResponseDto 綜合損益表
type TaiwanStockFinancialStatementsResponseDto struct {
	Msg    string                               `json:"msg"`
	Status int                                  `json:"status"`
	Data   []TaiwanStockFinancialStatementsData `json:"data"`
}

type TaiwanStockFinancialStatementsData struct {
	Date       string  `json:"date"`
	StockID    string  `json:"stock_id"`
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	OriginName string  `json:"origin_name"`
}
