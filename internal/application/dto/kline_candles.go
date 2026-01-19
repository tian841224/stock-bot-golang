package dto

// KlineCandles 股票Ｋ線回應
type KlineCandles struct {
	// 日期
	Date string
	// 證券類型
	Type string
	// 交易所
	Exchange string
	// 市場別
	Market string
	// 股票代號
	Symbol string
	// Ｋ線週期
	Timeframe string
	// Ｋ線資料
	Data []KlineCandlesData
}

// KlineCandlesData 股票Ｋ線資料
type KlineCandlesData struct {
	// 日期（分 K 含時間）
	Date string
	// Ｋ線開盤價
	Open float64
	// Ｋ線最高價
	High float64
	// Ｋ線最低價
	Low float64
	// Ｋ線收盤價
	Close float64
	// Ｋ線成交量（股）
	Volume float64
	// Ｋ線成交金額（元）
	Turnover float64
	// Ｋ線漲跌
	Change float64
}
