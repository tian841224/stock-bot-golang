package dto

type TopVolumeItemsResponseDto struct {
	Stat   string          `json:"stat"`
	Date   string          `json:"date"`
	Title  string          `json:"title"`
	Fields []string        `json:"fields"`
	Data   [][]interface{} `json:"data"`
	Notes  []string        `json:"notes"`
	Total  int             `json:"total"`
}

type TopVolumeItemsData struct {
	Rank             string  `json:"rank"`
	StockId          string  `json:"stockId"`
	StockName        string  `json:"stockName"`
	Volume           string  `json:"volume"`
	Transaction      string  `json:"transaction"`
	OpenPrice        float64 `json:"openPrice"`
	HighPrice        float64 `json:"highPrice"`
	LowPrice         float64 `json:"lowPrice"`
	ClosePrice       float64 `json:"closePrice"`
	UpDownSign       string  `json:"upDownSign"`
	ChangeAmount     float64 `json:"changeAmount"`
	PercentageChange string  `json:"percentageChange"`
	BuyPrice         float64 `json:"buyPrice"`
	SellPrice        float64 `json:"sellPrice"`
}
