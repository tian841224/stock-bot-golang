package dto

// 原始 API 回應結構
type AfterTradingVolumeRawResponseDto struct {
	Tables []struct {
		Data [][]interface{} `json:"data"`
	} `json:"tables"`
}

// 處理後的資料結構
type AfterTradingVolumeResponseDto struct {
	StockId          string  `json:"stockId"`
	StockName        string  `json:"stockName"`
	Volume           string  `json:"volume"`
	Transaction      string  `json:"transaction"`
	Amount           string  `json:"amount"`
	OpenPrice        float64 `json:"openPrice"`
	ClosePrice       float64 `json:"closePrice"`
	HighPrice        float64 `json:"highPrice"`
	LowPrice         float64 `json:"lowPrice"`
	UpDownSign       string  `json:"upDownSign"`
	ChangeAmount     float64 `json:"changeAmount"`
	PercentageChange string  `json:"percentageChange"`
}
