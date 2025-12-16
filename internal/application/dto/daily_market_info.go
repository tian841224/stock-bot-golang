package dto

// 大盤每日成交資訊
type DailyMarketInfo struct {
	// 日期
	Date string
	// 成交股數
	Volume string
	// 成交金額
	Amount string
	// 成交筆數
	Transaction string
	// 發行量加權股價指數
	Index string
	// 漲跌點數
	Change string
}
