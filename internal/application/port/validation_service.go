package port

import (
	"context"

	"github.com/tian841224/stock-bot/internal/domain/entity"
)

// ValidationPort 封裝 bot usecase 驗證股票代號所需的介面。
type ValidationPort interface {
	// 驗證股票代號
	ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error)
}
