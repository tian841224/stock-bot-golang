package dto

type TaiwanStockAnalysisResponseDto struct {
	Msg    string                  `json:"msg"`
	Status int                     `json:"status"`
	Data   TaiwanStockAnalysisData `json:"data"`
}
type TaiwanStockAnalysisData struct {
	StockPrice                StockPriceData                `json:"StockPrice"`
	InstitutionalInvestor     TSAInstitutionalInvestorData  `json:"InstitutionalInvestor"`
	Shareholding              ShareholdingData              `json:"Shareholding"`
	MarginPurchaseShortSale   MarginPurchaseShortSaleData   `json:"MarginPurchaseShortSale"`
	TaiwanFinancialStatements TaiwanFinancialStatementsData `json:"TaiwanFinancialStatements"`
	TaiwanBalanceSheet        TaiwanBalanceSheetData        `json:"TaiwanBalanceSheet"`
	TaiwanCashFlowsStatement  TaiwanCashFlowsStatementData  `json:"TaiwanCashFlowsStatement"`
	TaiwanStockDividend       TSATaiwanStockDividendData    `json:"TaiwanStockDividend"`
	TaiwanNews                TaiwanNewsData                `json:"TaiwanNews"`
}

type StockPriceData struct {
	StockPrice    StockPriceDetail `json:"StockPrice"`
	TechIndex     TechIndex        `json:"TechIndex"`
	MovingAverage MovingAverage    `json:"MovingAverage"`
	UpdateDate    string           `json:"update_date"`
}

type StockPriceDetail struct {
	Open          float64 `json:"open"`
	Close         float64 `json:"close"`
	Max           float64 `json:"max"`
	Min           float64 `json:"min"`
	TradingVolume int     `json:"Trading_Volume"`
	TradingMoney  int     `json:"Trading_money"`
	Spread        float64 `json:"spread"`
	SpreadPer     float64 `json:"spread_per"`
	StockID       string  `json:"stock_id"`
	StockName     string  `json:"stock_name"`
}

type TechIndex struct {
	Rsv  float64 `json:"rsv"`
	Bias float64 `json:"bias"`
}

type MovingAverage struct {
	Week     float64 `json:"week"`
	TwoWeek  float64 `json:"two_week"`
	Month    float64 `json:"month"`
	Period   float64 `json:"period"`
	HalfYear float64 `json:"half_year"`
	Year     float64 `json:"year"`
}

type TSAInstitutionalInvestorData struct {
	InstitutionalInvestor []TSAInstitutionalInvestorItem `json:"InstitutionalInvestor"`
	UpdateDate            string                         `json:"update_date"`
}

type TSAInstitutionalInvestorItem struct {
	Buy    int    `json:"buy"`
	Sell   int    `json:"sell"`
	Date   string `json:"date"`
	Name   string `json:"name"`
	ZhName string `json:"zh_name"`
	Spread int    `json:"spread"`
}

type ShareholdingData struct {
	Shareholding []ShareholdingItem `json:"Shareholding"`
	UpdateDate   string             `json:"update_date"`
}

type ShareholdingItem struct {
	Date                       string  `json:"date"`
	ForeignInvestmentSharesPer float64 `json:"ForeignInvestmentSharesPer"`
}

type MarginPurchaseShortSaleData struct {
	MarginPurchaseShortSale []MarginPurchaseShortSaleItem `json:"MarginPurchaseShortSale"`
	UpdateDate              string                        `json:"update_date"`
}

type MarginPurchaseShortSaleItem struct {
	ZhName string  `json:"zh_name"`
	Value  float64 `json:"value"`
}

type TaiwanFinancialStatementsData struct {
	TaiwanFinancialStatements []TaiwanFinancialStatementItem `json:"TaiwanFinancialStatements"`
	UpdateDate                string                         `json:"update_date"`
}

type TaiwanFinancialStatementItem struct {
	OriginName string  `json:"origin_name"`
	Value      float64 `json:"value"`
	ValueQoQ   float64 `json:"valueQoQ"`
	ValueYoY   float64 `json:"valueYoY"`
}

type TaiwanBalanceSheetData struct {
	TaiwanBalanceSheet []TaiwanBalanceSheetItem `json:"TaiwanBalanceSheet"`
	UpdateDate         string                   `json:"update_date"`
}

type TaiwanBalanceSheetItem struct {
	OriginName string  `json:"origin_name"`
	Value      float64 `json:"value"`
	ValueQoQ   float64 `json:"valueQoQ"`
	ValueYoY   float64 `json:"valueYoY"`
}

type TaiwanCashFlowsStatementData struct {
	TaiwanCashFlowsStatement []TaiwanCashFlowsStatementItem `json:"TaiwanCashFlowsStatement"`
	UpdateDate               string                         `json:"update_date"`
}

type TaiwanCashFlowsStatementItem struct {
	OriginName string  `json:"origin_name"`
	Value      float64 `json:"value"`
	ValueQoQ   float64 `json:"valueQoQ"`
	ValueYoY   float64 `json:"valueYoY"`
}

type TSATaiwanStockDividendData struct {
	StockDividend []StockDividendItem `json:"StockDividend"`
	CashDividend  []CashDividendItem  `json:"CashDividend"`
	UpdateDate    string              `json:"update_date"`
}

type StockDividendItem struct {
	Year                  int     `json:"year"`
	StockDividendDealDate string  `json:"StockDividendDealDate"`
	StockDividend         float64 `json:"StockDividend"`
}

type CashDividendItem struct {
	Year                    int     `json:"year"`
	CashDividendDealDate    string  `json:"CashDividendDealDate"`
	CashDividendReleaseDate string  `json:"CashDividendReleaseDate"`
	CashDividend            float64 `json:"CashDividend"`
}

type TaiwanNewsData struct {
	TaiwanNews []TaiwanNewsItem `json:"TaiwanNews"`
	UpdateDate string           `json:"update_date"`
}

type TaiwanNewsItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}
