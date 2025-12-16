// Package dto 提供 Fugle API 的 DTO 定義
package dto

// FugleActivesRequestDto 股票成交量值排行（依市場別）
type FugleActivesRequestDto struct {
	// 市場別 可選 TSE 上市；OTC 上櫃；ESB 興櫃一般板；TIB 臺灣創新板；PSB 興櫃戰略新板
	Market string `json:"market"`
	// 成交量／成交值 可選 volume 成交量；value 成交值
	Trade string `json:"trade"`
	// 標的類型 可選 ALLBUT0999 包含一般股票、特別股及 ETF ； COMMONSTOCK 為一般股票
	Type string `json:"type"`
}

// FugleActivesResponseDto 股票成交量值排行（依市場別）回應
type FugleActivesResponseDto struct {
	// 日期
	Date string `json:"date"`
	// 時間
	Time string `json:"time"`
	// 市場別
	Market string `json:"market"`
	// 成交量／成交值
	Trade string `json:"trade"`
	// 快照資料
	Data []FugleActivesDataDto `json:"data"`
}

// FugleActivesDataDto 股票成交量值排行（依市場別）資料
type FugleActivesDataDto struct {
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
