package stock

import (
	"context"
	"fmt"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
)

type finmindStockInfoAdapter struct {
	finmindAPI *finmindtrade.FinmindTradeAPI
}

var _ port.StockInfoProvider = (*finmindStockInfoAdapter)(nil)

func NewFinmindStockInfoAdapter(finmindAPI *finmindtrade.FinmindTradeAPI) port.StockInfoProvider {
	return &finmindStockInfoAdapter{
		finmindAPI: finmindAPI,
	}
}

func (a *finmindStockInfoAdapter) GetTaiwanStockInfo(ctx context.Context) ([]*entity.StockSymbol, error) {
	response, err := a.finmindAPI.GetTaiwanStockInfo()
	if err != nil {
		return nil, fmt.Errorf("呼叫 FinMind API 失敗: %w", err)
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("FinMind API 回應錯誤，狀態碼: %d, 訊息: %s", response.Status, response.Msg)
	}

	symbols := make([]*entity.StockSymbol, 0, len(response.Data))
	for _, stockInfo := range response.Data {
		symbol := &entity.StockSymbol{
			Symbol: stockInfo.StockID,
			Name:   stockInfo.StockName,
			Market: "TW",
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (a *finmindStockInfoAdapter) GetUSStockInfo(ctx context.Context) ([]*entity.StockSymbol, error) {
	response, err := a.finmindAPI.GetUSStockInfo()
	if err != nil {
		return nil, fmt.Errorf("呼叫 FinMind API 失敗: %w", err)
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("FinMind API 回應錯誤，狀態碼: %d, 訊息: %s", response.Status, response.Msg)
	}

	symbols := make([]*entity.StockSymbol, 0, len(response.Data))
	for _, stockInfo := range response.Data {
		symbol := &entity.StockSymbol{
			Symbol: stockInfo.StockID,
			Name:   stockInfo.StockName,
			Market: "US",
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}
