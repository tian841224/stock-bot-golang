package dto

type StockQuoteInfo struct {
	// 基本資訊
	StockID   string `json:"stock_id"`   // 股票代碼
	StockName string `json:"stock_name"` // 股票名稱
	Industry  string `json:"industry"`   // 產業別
	Market    string `json:"market"`     // 市場別

	// 價格資訊
	CurrentPrice float64 `json:"current_price"` // 現價
	Change       float64 `json:"change"`        // 漲跌
	ChangeRate   float64 `json:"change_rate"`   // 漲跌幅 (%)
	OpenPrice    float64 `json:"open_price"`    // 開盤價
	HighPrice    float64 `json:"high_price"`    // 最高價
	LowPrice     float64 `json:"low_price"`     // 最低價
	PrevClose    float64 `json:"prev_close"`    // 昨收價

	// 成交量資訊
	Volume      int64   `json:"volume"`       // 成交量
	Turnover    float64 `json:"turnover"`     // 成交額 (億)
	VolumeRatio float64 `json:"volume_ratio"` // 週轉率 (%)
	Amplitude   float64 `json:"amplitude"`    // 振幅 (%)

	// 財務指標
	PE           float64 `json:"pe"`            // 本益比
	PB           float64 `json:"pb"`            // 本淨比
	MarketCap    float64 `json:"market_cap"`    // 市值 (兆)
	BookValue    float64 `json:"book_value"`    // 每股淨值
	EPS          float64 `json:"eps"`           // 近四季EPS
	QuarterEPS   float64 `json:"quarter_eps"`   // 營季EPS
	Dividend     float64 `json:"dividend"`      // 年股利
	DividendRate float64 `json:"dividend_rate"` // 殖利率 (%)
	GrossMargin  float64 `json:"gross_margin"`  // 毛利率 (%)
	OperMargin   float64 `json:"oper_margin"`   // 營益率 (%)
	NetMargin    float64 `json:"net_margin"`    // 淨利率 (%)

	// 價位區間
	UpperLimit  float64 `json:"upper_limit"`   // 漲停價
	LowerLimit  float64 `json:"lower_limit"`   // 跌停價
	High52W     float64 `json:"high_52w"`      // 52週高
	Low52W      float64 `json:"low_52w"`       // 52週低
	High52WDate string  `json:"high_52w_date"` // 52週高日期
	Low52WDate  string  `json:"low_52w_date"`  // 52週低日期

	// 五檔資訊
	BidPrices []float64 `json:"bid_prices"` // 買盤價格 [買1, 買2, 買3, 買4, 買5]
	AskPrices []float64 `json:"ask_prices"` // 賣盤價格 [賣1, 賣2, 賣3, 賣4, 賣5]

	// 內外盤資訊
	OutVolume int64   `json:"out_volume"` // 外盤量
	InVolume  int64   `json:"in_volume"`  // 內盤量
	OutRatio  float64 `json:"out_ratio"`  // 外盤比 (%)
	InRatio   float64 `json:"in_ratio"`   // 內盤比 (%)
}
