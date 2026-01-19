package presenter

import (
	"context"
	"fmt"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
)

type validationGateway struct {
	validationPort         port.ValidationPort
	stockSymbolsRepository port.StockSymbolReader
}

func NewValidationGateway(validationPort port.ValidationPort, stockSymbolsRepository port.StockSymbolReader) *validationGateway {
	return &validationGateway{validationPort: validationPort, stockSymbolsRepository: stockSymbolsRepository}
}

var _ port.ValidationPort = (*validationGateway)(nil)

func (v *validationGateway) ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	stockSymbol, err := v.stockSymbolsRepository.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if stockSymbol == nil {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}
	return stockSymbol, nil
}
