package mappers

import (
	"time"

	"stock-bot/internal/domain/stock"
	cnyesDto "stock-bot/internal/infrastructure/cnyes/dto"
	stockDto "stock-bot/internal/service/twstock/dto"
)

// RevenueMapper 營收資料轉換器
type RevenueMapper struct{}

// NewRevenueMapper 建立營收轉換器
func NewRevenueMapper() *RevenueMapper {
	return &RevenueMapper{}
}

// FromCnyesDto 從鉅亨網 DTO 轉換為領域模型
func (m *RevenueMapper) FromCnyesDto(dto cnyesDto.CnyesRevenueDataDto) *stock.Revenue {
	return &stock.Revenue{
		StockID:         dto.Code,
		Time:            dto.Time,
		StockPrice:      dto.Datasets.C,
		SaleMonth:       dto.Datasets.SaleMonth,
		SaleAccumulated: dto.Datasets.SaleAccumulated,
		YoY:             dto.Datasets.YoY,
		YoYAccumulated:  dto.Datasets.YoYAccumulated,
	}
}

// ToRevenueDto 轉換為營收 DTO
func (m *RevenueMapper) ToRevenueDto(revenue *stock.Revenue) *stockDto.RevenueDto {
	if revenue == nil {
		return nil
	}

	return &stockDto.RevenueDto{
		Time:            revenue.Time,
		Code:            revenue.StockID,
		Name:            "", // 需要從其他地方取得
		StockPrice:      revenue.StockPrice,
		SaleMonth:       revenue.SaleMonth,
		SaleAccumulated: revenue.SaleAccumulated,
		YoY:             revenue.YoY,
		YoYAccumulated:  revenue.YoYAccumulated,
	}
}

// ToChartData 轉換為圖表資料
func (m *RevenueMapper) ToChartData(revenue *stock.Revenue) []ChartData {
	if revenue == nil || len(revenue.Time) == 0 {
		return []ChartData{}
	}

	chartData := make([]ChartData, len(revenue.Time))

	// 取得最新的營收和年增率
	latestRevenue := int64(0)
	latestYoY := 0.0
	if len(revenue.SaleMonth) > 0 {
		latestRevenue = revenue.SaleMonth[len(revenue.SaleMonth)-1]
	}
	if len(revenue.YoY) > 0 {
		latestYoY = revenue.YoY[len(revenue.YoY)-1]
	}

	for i, timestamp := range revenue.Time {
		// 轉換時間戳記為日期格式
		t := time.Unix(timestamp, 0)
		period := t.Format("2006/01")
		periodName := t.Format("2006/01")

		// 取得對應的資料
		revenueValue := int64(0)
		yoy := 0.0
		stockPrice := 0.0

		if i < len(revenue.SaleMonth) {
			revenueValue = revenue.SaleMonth[i]
		}
		if i < len(revenue.YoY) {
			yoy = revenue.YoY[i]
		}
		if i < len(revenue.StockPrice) {
			stockPrice = revenue.StockPrice[i]
		}

		chartData[i] = ChartData{
			Period:        period,
			PeriodName:    periodName,
			Revenue:       revenueValue,
			YoY:           yoy,
			StockPrice:    stockPrice,
			LatestRevenue: latestRevenue,
			LatestYoY:     latestYoY,
		}
	}

	return chartData
}

// ChartData 圖表資料結構
type ChartData struct {
	Period        string
	PeriodName    string
	Revenue       int64
	YoY           float64
	StockPrice    float64
	LatestRevenue int64
	LatestYoY     float64
}
