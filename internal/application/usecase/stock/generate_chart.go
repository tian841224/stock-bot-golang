package stock

import (
	"context"
	"fmt"

	"github.com/tian841224/stock-bot/internal/application/dto"
	port "github.com/tian841224/stock-bot/internal/application/port"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type MarketChartUsecase interface {
	GetRevenueChart(ctx context.Context, symbol string) (*dto.RevenueChart, error)
	GetHistoricalCandlesChart(ctx context.Context, symbol string) (*dto.KlineCandlesChart, error)
	GetPerformanceChart(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error)
}

type marketChartUsecase struct {
	marketChart port.MarketChartPort
	validation  port.ValidationPort
	logger      logger.Logger
}

func NewMarketDataChartUsecase(
	marketChart port.MarketChartPort,
	validation port.ValidationPort,
	logger logger.Logger,
) *marketChartUsecase {
	return &marketChartUsecase{marketChart: marketChart, validation: validation, logger: logger}
}

// GetRevenueChart 取得股票營收圖表
func (uc *marketChartUsecase) GetRevenueChart(ctx context.Context, symbol string) (*dto.RevenueChart, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	chartBytes, err := uc.marketChart.GetRevenueChart(ctx, stock.Symbol)
	if err != nil {
		uc.logger.Error("取得營收圖表失敗", logger.Error(err))
		return nil, fmt.Errorf("取得營收圖表失敗: %w", err)
	}

	if chartBytes == nil {
		return nil, fmt.Errorf("取得營收圖表失敗:查無資料，請確認後再試")
	}

	return &dto.RevenueChart{
		ChartData: chartBytes,
		StockName: stock.Name,
	}, nil
}

// GetHistoricalCandlesChart 取得股票歷史K線圖
func (uc *marketChartUsecase) GetHistoricalCandlesChart(ctx context.Context, symbol string) (*dto.KlineCandlesChart, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	chartBytes, stockName, err := uc.marketChart.GetHistoricalCandlesChart(ctx, stock.Symbol)
	if err != nil {
		uc.logger.Error("取得歷史K線圖失敗", logger.Error(err))
		return nil, fmt.Errorf("取得歷史K線圖失敗: %w", err)
	}

	if chartBytes == nil {
		return nil, fmt.Errorf("取得歷史K線圖失敗:查無資料，請確認後再試")
	}

	return &dto.KlineCandlesChart{
		ChartData: chartBytes,
		StockName: stockName,
	}, nil
}

func (uc *marketChartUsecase) GetPerformanceChart(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error) {
	stock, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	chart, err := uc.marketChart.GetPerformanceChart(ctx, stock.Symbol)
	if err != nil {
		uc.logger.Error("取得績效圖表失敗", logger.Error(err))
		return nil, fmt.Errorf("查無資料，請確認後再試")
	}

	chartData := make([]dto.StockPerformanceData, len(chart.Data))
	for i, data := range chart.Data {
		chartData[i] = dto.StockPerformanceData{
			Period:      data.Period,
			PeriodName:  data.PeriodName,
			Performance: data.Performance,
		}
	}

	return &dto.StockPerformanceChart{
		Symbol:    stock.Symbol,
		StockName: stock.Name,
		Data:      chartData,
		ChartData: chart.ChartData,
	}, nil
}
