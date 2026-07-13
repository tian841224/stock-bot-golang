package dto

// 股票公司資訊
type StockCompanyInfo struct {
	// 基本資訊
	Symbol   string // 股票代碼
	Name     string // 股票名稱
	Industry string // 產業別
	Market   string // 市場別

	// 財務指標
	PE           float64 // 本益比
	PB           float64 // 本淨比
	MarketCap    float64 // 市值 (兆)
	BookValue    float64 // 每股淨值
	EPS          float64 // 近四季EPS
	QuarterEPS   float64 // 營季EPS
	Dividend     float64 // 年股利
	DividendRate float64 // 殖利率 (%)
	GrossMargin  float64 // 毛利率 (%)
	OperMargin   float64 // 營益率 (%)
	NetMargin    float64 // 淨利率 (%)
}
