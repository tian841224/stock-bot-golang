package stock

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// mockUserAccountPort 用於測試的 UserAccountPort mock
type mockUserAccountPort struct {
	GetOrCreateUserFunc func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
}

type mockUserSubscriptionPort struct {
	GetUserSubscriptionListFunc      func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	AddUserSubscriptionStockFunc     func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	AddUserSubscriptionItemFunc      func(ctx context.Context, userID uint, item valueobject.SubscriptionType) error
	DeleteUserSubscriptionStockFunc  func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	GetUserSubscriptionStockListFunc func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
	GetUserSubscriptionItemsFunc     func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	IsStockSubscribedFunc            func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	UpdateUserSubscriptionItemFunc   func(ctx context.Context, userID uint, item valueobject.SubscriptionType, status bool) (bool, error)
}

func (m *mockUserAccountPort) GetOrCreateUser(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
	if m.GetOrCreateUserFunc != nil {
		return m.GetOrCreateUserFunc(ctx, accountID, userType)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	if m.GetUserSubscriptionListFunc != nil {
		return m.GetUserSubscriptionListFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.AddUserSubscriptionStockFunc != nil {
		return m.AddUserSubscriptionStockFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.DeleteUserSubscriptionStockFunc != nil {
		return m.DeleteUserSubscriptionStockFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
	if m.GetUserSubscriptionStockListFunc != nil {
		return m.GetUserSubscriptionStockListFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionItems(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	if m.GetUserSubscriptionItemsFunc != nil {
		return m.GetUserSubscriptionItemsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) IsStockSubscribed(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.IsStockSubscribedFunc != nil {
		return m.IsStockSubscribedFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) UpdateUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType, status bool) (bool, error) {
	if m.UpdateUserSubscriptionItemFunc != nil {
		return m.UpdateUserSubscriptionItemFunc(ctx, userID, item, status)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error {
	if m.AddUserSubscriptionItemFunc != nil {
		return m.AddUserSubscriptionItemFunc(ctx, userID, item)
	}
	return nil
}

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

type mockTradeDateRepository struct {
	GetByIDFunc               func(ctx context.Context, id uint) (*entity.TradeDate, error)
	GetByDateFunc             func(ctx context.Context, date time.Time) (*entity.TradeDate, error)
	GetByDateRangeFunc        func(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error)
	CreateFunc                func(ctx context.Context, tradeDate *entity.TradeDate) error
	BatchCreateTradeDatesFunc func(ctx context.Context, tradeDates []*entity.TradeDate) error
}

func (m *mockTradeDateRepository) GetByID(ctx context.Context, id uint) (*entity.TradeDate, error) {
	if m != nil && m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) GetByDate(ctx context.Context, date time.Time) (*entity.TradeDate, error) {
	if m != nil && m.GetByDateFunc != nil {
		return m.GetByDateFunc(ctx, date)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error) {
	if m != nil && m.GetByDateRangeFunc != nil {
		return m.GetByDateRangeFunc(ctx, startDate, endDate)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) Create(ctx context.Context, tradeDate *entity.TradeDate) error {
	if m != nil && m.CreateFunc != nil {
		return m.CreateFunc(ctx, tradeDate)
	}
	return nil
}

func (m *mockTradeDateRepository) BatchCreateTradeDates(ctx context.Context, tradeDates []*entity.TradeDate) error {
	if m != nil && m.BatchCreateTradeDatesFunc != nil {
		return m.BatchCreateTradeDatesFunc(ctx, tradeDates)
	}
	return nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Field) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Debug(msg string, fields ...logger.Field) {}
func (m *mockLogger) Panic(msg string, fields ...logger.Field) {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Field) {}
func (m *mockLogger) Sync() error                              { return nil }
