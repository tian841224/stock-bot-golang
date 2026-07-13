package dto

// FugleMoversRequestDto 股票漲跌幅排行
type FugleMoversRequestDto struct {
	// 市場別
	Market string `json:"market"`
	// 方向
	Direction string `json:"direction"`
	// 漲跌／漲跌幅
	Change string `json:"change"`
	// 標的類型
	Type string `json:"type"`
	// 篩選大於漲跌／漲跌幅的股票
	Gt string `json:"gt"`
	// 篩選大於或等於漲跌／漲跌幅的股票
	Gte string `json:"gte"`
	// 篩選小於漲跌／漲跌幅的股票
	Lt string `json:"lt"`
	// 篩選小於或等於漲跌／漲跌幅的股票
	Lte string `json:"lte"`
	// 篩選等於漲跌／漲跌幅的股票
	Eq string `json:"eq"`
}

// FugleMoversResponseDto 股票漲跌幅排行回應
type FugleMoversResponseDto struct {
	// 日期
	Date string `json:"date"`
	// 時間
	Time string `json:"time"`
	// 市場別
	Market string `json:"market"`
	// 漲跌／漲跌幅
	Change string `json:"change"`
	// 快照資料
	Data []FugleMoversDataDto `json:"data"`
}

// FugleMoversDataDto 股票漲跌幅排行資料
type FugleMoversDataDto struct {
	// Ticker 類型
	Type string `json:"type"`
	// 股票代碼
	Symbol string `json:"symbol"`
	// 股票簡稱
	Name string `json:"name"`
	// 開盤價
	OpenPrice float64 `json:"openPrice"`
	// 最高價
	HighPrice float64 `json:"highPrice"`
	// 最低價
	LowPrice float64 `json:"lowPrice"`
	// 收盤價
	ClosePrice float64 `json:"closePrice"`
	// 漲跌
	Change float64 `json:"change"`
	// 漲跌幅
	ChangePercent float64 `json:"changePercent"`
	// 成交量
	TradeVolume float64 `json:"tradeVolume"`
	// 成交金額
	TradeValue float64 `json:"tradeValue"`
	// 最後更新時間
	LastUpdated float64 `json:"lastUpdated"`
}
