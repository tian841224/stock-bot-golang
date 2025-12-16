package dto

// 大盤每日成交資訊（採用 TWSE MI_INDEX20 回應格式）
type DailyMarketInfoResponseDto struct {
	Stat   string     `json:"stat"`
	Date   string     `json:"date"`
	Title  string     `json:"title"`
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data"` // 修改為字串陣列的陣列
	Notes  []string   `json:"notes"`
	Total  int        `json:"total"`
}

// 大盤每日成交資訊
type DailyMarketInfoData struct {
	// 日期
	Date string `json:"date"`
	// 成交股數
	Volume string `json:"volume"`
	// 成交金額
	Amount string `json:"amount"`
	// 成交筆數
	Transaction string `json:"transaction"`
	// 發行量加權股價指數
	Index string `json:"index"`
	// 漲跌點數
	Change string `json:"change"`
}
