package entity

// StockQuote 股票報價資料
type StockQuote struct {
	// 日期
	Date string
	// 類型
	Type string
	// 交易所
	Exchange string
	// 市場
	Market string
	// 股票代碼
	Symbol string
	// 股票簡稱
	Name string
	// 今日參考價
	ReferencePrice float64
	// 昨日收盤價
	PreviousClose float64
	// 開盤價
	OpenPrice float64
	// 開盤價成交時間
	OpenTime float64
	// 最高價
	HighPrice float64
	// 最高價成交時間
	HighTime float64
	// 最低價
	LowPrice float64
	// 最低價成交時間
	LowTime float64
	// 收盤價（最後成交價）
	ClosePrice float64
	// 收盤價（最後成交價）成交時間
	CloseTime float64
	// 當日成交均價
	AvgPrice float64
	// 最後一筆成交漲跌（含試撮）
	Change float64
	// 最後一筆成交漲跌幅（含試撮）
	ChangePercent float64
	// 當日振幅
	Amplitude float64
	LastPrice float64
	// 最後一筆成交數量（含試撮）
	LastSize float64
	// 最佳五檔委買
	Bids []struct {
		Price float64
		Size  float64
	}
	// 最佳五檔委賣
	Asks []struct {
		Price float64
		Size  float64
	}
	// 統計資訊
	Total struct {
		TradeValue       float64
		TradeVolume      float64
		TradeVolumeAtBid float64
		TradeVolumeAtAsk float64
		Transaction      float64
		Time             float64
	}
	// 最後一筆成交資訊
	LastTrade struct {
		Bid    float64
		Ask    float64
		Price  float64
		Size   float64
		Time   float64
		Serial float64
	}
	// 最後一筆試撮資訊
	LastTrial struct {
		Bid    float64
		Ask    float64
		Price  float64
		Size   float64
		Time   float64
		Serial float64
	}
	// 最後成交為逐筆交易
	IsContinuous bool
	// 最後更新時間
	LastUpdated float64
	// 流水號
	Serial float64
}
