package dto

// TaiwanVariousIndicatorsResponseDto 台股各種指標
type TaiwanVariousIndicatorsResponseDto struct {
	Msg    string                        `json:"msg"`
	Status int                           `json:"status"`
	Data   []TaiwanVariousIndicatorsData `json:"data"`
}
type TaiwanVariousIndicatorsData struct {
	Date  string  `json:"date"`
	TAIEX float64 `json:"TAIEX"`
}
