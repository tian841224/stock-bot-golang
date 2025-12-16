package port

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type BotCommandPort interface {
	GetUseGuideMessage() string
	GetDailyMarketInfo(ctx context.Context, userType valueobject.UserType, count int) (string, error)
	GetStockPerformance(ctx context.Context, userType valueobject.UserType, symbol string) (string, error)
	GetStockPerformanceChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetTopVolumeStock(ctx context.Context, userType valueobject.UserType) (string, error)
	GetStockPrice(ctx context.Context, userType valueobject.UserType, symbol string, date *time.Time) (string, error)
	GetStockRevenueChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetHistoricalCandlesChart(ctx context.Context, symbol string) (*dto.ChartAsset, error)
	GetStockCompanyInfo(ctx context.Context, userType valueobject.UserType, symbol string) (string, error)
	GetStockNewsForLine(ctx context.Context, symbol string) (*dto.LineStockNewsMessage, error)
	GetStockNewsForTelegram(ctx context.Context, symbol string) (*dto.TgStockNewsMessage, error)
}
