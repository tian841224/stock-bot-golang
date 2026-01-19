package port

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/dto"
)

// MarketDataPort 封裝 bot usecase 取用市場/股票資料所需的介面。
type MarketDataPort interface {
	// 取得大盤資訊
	GetDailyMarketInfo(ctx context.Context, count int) (*[]dto.DailyMarketInfo, error)
	// 取得股票績效
	GetStockPerformance(ctx context.Context, symbol string) ([]dto.StockPerformanceData, error)
	// 取得交易量排行
	GetTopVolumeStock(ctx context.Context) ([]*dto.TopVolume, error)
	// 取得股票價格
	GetStockPrice(ctx context.Context, symbol string, dates ...*time.Time) (*[]dto.StockPrice, error)
	// 取得股票公司資訊
	GetStockCompanyInfo(ctx context.Context, symbol string) (*dto.StockCompanyInfo, error)
	// 取得股票營收資料
	GetStockRevenue(ctx context.Context, symbol string) (*dto.StockRevenue, error)
	// 取得最後交易日
	GetLatestTradeDate(ctx context.Context) (time.Time, error)
	// 取得最後交易日 by 日期範圍
	GetLatestTradeDateByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) ([]time.Time, error)
	GetStockNews(ctx context.Context, symbol string) ([]dto.StockNews, error)
}
