package stock

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// mockValidationPort 用於測試的 ValidationPort mock
type mockValidationPort struct {
	ValidateSymbolFunc func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
}

func (m *mockValidationPort) ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	if m.ValidateSymbolFunc != nil {
		return m.ValidateSymbolFunc(ctx, symbol)
	}
	return nil, nil
}

type mockMarketDataPort struct {
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

func (m *mockMarketDataPort) GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error) {
	if m != nil && m.GetDailyMarketInfoFunc != nil {
		return m.GetDailyMarketInfoFunc(ctx, count)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetStockPerformance(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error) {
	if m != nil && m.GetStockPerformanceFunc != nil {
		return m.GetStockPerformanceFunc(ctx, symbol)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetTopVolumeStock(ctx context.Context) ([]*dto.TopVolume, error) {
	if m != nil && m.GetTopVolumeStockFunc != nil {
		return m.GetTopVolumeStockFunc(ctx)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetStockPrice(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error) {
	if m != nil && m.GetStockPriceFunc != nil {
		return m.GetStockPriceFunc(ctx, symbol, dates...)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error) {
	if m != nil && m.GetStockCompanyInfoFunc != nil {
		return m.GetStockCompanyInfoFunc(ctx, symbol)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetStockRevenue(ctx context.Context, symbol string) (*dto.StockRevenue, error) {
	if m != nil && m.GetStockRevenueFunc != nil {
		return m.GetStockRevenueFunc(ctx, symbol)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetLatestTradeDate(ctx context.Context) (time.Time, error) {
	if m != nil && m.GetLatestTradeDateFunc != nil {
		return m.GetLatestTradeDateFunc(ctx)
	}
	return time.Time{}, nil
}

func (m *mockMarketDataPort) GetLatestTradeDateByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error) {
	if m != nil && m.GetLatestTradeDateByDateRangeFunc != nil {
		return m.GetLatestTradeDateByDateRangeFunc(ctx, startDate, endDate)
	}
	return nil, nil
}

func (m *mockMarketDataPort) GetStockNews(ctx context.Context, symbol string) ([]dto.StockNews, error) {
	if m != nil && m.GetStockNewsFunc != nil {
		return m.GetStockNewsFunc(ctx, symbol)
	}
	return nil, nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Field)   {}
func (m *mockLogger) Error(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Warn(msg string, fields ...logger.Field)   {}
func (m *mockLogger) Debug(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Panic(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Sync() error                               { return nil }
func (m *mockLogger) With(fields ...logger.Field) logger.Logger { return m }
