package bot

import (
	"context"
	"errors"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type mockMarketDataUsecase struct {
	GetDailyMarketInfoFunc            func(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error)
	GetStockPerformanceFunc           func(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error)
	GetTopVolumeStockFunc             func(ctx context.Context) ([]*dto.TopVolume, error)
	GetStockPriceFunc                 func(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error)
	GetStockCompanyInfoFunc           func(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error)
	GetStockRevenueFunc               func(ctx context.Context, symbol string) (*dto.StockRevenue, error)
	GetLatestTradeDateFunc            func(ctx context.Context) (time.Time, error)
	GetLatestTradeDateByDateRangeFunc func(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error)
	GetStockNewsFunc                  func(ctx context.Context, symbol string) ([]dto.StockNews, error)
}

type mockMarketChartUsecase struct {
	GetRevenueChartFunc           func(ctx context.Context, symbol string) (*dto.RevenueChart, error)
	GetHistoricalCandlesChartFunc func(ctx context.Context, symbol string) (*dto.KlineCandlesChart, error)
	GetPerformanceChartFunc       func(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error)
}

type mockFormatterPort struct {
	FormatDailyMarketInfoFunc  func(marketData *[]dto.DailyMarketInfo, userType valueobject.UserType) string
	FormatStockPerformanceFunc func(stockPerformance *dto.StockPerformanceData, userType valueobject.UserType) string
	FormatStockPriceFunc       func(stockPrice *dto.StockPrice, userType valueobject.UserType) string
	FormatStockRevenueFunc     func(stockRevenue *dto.StockRevenue, userType valueobject.UserType) string
	FormatStockCompanyInfoFunc func(stockCompanyInfo *dto.StockCompanyInfo, userType valueobject.UserType) string
	FormatTopVolumeStockFunc   func(topVolume *[]dto.TopVolume, userType valueobject.UserType) string
}

func (m *mockMarketDataUsecase) GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {
	if m.GetDailyMarketInfoFunc != nil {
		return m.GetDailyMarketInfoFunc(ctx, count)
	}
	return nil, errors.New("GetDailyMarketInfoFunc is not implemented")
}

func (m *mockMarketDataUsecase) GetStockPerformance(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error) {
	if m.GetStockPerformanceFunc != nil {
		return m.GetStockPerformanceFunc(ctx, symbol)
	}
	return nil, errors.New("GetStockPerformanceFunc is not implemented")
}

func (m *mockMarketDataUsecase) GetTopVolumeStock(ctx context.Context) ([]*dto.TopVolume, error) {
	if m.GetTopVolumeStockFunc != nil {
		return m.GetTopVolumeStockFunc(ctx)
	}
	return nil, errors.New("GetTopVolumeStockFunc is not implemented")
}

func (m *mockMarketDataUsecase) GetStockPrice(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error) {
	if m.GetStockPriceFunc != nil {
		return m.GetStockPriceFunc(ctx, symbol, dates...)
	}
	return nil, errors.New("GetStockPriceFunc is not implemented")
}

func (m *mockMarketDataUsecase) GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error) {
	if m.GetStockCompanyInfoFunc != nil {
		return m.GetStockCompanyInfoFunc(ctx, symbol)
	}
	return nil, errors.New("GetStockCompanyInfoFunc is not implemented")
}

func (m *mockMarketDataUsecase) GetStockRevenue(ctx context.Context, symbol string) (*dto.StockRevenue, error) {
	if m.GetStockRevenueFunc != nil {
		return m.GetStockRevenueFunc(ctx, symbol)
	}
	return nil, errors.New("GetStockRevenueFunc is not implemented")
}
