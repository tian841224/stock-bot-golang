package dto

import "time"

// StockPrice 股票價格
type StockPrice struct {
	// 股票代號
	Symbol string
	// 股票名稱
	Name string
	// 日期
	Date time.Time
	// 前一日收盤價
	PrevClosePrice float64
	// 開盤價
	OpenPrice float64
	// 收盤價
	ClosePrice float64
	// 最高價
	HighPrice float64
	// 最低價
	LowPrice float64
	// 漲跌金額
	ChangeAmount float64
	// 漲跌幅
	ChangeRate float64
	// 漲跌方向
	UpDownSign string
	// 交易量
	Volume int64
	// 交易筆數
	Transactions int64
	// 成交金額
	Amount float64
}
