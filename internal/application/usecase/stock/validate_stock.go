package stock

import (
	"context"
	"fmt"

	port "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
)

type StockValidationUsecase interface {
	ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error)
}

type stockValidationUsecase struct {
	validation port.ValidationPort
}

func NewStockValidationUsecase(
	validation port.ValidationPort,
) *stockValidationUsecase {
	return &stockValidationUsecase{validation: validation}
}

func (uc *stockValidationUsecase) ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	stockSymbol, err := uc.validation.ValidateSymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if stockSymbol == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}
	return stockSymbol, nil
}
