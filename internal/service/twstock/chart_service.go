package twstock

import (
	"fmt"

	fugleDto "stock-bot/internal/infrastructure/fugle/dto"
	stockDto "stock-bot/internal/service/twstock/dto"
	"stock-bot/pkg/imageutil"
	"stock-bot/pkg/logger"

	"go.uber.org/zap"
)

// ========== 圖表生成相關方法 ==========

// GetStockPerformanceWithChart 取得股票績效並生成圖表
func (s *StockService) GetStockPerformanceWithChart(stockID string, chartType string) (*stockDto.StockPerformanceResponseDto, error) {
	logger.Log.Info("取得股票績效並生成圖表", zap.String("stockID", stockID), zap.String("chartType", chartType))

	// 先取得績效資料
	performanceData, err := s.GetStockPriceHistory(stockID)
	if err != nil {
		return nil, err
	}

	performanceResponse := &stockDto.StockPerformanceResponseDto{
		Data: performanceData,
	}

	// 取得股票名稱用於圖表標題
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	if err != nil || symbol == nil {
		logger.Log.Error("取得股票名稱失敗", zap.Error(err))
		return nil, fmt.Errorf("查無股票名稱")
	}

	// 轉換資料格式以供圖表使用
	chartData := make([]imageutil.PerformanceData, len(performanceData))
	for i, data := range performanceData {
		chartData[i] = imageutil.PerformanceData{
			Period:      data.Period,
			PeriodName:  data.PeriodName,
			Performance: data.Performance,
		}
	}

	// 生成圖表
	title := fmt.Sprintf("%s (%s) 績效表現", symbol.Name, stockID)
	var chartBytes []byte

	// 只支援折線圖
	chartBytes, err = imageutil.GeneratePerformanceLineChart(chartData, title)

	if err != nil {
		logger.Log.Error("生成圖表失敗", zap.Error(err))
		// 即使圖表生成失敗，仍然回傳績效資料
		return performanceResponse, nil
	}

	// 將圖表資料加入回應
	performanceResponse.ChartData = chartBytes

	return performanceResponse, nil
}

// GetStockHistoricalCandlesChart 取得股票歷史 K 線圖
func (s *StockService) GetStockHistoricalCandlesChart(dto fugleDto.FugleCandlesRequestDto) ([]byte, string, error) {
	response, err := s.fugleClient.GetStockHistoricalCandles(dto)
	if err != nil {
		return nil, "", err
	}
	if len(response.Data) == 0 {
		return nil, "", fmt.Errorf("查無K線資料")
	}

	// 取得股票名稱
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(dto.Symbol, "TW")
	stockName := dto.Symbol
	if err == nil && symbol != nil {
		stockName = symbol.Name
	}

	// 轉換資料
	chartData := make([]imageutil.CandlestickData, len(response.Data))
	for i, d := range response.Data {
		chartData[i] = imageutil.CandlestickData{
			Date:   d.Date,
			Open:   d.Open,
			High:   d.High,
			Low:    d.Low,
			Close:  d.Close,
			Volume: d.Volume,
		}
	}

	// 產生圖表
	chartBytes, err := imageutil.GenerateCandlestickChartPNG(chartData, stockName, symbol.Symbol)
	if err != nil {
		return nil, stockName, fmt.Errorf("產生K線圖失敗: %v", err)
	}

	return chartBytes, stockName, nil
}

// GetStockRevenueChart 取得股票營收圖表
func (s *StockService) GetStockRevenueChart(stockID string) ([]byte, error) {
	logger.Log.Info("產生股票營收圖表", zap.String("stockID", stockID))

	// 取得營收資料
	revenueData, err := s.GetStockRevenue(stockID)
	if err != nil {
		return nil, fmt.Errorf("取得營收資料失敗: %v", err)
	}

	// 轉換為圖表資料格式
	chartData := s.convertToChartData(revenueData)

	// 產生圖表
	chartBytes, err := imageutil.GenerateRevenueChartPNG(chartData, revenueData.Name, stockID)
	if err != nil {
		return nil, fmt.Errorf("產生營收圖表失敗: %v", err)
	}

	return chartBytes, nil
}
