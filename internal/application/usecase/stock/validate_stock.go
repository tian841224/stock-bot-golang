package stock

import (
	"context"
	"fmt"

	port "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type StockValidationUsecase interface {
	ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error)
}

type stockValidationUsecase struct {
	validation port.ValidationPort
	logger     logger.Logger
}

func NewStockValidationUsecase(
	validation port.ValidationPort,
	log logger.Logger,
) *stockValidationUsecase {
	return &stockValidationUsecase{
		validation: validation,
		logger:     log,
	}
}

func (uc *stockValidationUsecase) ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	stockSymbol, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil {
		uc.logger.Error("驗證股票代號失敗",
			logger.String("symbol", symbol),
			logger.Error(err))
		return nil, err
	}
	if stockSymbol == nil {
		err := fmt.Errorf("查無此股票代號，請重新確認")
		uc.logger.Warn("股票代號不存在",
			logger.String("symbol", symbol))
		return nil, err
	}
	return stockSymbol, nil
}
