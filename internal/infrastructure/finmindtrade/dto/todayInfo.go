package dto

// TodayInfoResponseDto 大盤資訊(法人/資券/美股大盤)
type TodayInfoResponseDto struct {
	Msg    string        `json:"msg"`
	Status int           `json:"status"`
	Data   TodayInfoData `json:"data"`
}

type TodayInfoData struct {
	InstitutionalInvestor        []InstitutionalInvestorData        `json:"InstitutionalInvestor"`
	TotalMarginPurchaseShortSale []TotalMarginPurchaseShortSaleData `json:"TotalMarginPurchaseShortSale"`
	USStockPrice                 []TodayInfoUSStockPriceData        `json:"USStockPrice"`
}

type InstitutionalInvestorData struct {
	Buy    int     `json:"buy"`
	Sell   int     `json:"sell"`
	Date   string  `json:"date"`
	Name   string  `json:"name"`
	ZhName string  `json:"zh_name"`
	Spread float64 `json:"spread"`
}

type TotalMarginPurchaseShortSaleData struct {
	TodayBalance int    `json:"TodayBalance"`
	YesBalance   int    `json:"YesBalance"`
	Buy          int    `json:"buy"`
	Sell         int    `json:"sell"`
	Date         string `json:"date"`
	Name         string `json:"name"`
	ZhName       string `json:"zh_name"`
}

type TodayInfoUSStockPriceData struct {
	AdjClose  float64 `json:"Adj_Close"`
	Close     float64 `json:"Close"`
	High      float64 `json:"High"`
	Low       float64 `json:"Low"`
	Open      float64 `json:"Open"`
	Spread    float64 `json:"Spread"`
	SpreadPer string  `json:"SpreadPer"`
	Date      string  `json:"date"`
	StockID   string  `json:"stock_id"`
	ZhName    string  `json:"zh_name"`
}
