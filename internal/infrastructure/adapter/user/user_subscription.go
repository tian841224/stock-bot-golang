package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type userSubscriptionGateway struct {
	subscriptionRepo             port.SubscriptionReader
	subscriptionSymbolRepo       port.SubscriptionSymbolReader
	stockSymbolRepo              port.StockSymbolReader
	subscriptionRepoWriter       port.SubscriptionWriter
	subscriptionSymbolRepoWriter port.SubscriptionSymbolWriter
	featureReader                port.FeatureReader
	userAccountPort              port.UserAccountPort
}

func NewUserSubscriptionGateway(
	subscriptionRepo port.SubscriptionReader,
	subscriptionSymbolRepo port.SubscriptionSymbolReader,
	stockSymbolRepo port.StockSymbolReader,
	subscriptionRepoWriter port.SubscriptionWriter,
	subscriptionSymbolRepoWriter port.SubscriptionSymbolWriter,
	featureReader port.FeatureReader,
	userAccountPort port.UserAccountPort,
) port.UserSubscriptionPort {
	return &userSubscriptionGateway{
		subscriptionRepo:             subscriptionRepo,
		subscriptionSymbolRepo:       subscriptionSymbolRepo,
		stockSymbolRepo:              stockSymbolRepo,
		subscriptionRepoWriter:       subscriptionRepoWriter,
		subscriptionSymbolRepoWriter: subscriptionSymbolRepoWriter,
		featureReader:                featureReader,
		userAccountPort:              userAccountPort,
	}
}

var _ port.UserSubscriptionPort = (*userSubscriptionGateway)(nil)

func (p *userSubscriptionGateway) GetUserSubscriptionItemList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	subscriptions, err := p.subscriptionRepo.GetUserSubscriptionList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("取得使用者訂閱項目列表失敗: %w", err)
	}

	var items []*dto.UserSubscriptionItem
	for _, sub := range subscriptions {
		items = append(items, &dto.UserSubscriptionItem{
			Item:   sub.Item,
			Status: sub.Active,
		})
	}
	return items, nil
}

func (p *userSubscriptionGateway) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
	subscriptionSymbols, err := p.subscriptionSymbolRepo.GetUserSubscriptionStockList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("取得使用者訂閱股票列表失敗: %w", err)
	}

	var stocks []*dto.UserSubscriptionStock
	for _, subSymbol := range subscriptionSymbols {
		if subSymbol.StockSymbol != nil && subSymbol.Subscription != nil {
			stocks = append(stocks, &dto.UserSubscriptionStock{
				Symbol: subSymbol.StockSymbol.Symbol,
				Name:   subSymbol.StockSymbol.Name,
				Status: subSymbol.Subscription.Active,
			})
		}
	}
	return stocks, nil
}

func (p *userSubscriptionGateway) GetUserSubscriptionDetail(ctx context.Context, userID uint) (*dto.UserSubscriptionDetail, error) {
	userSubscriptionItemList, err := p.GetUserSubscriptionItemList(ctx, userID)
	if err != nil {
		return nil, err
	}

	userSubscriptionStockList, err := p.GetUserSubscriptionStockList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserSubscriptionDetail{
		Items:  userSubscriptionItemList,
		Stocks: userSubscriptionStockList,
	}, nil
}

func (p *userSubscriptionGateway) AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error {
	// 將 SubscriptionType 轉換為 Feature Code（字串格式）
	featureCode := strconv.Itoa(int(item))

	// 檢查資料庫中是否存在對應的 Feature
	feature, err := p.featureReader.GetByCode(ctx, featureCode)
	if err != nil {
		return fmt.Errorf("查詢功能失敗: %w", err)
	}
	if feature == nil {
		return fmt.Errorf("查無對應的功能代碼: %s (訂閱類型: %d)", featureCode, item)
	}

	subscription := &entity.Subscription{
		UserID:    userID,
		Item:      item,
		FeatureID: feature.ID,
		Active:    true,
	}
	return p.subscriptionRepoWriter.Create(ctx, subscription)
}

func (p *userSubscriptionGateway) AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) error {
	// 取得股票符號 ID
	symbol, err := p.stockSymbolRepo.GetBySymbol(ctx, stockSymbol)
	if err != nil {
		return fmt.Errorf("取得股票符號失敗: %w", err)
	}
	if symbol == nil {
		return fmt.Errorf("查無此股票代號: %s", stockSymbol)
	}

	// 取得或建立股票訂閱
	subscription, err := p.subscriptionRepo.GetByUserAndFeature(ctx, userID, uint(valueobject.SubscriptionTypeStockInfo))
	if err != nil {
		return fmt.Errorf("取得訂閱失敗: %w", err)
	}

	if subscription == nil {
		return fmt.Errorf("請先訂閱股票資訊")
	}

	// 建立訂閱股票關聯
	subscriptionSymbol := &entity.SubscriptionSymbol{
		SubscriptionID: subscription.ID,
		SymbolID:       symbol.ID,
	}
	return p.subscriptionSymbolRepoWriter.Create(ctx, subscriptionSymbol)
}

func (p *userSubscriptionGateway) DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) error {
	// 取得股票符號 ID
	symbol, err := p.stockSymbolRepo.GetBySymbol(ctx, stockSymbol)
	if err != nil {
		return fmt.Errorf("取得股票符號失敗: %w", err)
	}
	if symbol == nil {
		return fmt.Errorf("查無此股票代號: %s", stockSymbol)
	}

	// 取得訂閱
	subscription, err := p.subscriptionRepo.GetByUserAndFeature(ctx, userID, uint(valueobject.SubscriptionTypeStockInfo))
	if err != nil {
		return fmt.Errorf("取得訂閱失敗: %w", err)
	}
	if subscription == nil {
		return fmt.Errorf("未找到訂閱")
	}

	// 取得訂閱股票關聯
	subscriptionSymbol, err := p.subscriptionSymbolRepo.GetBySubscriptionAndSymbol(ctx, subscription.ID, symbol.ID)
	if err != nil {
		return fmt.Errorf("取得訂閱股票關聯失敗: %w", err)
	}
	if subscriptionSymbol == nil {
		return fmt.Errorf("未找到訂閱股票")
	}

	// 刪除訂閱股票關聯
	return p.subscriptionSymbolRepoWriter.Delete(ctx, subscriptionSymbol.ID)
}

func (p *userSubscriptionGateway) DeleteUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error {
	subscription, err := p.subscriptionRepo.GetByUserAndFeature(ctx, userID, uint(item))
	if err != nil {
		return fmt.Errorf("取得訂閱失敗: %w", err)
	}
	if subscription == nil {
		return fmt.Errorf("未找到訂閱項目")
	}
	return p.subscriptionRepoWriter.Delete(ctx, subscription.ID)
}
