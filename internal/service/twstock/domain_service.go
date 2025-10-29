package twstock

import (
	"github.com/tian841224/stock-bot/internal/domain/stock"
	"github.com/tian841224/stock-bot/internal/domain/stock/mappers"
	"github.com/tian841224/stock-bot/internal/domain/stock/services"
)

// DomainService 領域服務包裝器
type DomainService struct {
	stockMapper    *mappers.StockMapper
	revenueMapper  *mappers.RevenueMapper
	stockDomainSvc *services.StockDomainService
}

// NewDomainService 建立領域服務
func NewDomainService() *DomainService {
	return &DomainService{
		stockMapper:    mappers.NewStockMapper(),
		revenueMapper:  mappers.NewRevenueMapper(),
		stockDomainSvc: services.NewStockDomainService(),
	}
}

// GetStockMapper 取得股票轉換器
func (d *DomainService) GetStockMapper() *mappers.StockMapper {
	return d.stockMapper
}

// GetRevenueMapper 取得營收轉換器
func (d *DomainService) GetRevenueMapper() *mappers.RevenueMapper {
	return d.revenueMapper
}

// GetStockDomainService 取得股票領域服務
func (d *DomainService) GetStockDomainService() *services.StockDomainService {
	return d.stockDomainSvc
}

// ValidateStock 驗證股票資料
func (d *DomainService) ValidateStock(stock *stock.Stock) error {
	return d.stockDomainSvc.ValidateStockData(stock)
}

// FormatStockInfo 格式化股票資訊
func (d *DomainService) FormatStockInfo(stock *stock.Stock) map[string]interface{} {
	if stock == nil {
		return map[string]interface{}{}
	}

	info := map[string]interface{}{
		"stock_id":   stock.ID,
		"stock_name": stock.Name,
		"symbol":     stock.Symbol,
		"industry":   stock.Industry,
		"market":     stock.Market,
		"status":     stock.GetPriceChangeStatus(),
		"is_trading": stock.IsTradingDay(),
	}

	if stock.CurrentInfo != nil {
		info["current_price"] = d.stockDomainSvc.FormatPrice(stock.CurrentInfo.CurrentPrice)
		info["change"] = stock.CurrentInfo.Change
		info["change_rate"] = d.stockDomainSvc.FormatPercentage(stock.CurrentInfo.ChangeRate)
		info["volume"] = d.stockDomainSvc.FormatVolume(stock.CurrentInfo.Volume)
		info["turnover"] = stock.GetTurnoverInBillions()
	}

	if stock.Financials != nil {
		info["market_cap"] = stock.GetMarketCapInTrillions()
		info["pe"] = stock.Financials.PE
		info["pb"] = stock.Financials.PB
		info["eps"] = stock.Financials.EPS
		info["dividend_rate"] = stock.Financials.DividendRate
	}

	return info
}
