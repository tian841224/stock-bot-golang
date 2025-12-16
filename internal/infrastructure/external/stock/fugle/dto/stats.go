// Package dto 提供 Fugle API 的 DTO 定義
package dto

// FugleStatsRequestDto 股票統計資訊
type FugleStatsRequestDto struct {
	// 股票代碼
	Symbol string `json:"symbol"`
}

// FugleStatsResponseDto 股票統計資訊回應
type FugleStatsResponseDto struct {
	// 日期
	Date string `json:"date"`
	// Ticker 類型
	Type string `json:"type"`
	// 交易所
	Exchange string `json:"exchange"`
	// 市場別
	Market string `json:"market"`
	// 股票代碼
	Symbol string `json:"symbol"`
	// 股票簡稱
	Name string `json:"name"`
	// 最後交易日開盤價
	OpenPrice float64 `json:"openPrice"`
	// 最後交易日最高價
	HighPrice float64 `json:"highPrice"`
	// 最後交易日最低價
	LowPrice float64 `json:"lowPrice"`
	// 最後交易日收盤價
	ClosePrice float64 `json:"closePrice"`
	// 最後交易日漲跌
	Change float64 `json:"change"`
	// 最後交易日漲跌幅
	ChangePercent float64 `json:"changePercent"`
	// 最後交易日成交量
	TradeVolume float64 `json:"tradeVolume"`
	// 最後交易日成交金額
	TradeValue float64 `json:"tradeValue"`
	// 前一交易日收盤價
	PreviousClose float64 `json:"previousClose"`
	// 近 52 週高點
	Week52High float64 `json:"week52High"`
	// 近 52 週低點
	Week52Low float64 `json:"week52Low"`
}
