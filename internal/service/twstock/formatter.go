package twstock

import (
	"time"

	cnyesInfraDto "github.com/tian841224/stock-bot/internal/infrastructure/cnyes/dto"
	stockDto "github.com/tian841224/stock-bot/internal/service/twstock/dto"
	"github.com/tian841224/stock-bot/pkg/imageutil"
)

// ========== 轉換格式相關方法 ==========

// formatStockInfo 格式化股票資訊（已重構為使用領域模型）
func (s *StockService) formatStockInfo(data cnyesInfraDto.CnyesStockQuoteDataDto) *stockDto.StockQuoteInfo {
	// 使用領域服務進行轉換
	stock, err := s.domainService.GetStockMapper().FromCnyesDto(data)
	if err != nil {
		// 如果轉換失敗，回傳基本資訊
		return &stockDto.StockQuoteInfo{
			StockID:   data.StockID,
			StockName: data.StockName,
			Industry:  data.Industry,
			Market:    data.Market,
		}
	}

	// 驗證股票資料
	if err := s.domainService.ValidateStock(stock); err != nil {
		// 驗證失敗時回傳基本資訊
		return &stockDto.StockQuoteInfo{
			StockID:   data.StockID,
			StockName: data.StockName,
			Industry:  data.Industry,
			Market:    data.Market,
		}
	}

	// 轉換為 DTO
	return s.domainService.GetStockMapper().ToStockQuoteDto(stock)
}

// formatRevenue 格式化財報資料
func (s *StockService) formatRevenue(data cnyesInfraDto.CnyesRevenueDataDto) *stockDto.RevenueDto {
	return &stockDto.RevenueDto{
		Time:            data.Time,
		Code:            data.Code,
		Name:            data.Name,
		StockPrice:      data.Datasets.C,
		SaleMonth:       data.Datasets.SaleMonth,
		SaleAccumulated: data.Datasets.SaleAccumulated,
		YoY:             data.Datasets.YoY,
		YoYAccumulated:  data.Datasets.YoYAccumulated,
	}
}

// formatStockQuote 格式化股票報價資料
func (s *StockService) formatStockQuote(data cnyesInfraDto.CnyesStockQuoteDataDto) *stockDto.StockQuoteInfo {
	return &stockDto.StockQuoteInfo{
		// 基本資訊
		StockID:   data.StockID,
		StockName: data.StockName,
		Industry:  data.Industry,
		Market:    data.Market,

		// 價格資訊
		CurrentPrice: data.CurrentPrice,
		Change:       data.Change,
		ChangeRate:   data.ChangeRate,
		OpenPrice:    data.OpenPrice,
		HighPrice:    data.HighPrice,
		LowPrice:     data.LowPrice,
		PrevClose:    data.PrevClose,

		// 成交量資訊 (轉換單位)
		Volume:      int64(data.Volume),
		Turnover:    data.Turnover / 1e8,    // 轉換為億元
		VolumeRatio: data.VolumeRatio * 100, // 轉換為百分比
		Amplitude:   data.Amplitude,

		// 財務指標
		PE:           data.PE,
		PB:           data.PB,
		MarketCap:    data.MarketCap / 1e12, // 轉換為兆元
		BookValue:    data.BookValue,
		EPS:          data.EPS,
		QuarterEPS:   data.QuarterEPS,
		Dividend:     data.Dividend,
		DividendRate: data.DividendRate,
		GrossMargin:  data.GrossMargin,
		OperMargin:   data.OperMargin,
		NetMargin:    data.NetMargin,

		// 價位區間
		UpperLimit:  data.UpperLimit,
		LowerLimit:  data.LowerLimit,
		High52W:     data.High52W,
		Low52W:      data.Low52W,
		High52WDate: data.High52WDate,
		Low52WDate:  data.Low52WDate,

		// 五檔資訊
		BidPrices: []float64{
			data.BidPrice1, data.BidPrice2, data.BidPrice3, data.BidPrice4, data.BidPrice5,
		},
		AskPrices: []float64{
			data.AskPrice1, data.AskPrice2, data.AskPrice3, data.AskPrice4, data.AskPrice5,
		},

		// 內外盤資訊
		OutVolume: int64(data.OutVolume),
		InVolume:  int64(data.InVolume),
		OutRatio:  data.OutRatio,
		InRatio:   data.InRatio,
	}
}

// convertToChartData 轉換營收資料為圖表格式
func (s *StockService) convertToChartData(revenueData *stockDto.RevenueDto) []imageutil.RevenueChartData {
	if revenueData == nil || len(revenueData.Time) == 0 {
		return []imageutil.RevenueChartData{}
	}

	chartData := make([]imageutil.RevenueChartData, len(revenueData.Time))

	// 取得最新的營收和年增率
	latestRevenue := int64(0)
	latestYoY := 0.0
	if len(revenueData.SaleMonth) > 0 {
		latestRevenue = revenueData.SaleMonth[len(revenueData.SaleMonth)-1]
	}
	if len(revenueData.YoY) > 0 {
		latestYoY = revenueData.YoY[len(revenueData.YoY)-1]
	}

	for i, timestamp := range revenueData.Time {
		// 轉換時間戳記為日期格式
		t := time.Unix(timestamp, 0)
		period := t.Format("2006/01")
		periodName := t.Format("2006/01")

		// 取得對應的資料
		revenue := int64(0)
		yoy := 0.0
		stockPrice := 0.0

		if i < len(revenueData.SaleMonth) {
			revenue = revenueData.SaleMonth[i]
		}
		if i < len(revenueData.YoY) {
			yoy = revenueData.YoY[i]
		}
		if i < len(revenueData.StockPrice) {
			stockPrice = revenueData.StockPrice[i]
		}

		chartData[i] = imageutil.RevenueChartData{
			Period:        period,
			PeriodName:    periodName,
			Revenue:       revenue,
			YoY:           yoy,
			StockPrice:    stockPrice,
			LatestRevenue: latestRevenue,
			LatestYoY:     latestYoY,
		}
	}

	return chartData
}
