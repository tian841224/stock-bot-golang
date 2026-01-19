package dto

type FinmindtradeRequestDto struct {
	DataSet   string `json:"dataset"`
	StockID   string `json:"stock_id"`
	DataID    string `json:"data_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
