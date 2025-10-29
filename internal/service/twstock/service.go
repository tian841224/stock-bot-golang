package twstock

import (
	"stock-bot/internal/infrastructure/cnyes"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/infrastructure/fugle"
	"stock-bot/internal/infrastructure/twse"
	"stock-bot/internal/repository"
)

// StockService 股票服務
type StockService struct {
	finmindClient finmindtrade.FinmindTradeAPIInterface
	twseAPI       *twse.TwseAPI
	cnyesAPI      *cnyes.CnyesAPI
	fugleClient   *fugle.FugleAPI
	symbolsRepo   repository.SymbolRepository
	domainService *DomainService
}

// NewStockService 建立股票服務實例
func NewStockService(
	finmindClient finmindtrade.FinmindTradeAPIInterface,
	twseAPI *twse.TwseAPI,
	cnyesAPI *cnyes.CnyesAPI,
	fugleClient *fugle.FugleAPI,
	symbolsRepo repository.SymbolRepository,
) *StockService {
	return &StockService{
		finmindClient: finmindClient,
		twseAPI:       twseAPI,
		cnyesAPI:      cnyesAPI,
		fugleClient:   fugleClient,
		symbolsRepo:   symbolsRepo,
		domainService: NewDomainService(),
	}
}
