// Package dto 提供鉅亨網 API 的資料傳輸物件
package dto

// CnyesRevenueResponseDto 營收
type CnyesRevenueResponseDto struct {
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Data       CnyesRevenueDataDto `json:"data"`
}

type CnyesRevenueDataDto struct {
	Time     []int64                 `json:"time"`
	Code     string                  `json:"code"`
	Name     string                  `json:"name"`
	Datasets CnyesRevenueDatasetsDto `json:"datasets"`
}

type CnyesRevenueDatasetsDto struct {
	// 股價
	C []float64 `json:"c"`
	// 月營收
	SaleMonth []int64 `json:"saleMonth"`
	// 累積營收
	SaleAccumulated []int64 `json:"saleAccumulated"`
	// 月營收年增率
	YoY []float64 `json:"yoy"`
	// 累積營收年增率
	YoYAccumulated []float64 `json:"yoyAccumulated"`
}
