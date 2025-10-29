package dto

// FugleStockQuoteRequestDto 股票即時報價
type FugleStockQuoteRequestDto struct {
	// 股票代碼
	Symbol string `json:"symbol"`
	// 類型，可選 oddlot 盤中零股
	Type string `json:"type,omitempty"`
}

// FugleStockQuoteResponseDto 股票即時報價回應
type FugleStockQuoteResponseDto struct {
	// 日期
	Date string `json:"date"`
	// 類型
	Type string `json:"type"`
	// 交易所
	Exchange string `json:"exchange"`
	// 市場
	Market string `json:"market"`
	// 股票代碼
	Symbol string `json:"symbol"`
	// 股票簡稱
	Name string `json:"name"`
	// 今日參考價
	ReferencePrice float64 `json:"referencePrice"`
	// 昨日收盤價
	PreviousClose float64 `json:"previousClose"`
	// 開盤價
	OpenPrice float64 `json:"openPrice"`
	// 開盤價成交時間
	OpenTime float64 `json:"openTime"`
	// 最高價
	HighPrice float64 `json:"highPrice"`
	// 最高價成交時間
	HighTime float64 `json:"highTime"`
	// 最低價
	LowPrice float64 `json:"lowPrice"`
	// 最低價成交時間
	LowTime float64 `json:"lowTime"`
	// 收盤價（最後成交價）
	ClosePrice float64 `json:"closePrice"`
	// 收盤價（最後成交價）成交時間
	CloseTime float64 `json:"closeTime"`
	// 當日成交均價
	AvgPrice float64 `json:"avgPrice"`
	// 最後一筆成交漲跌（含試撮）
	Change float64 `json:"change"`
	// 最後一筆成交漲跌幅（含試撮）
	ChangePercent float64 `json:"changePercent"`
	// 當日振幅
	Amplitude float64 `json:"amplitude"`
	LastPrice float64 `json:"lastPrice"`
	// 最後一筆成交數量（含試撮）
	LastSize float64 `json:"lastSize"`
	// 最佳五檔委買
	Bids []struct {
		Price float64 `json:"price"`
		Size  float64 `json:"size"`
	} `json:"bids"`
	// 最佳五檔委賣
	Asks []struct {
		Price float64 `json:"price"`
		Size  float64 `json:"size"`
	} `json:"asks"`
	// 統計資訊
	Total struct {
		TradeValue       float64 `json:"tradeValue"`
		TradeVolume      float64 `json:"tradeVolume"`
		TradeVolumeAtBid float64 `json:"tradeVolumeAtBid"`
		TradeVolumeAtAsk float64 `json:"tradeVolumeAtAsk"`
		Transaction      float64 `json:"transaction"`
		Time             float64 `json:"time"`
	} `json:"total"`
	// 最後一筆成交資訊
	LastTrade struct {
		Bid    float64 `json:"bid"`
		Ask    float64 `json:"ask"`
		Price  float64 `json:"price"`
		Size   float64 `json:"size"`
		Time   float64 `json:"time"`
		Serial float64 `json:"serial"`
	} `json:"lastTrade"`
	// 最後一筆試撮資訊
	LastTrial struct {
		Bid    float64 `json:"bid"`
		Ask    float64 `json:"ask"`
		Price  float64 `json:"price"`
		Size   float64 `json:"size"`
		Time   float64 `json:"time"`
		Serial float64 `json:"serial"`
	} `json:"lastTrial"`
	// 最後成交為逐筆交易
	IsContinuous bool `json:"isContinuous"`
	// 最後更新時間
	LastUpdated float64 `json:"lastUpdated"`
	// 流水號
	Serial float64 `json:"serial"`
}
