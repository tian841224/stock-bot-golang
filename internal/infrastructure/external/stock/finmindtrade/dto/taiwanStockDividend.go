package dto

// TaiwanStockDividendResponseDto 股利發放
type TaiwanStockDividendResponseDto struct {
	Msg    string                  `json:"msg"`
	Status int                     `json:"status"`
	Data   TaiwanStockDividendData `json:"data"`
}
type TaiwanStockDividendData struct {
	Date                                  string  `json:"date"`
	StockID                               string  `json:"stock_id"`
	Year                                  string  `json:"year"`
	StockEarningsDistribution             float64 `json:"StockEarningsDistribution"`
	StockStatutorySurplus                 float64 `json:"StockStatutorySurplus"`
	StockExDividendTradingDate            string  `json:"StockExDividendTradingDate"`
	TotalEmployeeStockDividend            float64 `json:"TotalEmployeeStockDividend"`
	TotalEmployeeStockDividendAmount      float64 `json:"TotalEmployeeStockDividendAmount"`
	RatioOfEmployeeStockDividendOfTotal   float64 `json:"RatioOfEmployeeStockDividendOfTotal"`
	RatioOfEmployeeStockDividend          float64 `json:"RatioOfEmployeeStockDividend"`
	CashEarningsDistribution              float64 `json:"CashEarningsDistribution"`
	CashStatutorySurplus                  float64 `json:"CashStatutorySurplus"`
	CashExDividendTradingDate             string  `json:"CashExDividendTradingDate"`
	CashDividendPaymentDate               string  `json:"CashDividendPaymentDate"`
	TotalEmployeeCashDividend             float64 `json:"TotalEmployeeCashDividend"`
	TotalNumberOfCashCapitalIncrease      float64 `json:"TotalNumberOfCashCapitalIncrease"`
	CashIncreaseSubscriptionRate          float64 `json:"CashIncreaseSubscriptionRate"`
	CashIncreaseSubscriptionpRrice        float64 `json:"CashIncreaseSubscriptionpRrice"`
	RemunerationOfDirectorsAndSupervisors float64 `json:"RemunerationOfDirectorsAndSupervisors"`
	ParticipateDistributionOfTotalShares  float64 `json:"ParticipateDistributionOfTotalShares"`
	AnnouncementDate                      string  `json:"AnnouncementDate"`
	AnnouncementTime                      string  `json:"AnnouncementTime"`
}
