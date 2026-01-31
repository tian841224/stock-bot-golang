package entity

// StockPrice 股票價格資訊
type StockPrice struct {
	CurrentPrice float64
	Change       float64
	ChangeRate   float64
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	PrevClose    float64
	Volume       int64
	Turnover     float64
	VolumeRatio  float64
	Amplitude    float64
}