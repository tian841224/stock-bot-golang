// Package dto 提供 Fugle API 的 DTO 定義
package dto

// FugleCandlesRequestDto 股票Ｋ線
type FugleCandlesRequestDto struct {
	// 股票代碼
	Symbol string `json:"symbol"`
	// 開始日期（格式：yyyy-MM-dd）
	From string `json:"from"`
	// 結束日期（格式：yyyy-MM-dd）
	To string `json:"to"`
	// Ｋ線週期，可選 1 1分Ｋ；3 3分Ｋ；5 5分Ｋ；10 10分Ｋ；15 15分Ｋ；30 30分Ｋ；60 60分Ｋ；D 日Ｋ；W 週Ｋ；M 月Ｋ
	Timeframe string `json:"timeframe"`
	// 欄位，可選：open,high,low,close,volume,turnover,change
	Fields string `json:"fields"`
	// 時間排序，預設為 desc 降冪 ；可選 asc 升冪
	Sort string `json:"sort"`
}

// FugleCandlesResponseDto 股票Ｋ線回應
type FugleCandlesResponseDto struct {
	// 日期
	Date string `json:"date"`
	// 證券類型
	Type string `json:"type"`
	// 交易所
	Exchange string `json:"exchange"`
	// 市場別
	Market string `json:"market"`
	// 股票代號
	Symbol string `json:"symbol"`
	// Ｋ線週期
	Timeframe string `json:"timeframe"`
	// Ｋ線資料
	Data []FugleCandlesDataDto `json:"data"`
}

// FugleCandlesDataDto 股票Ｋ線資料
type FugleCandlesDataDto struct {
	// 日期（分 K 含時間）
	Date string `json:"date"`
	// Ｋ線開盤價
	Open float64 `json:"open"`
	// Ｋ線最高價
	High float64 `json:"high"`
	// Ｋ線最低價
	Low float64 `json:"low"`
	// Ｋ線收盤價
	Close float64 `json:"close"`
	// Ｋ線成交量（股）
	Volume float64 `json:"volume"`
	// Ｋ線成交金額（元）
	Turnover float64 `json:"turnover"`
	// Ｋ線漲跌
	Change float64 `json:"change"`
}
