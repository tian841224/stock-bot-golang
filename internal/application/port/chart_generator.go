package port

import (
	"context"

	dto "github.com/tian841224/stock-bot/internal/application/dto"
)

// MarketChartPort 封裝 bot usecase 取用市場/股票圖表所需的介面。
type MarketChartPort interface {
	// 取得股票K線圖
	GetHistoricalCandlesChart(ctx context.Context, symbol string) ([]byte, string, error)
	// 取得股票營收圖表
	GetRevenueChart(ctx context.Context, symbol string) ([]byte, error)
	// 取得股票績效圖表
	GetPerformanceChart(ctx context.Context, symbol string) (*dto.StockPerformanceChart, error)
}
