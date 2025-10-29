// Package dto 提供鉅亨網 API 的資料傳輸物件
package dto

// CnyesStockQuoteResponseDto 鉅亨網股票報價回應
type CnyesStockQuoteResponseDto struct {
	StatusCode int                      `json:"statusCode"`
	Message    string                   `json:"message"`
	Data       []CnyesStockQuoteDataDto `json:"data"`
}

// CnyesStockQuoteDataDto 股票報價資料
type CnyesStockQuoteDataDto struct {
	// 基本資訊
	Symbol    string `json:"0"`      // 股票代碼 (TWS:2330:STOCK)
	StockID   string `json:"200010"` // 股票代碼 (2330)
	StockName string `json:"200009"` // 股票名稱 (台積電)
	Industry  string `json:"200087"` // 產業別 (半導體業)
	Market    string `json:"200222"` // 市場別 (上市)
	Exchange  string `json:"200011"` // 交易所 (TSE)
	StockType string `json:"200061"` // 股票類型 (COMMON)

	// 價格資訊
	CurrentPrice float64 `json:"6"`  // 現價
	Change       float64 `json:"11"` // 漲跌
	ChangeRate   float64 `json:"56"` // 漲跌幅 (%)
	OpenPrice    float64 `json:"19"` // 開盤價
	HighPrice    float64 `json:"12"` // 最高價
	LowPrice     float64 `json:"13"` // 最低價
	PrevClose    float64 `json:"21"` // 昨收價

	// 成交量資訊
	Volume      float64 `json:"800001"` // 成交量
	Turnover    float64 `json:"200067"` // 成交額
	VolumeRatio float64 `json:"200127"` // 量比/週轉率
	Amplitude   float64 `json:"200124"` // 振幅

	// 財務指標
	PE           float64 `json:"36"`     // 本益比
	PB           float64 `json:"700006"` // 本淨比
	MarketCap    float64 `json:"700005"` // 市值
	BookValue    float64 `json:"200216"` // 每股淨值
	EPS          float64 `json:"34"`     // 近四季EPS
	QuarterEPS   float64 `json:"200223"` // 營季EPS
	Dividend     float64 `json:"200224"` // 年股利
	DividendRate float64 `json:"200225"` // 殖利率
	GrossMargin  float64 `json:"200220"` // 毛利率
	OperMargin   float64 `json:"200221"` // 營益率
	NetMargin    float64 `json:"200219"` // 淨利率

	// 價位區間
	UpperLimit  float64 `json:"75"`     // 漲停價
	LowerLimit  float64 `json:"76"`     // 跌停價
	High52W     float64 `json:"200192"` // 52週高
	Low52W      float64 `json:"200196"` // 52週低
	High52WDate string  `json:"200193"` // 52週高日期
	Low52WDate  string  `json:"200197"` // 52週低日期

	// 五檔資訊 (買盤)
	BidPrice1 float64 `json:"436"` // 買1價
	BidPrice2 float64 `json:"437"` // 買2價
	BidPrice3 float64 `json:"438"` // 買3價
	BidPrice4 float64 `json:"439"` // 買4價
	BidPrice5 float64 `json:"440"` // 買5價

	// 五檔資訊 (賣盤)
	AskPrice1 float64 `json:"441"` // 賣1價
	AskPrice2 float64 `json:"442"` // 賣2價
	AskPrice3 float64 `json:"443"` // 賣3價
	AskPrice4 float64 `json:"444"` // 賣4價
	AskPrice5 float64 `json:"445"` // 賣5價

	// 其他資訊
	UpdateTime  int64   `json:"200007"` // 更新時間戳
	Status      int     `json:"800002"` // 狀態
	AvgPrice    float64 `json:"3404"`   // 均價
	TotalVolume float64 `json:"200013"` // 總成交量
	OutVolume   float64 `json:"200055"` // 外盤量
	InVolume    float64 `json:"200054"` // 內盤量
	OutRatio    float64 `json:"200057"` // 外盤比
	InRatio     float64 `json:"200056"` // 內盤比
}
