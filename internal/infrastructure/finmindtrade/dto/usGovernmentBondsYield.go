package dto

// USGovernmentBondsYieldResponseDto 美國公債殖利率
type USGovernmentBondsYieldResponseDto struct {
	Msg    string                     `json:"msg"`
	Status int                        `json:"status"`
	Data   USGovernmentBondsYieldData `json:"data"`
}
type USGovernmentBondsYieldData struct {
	Date  string  `json:"date"`
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}
