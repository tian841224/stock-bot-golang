package dto

// 交易量排行
type TopVolume struct {
	// 股票代號
	StockSymbol string
	// 股票名稱
	StockName string
	// 日期
	Date string
	// 開盤價
	OpenPrice float64
	// 收盤價
	ClosePrice float64
	// 最高價
	HighPrice float64
	// 最低價
	LowPrice float64
	// 交易量
	Volume string
	// 交易筆數
	Transaction string
	// 成交金額
	Amount string
	// 漲跌金額
	ChangeAmount float64
	// 漲跌幅
	PercentageChange string
	// 漲跌方向
	UpDownSign string
}
