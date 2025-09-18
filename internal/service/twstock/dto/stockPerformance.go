package twstock

type StockPerformanceResponseDto struct {
	Data      []StockPerformanceData `json:"data"`
	ChartData []byte                 `json:"chart_data,omitempty"`
}

type StockPerformanceData struct {
	StockID     string `json:"stock_id"`
	Period      string `json:"period"`
	PeriodName  string `json:"period_name"`
	Performance string `json:"performance"`
}
