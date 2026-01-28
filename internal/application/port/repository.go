package port

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// UserRepository 定義使用者資料存取介面
type UserRepository interface {
	UserReader
	UserWriter
}

type UserReader interface {
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByAccountID(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
}

type UserWriter interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userID uint) error
}

// SubscriptionRepository 定義訂閱資料存取介面
type SubscriptionRepository interface {
	SubscriptionReader
	SubscriptionWriter
}

type SubscriptionReader interface {
	GetByID(ctx context.Context, id uint) (*entity.Subscription, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Subscription, error)
	GetByFeatureID(ctx context.Context, featureID uint) ([]*entity.Subscription, error)
	GetByUserAndFeature(ctx context.Context, userID, featureID uint) (*entity.Subscription, error)
	GetByStatus(ctx context.Context, status bool) ([]*entity.Subscription, error)
	List(ctx context.Context, offset, limit int) ([]*entity.Subscription, error)
	GetActiveSubscriptions(ctx context.Context) ([]*entity.Subscription, error)
	GetBySchedule(ctx context.Context, scheduleCron string) ([]*entity.Subscription, error)
	GetUserSubscriptionList(ctx context.Context, userID uint) ([]*entity.Subscription, error)
}

type SubscriptionWriter interface {
	Create(ctx context.Context, subscription *entity.Subscription) error
	Update(ctx context.Context, subscription *entity.Subscription) error
	UpdateStatus(ctx context.Context, id uint, status bool) error
	Delete(ctx context.Context, id uint) error
}

// StockSymbolRepository 定義股票代號資料存取介面
type StockSymbolRepository interface {
	StockSymbolReader
	StockSymbolWriter
}

type StockSymbolReader interface {
	GetByID(ctx context.Context, id uint) (*entity.StockSymbol, error)
	GetBySymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.StockSymbol, error)
	GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.StockSymbol, error)
	GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.StockSymbol, error)
	GetMarketStats(ctx context.Context) (map[string]int, error)
}

type StockSymbolWriter interface {
	Create(ctx context.Context, stockSymbol *entity.StockSymbol) error
	Update(ctx context.Context, stockSymbol *entity.StockSymbol) error
	Delete(ctx context.Context, id uint) error
	BatchUpsert(ctx context.Context, symbols []*entity.StockSymbol) (successCount, errorCount int, err error)
}

// SubscriptionSymbolRepository 定義訂閱股票資料存取介面
type SubscriptionSymbolRepository interface {
	SubscriptionSymbolReader
	SubscriptionSymbolWriter
}

type SubscriptionSymbolReader interface {
	GetByID(ctx context.Context, id uint) (*entity.SubscriptionSymbol, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.SubscriptionSymbol, error)
	GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.SubscriptionSymbol, error)
	GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.SubscriptionSymbol, error)
	GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*entity.SubscriptionSymbol, error)
	List(ctx context.Context, offset, limit int) ([]*entity.SubscriptionSymbol, error)
	GetAll(ctx context.Context, order string) ([]*entity.SubscriptionSymbol, error)
	GetByFeature(ctx context.Context, feature valueobject.SubscriptionType) ([]*entity.SubscriptionSymbol, error)
}

type SubscriptionSymbolWriter interface {
	Create(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error
	Update(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error
	Delete(ctx context.Context, id uint) error
}

// TradeDateRepository 定義交易日資料存取介面
type TradeDateRepository interface {
	TradeDateReader
	TradeDateWriter
}

type TradeDateReader interface {
	GetByID(ctx context.Context, id uint) (*entity.TradeDate, error)
	GetByDate(ctx context.Context, date time.Time) (*entity.TradeDate, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error)
}

type TradeDateWriter interface {
	Create(ctx context.Context, tradeDate *entity.TradeDate) error
	BatchCreateTradeDates(ctx context.Context, tradeDates []*entity.TradeDate) error
}

// FeatureRepository 定義功能資料存取介面
type FeatureRepository interface {
	FeatureReader
	FeatureWriter
}

type FeatureReader interface {
	GetByID(ctx context.Context, id uint) (*entity.Feature, error)
	GetByCode(ctx context.Context, code string) (*entity.Feature, error)
	GetByName(ctx context.Context, name string) (*entity.Feature, error)
	List(ctx context.Context, offset, limit int) ([]*entity.Feature, error)
}

type FeatureWriter interface {
	Create(ctx context.Context, feature *entity.Feature) error
	Update(ctx context.Context, feature *entity.Feature) error
	Delete(ctx context.Context, id uint) error
}

// StockInfoProvider 定義股票資訊提供者介面
type StockInfoProvider interface {
	GetTaiwanStockInfo(ctx context.Context) ([]*entity.StockSymbol, error)
	GetUSStockInfo(ctx context.Context) ([]*entity.StockSymbol, error)
	GetTaiwanStockTradingDate(ctx context.Context) ([]*entity.TradeDate, error)
}

// SyncMetadataRepository 定義同步元資料資料存取介面
type SyncMetadataRepository interface {
	GetByMarket(ctx context.Context, market string) (*entity.SyncMetadata, error)
	Upsert(ctx context.Context, metadata *entity.SyncMetadata) error
}
