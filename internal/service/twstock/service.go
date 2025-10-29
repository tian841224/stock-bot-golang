package twstock

import (
	"github.com/tian841224/stock-bot/internal/infrastructure/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/fugle"
	"github.com/tian841224/stock-bot/internal/infrastructure/twse"
	"github.com/tian841224/stock-bot/internal/repository"
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
