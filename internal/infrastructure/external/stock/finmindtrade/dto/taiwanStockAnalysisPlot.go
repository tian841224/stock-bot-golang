package dto

type TaiwanStockAnalysisPlotResponseDto struct {
	Msg    string                      `json:"msg"`
	Status int                         `json:"status"`
	Data   TaiwanStockAnalysisPlotData `json:"data"`
}

type TaiwanStockAnalysisPlotData struct {
	EPS                TaiwanStockAnalysisPlotEPS                `json:"EPS"`
	TaiwanMonthRevenue TaiwanStockAnalysisPlotTaiwanMonthRevenue `json:"TaiwanMonthRevenue"`
}

type TaiwanStockAnalysisPlotEPS struct {
	Data       PlotData `json:"data"`
	Title      string   `json:"title"`
	YoY        float64  `json:"YoY"`
	QoQ        float64  `json:"QoQ"`
	UpdateDate string   `json:"update_date"`
}

type TaiwanStockAnalysisPlotTaiwanMonthRevenue struct {
	Data       PlotData `json:"data"`
	Title      string   `json:"title"`
	YoY        float64  `json:"YoY"`
	MoM        float64  `json:"MoM"`
	UpdateDate string   `json:"update_date"`
}

type PlotData struct {
	Labels []string    `json:"labels"`
	Series [][]float64 `json:"series"`
}
