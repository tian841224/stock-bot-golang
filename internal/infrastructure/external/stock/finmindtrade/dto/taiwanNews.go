package dto

type TaiwanNewsResponseDto struct {
	Msg    string                   `json:"msg"`
	Status int                      `json:"status"`
	Data   []TaiwanNewsResponseData `json:"data"`
}

type TaiwanNewsResponseData struct {
	Date    string `json:"date"`
	StockID string `json:"stock_id"`
	Link    string `json:"link"`
	Source  string `json:"source"`
	Title   string `json:"title"`
}
