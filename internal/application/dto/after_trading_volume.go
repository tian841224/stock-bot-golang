package dto

// 處理後的資料結構
type AfterTradingVolume struct {
	StockId          string
	StockName        string
	Volume           string
	Transaction      string
	Amount           string
	OpenPrice        float64
	ClosePrice       float64
	HighPrice        float64
	LowPrice         float64
	UpDownSign       string
	ChangeAmount     float64
	PercentageChange string
}
