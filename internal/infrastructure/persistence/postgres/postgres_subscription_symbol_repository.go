package repository

import (
	"context"
	"fmt"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	"github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"
	repo "github.com/tian841224/stock-bot/internal/application/port"

	"gorm.io/gorm"
)

type subscriptionSymbolRepository struct {
	db *gorm.DB
}

var _ repo.SubscriptionSymbolReader = (*subscriptionSymbolRepository)(nil)
var _ repo.SubscriptionSymbolWriter = (*subscriptionSymbolRepository)(nil)

func NewSubscriptionSymbolRepository(db *gorm.DB) *subscriptionSymbolRepository {
	return &subscriptionSymbolRepository{db: db}
}

func (r *subscriptionSymbolRepository) toEntity(model *models.SubscriptionSymbol) *entity.SubscriptionSymbol {
	return &entity.SubscriptionSymbol{
		ID:       model.ID,
		SymbolID: model.SymbolID,
		StockSymbol: &entity.StockSymbol{
			ID:     model.StockSymbol.ID,
			Symbol: model.StockSymbol.Symbol,
			Market: model.StockSymbol.Market,
			Name:   model.StockSymbol.Name,
		},
		SubscriptionID: model.SubscriptionID,
		Subscription: &entity.Subscription{
			ID:           model.Subscription.ID,
			UserID:       model.Subscription.UserID,
			Item:         valueobject.SubscriptionType(model.Subscription.FeatureID),
			Active:       model.Subscription.Status,
			ScheduleCron: model.Subscription.ScheduleCron,
			FeatureID:    model.Subscription.FeatureID,
		},
	}
}

func (r *subscriptionSymbolRepository) toModel(entity *entity.SubscriptionSymbol) *models.SubscriptionSymbol {
	return &models.SubscriptionSymbol{
		SymbolID: entity.SymbolID,
		StockSymbol: &models.StockSymbol{
			Symbol: entity.StockSymbol.Symbol,
			Market: entity.StockSymbol.Market,
			Name:   entity.StockSymbol.Name,
		},
		SubscriptionID: entity.SubscriptionID,
		Subscription: &models.Subscription{
			UserID:       entity.Subscription.UserID,
			FeatureID:    entity.Subscription.FeatureID,
			Status:       entity.Subscription.Active,
			ScheduleCron: entity.Subscription.ScheduleCron,
		},
	}
}

// GetByID 根據 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetByID(ctx context.Context, id uint) (*entity.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("Subscription").Preload("StockSymbol").First(&subscriptionSymbol, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&subscriptionSymbol), nil
}

// GetBySubscriptionID 根據訂閱 ID 取得股票關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("StockSymbol").Where("subscription_id = ?", subscriptionID).Find(&subscriptionSymbols).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var entities []*entity.SubscriptionSymbol
	for _, subscriptionSymbol := range subscriptionSymbols {
		entities = append(entities, r.toEntity(subscriptionSymbol))
	}
	return entities, nil
}

// GetBySymbolID 根據股票 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("Subscription").Preload("StockSymbol").Where("symbol_id = ?", symbolID).Find(&subscriptionSymbols).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var entities []*entity.SubscriptionSymbol
	for _, subscriptionSymbol := range subscriptionSymbols {
		entities = append(entities, r.toEntity(subscriptionSymbol))
	}
	return entities, nil
}

// GetUserSubscriptionStockList 取得使用者訂閱股票列表
func (r *subscriptionSymbolRepository) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*entity.SubscriptionSymbol, error) {
	var userSubscription *models.Subscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userSubscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var subscriptionSymbols []*models.SubscriptionSymbol
	err = r.db.WithContext(ctx).Preload("StockSymbol").Where("subscription_id = ? AND feature_id = ?", userSubscription.ID, valueobject.SubscriptionTypeStockInfo).Find(&subscriptionSymbols).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.SubscriptionSymbol
	for _, subscriptionSymbol := range subscriptionSymbols {
		entities = append(entities, r.toEntity(subscriptionSymbol))
	}
	return entities, nil
}

// GetBySubscriptionAndSymbol 根據訂閱和股票取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("Subscription").Preload("StockSymbol").Where("subscription_id = ? AND symbol_id = ?", subscriptionID, symbolID).First(&subscriptionSymbol).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&subscriptionSymbol), nil
}

// List 取得訂閱股票關聯列表
func (r *subscriptionSymbolRepository) List(ctx context.Context, offset, limit int) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("Subscription").Preload("StockSymbol").Offset(offset).Limit(limit).Find(&subscriptionSymbols).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var entities []*entity.SubscriptionSymbol
	for _, subscriptionSymbol := range subscriptionSymbols {
		entities = append(entities, r.toEntity(subscriptionSymbol))
	}
	return entities, nil
}

// GetAll 取得所有訂閱股票關聯
func (r *subscriptionSymbolRepository) GetAll(ctx context.Context, order string) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("Subscription").Preload("StockSymbol").Order(order).Find(&subscriptionSymbols).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var entities []*entity.SubscriptionSymbol
	for _, subscriptionSymbol := range subscriptionSymbols {
		entities = append(entities, r.toEntity(subscriptionSymbol))
	}
	return entities, nil
}

// GetSymbolsBySubscriptionID 根據訂閱 ID 取得所有股票
func (r *subscriptionSymbolRepository) GetSymbolsBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Table("stock_symbol").
		Joins("JOIN subscription_symbols ON stock_symbol.id = subscription_symbols.symbol_id").
		Where("subscription_symbols.subscription_id = ?", subscriptionID).
		Find(&symbols).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var entities []*entity.StockSymbol
	for _, symbol := range symbols {
		entities = append(entities, &entity.StockSymbol{
			ID:     symbol.ID,
			Symbol: symbol.Symbol,
			Market: symbol.Market,
			Name:   symbol.Name,
		})
	}
	return entities, nil
}

// GetSubscriptionsBySymbolID 根據股票 ID 取得所有訂閱
func (r *subscriptionSymbolRepository) GetSubscriptionsBySymbolID(ctx context.Context, symbolID uint) ([]*entity.Subscription, error) {
	var subscriptions []*entity.Subscription
	err := r.db.WithContext(ctx).Table("subscriptions").
		Joins("JOIN subscription_symbols ON subscriptions.id = subscription_symbols.subscription_id").
		Where("subscription_symbols.symbol_id = ?", symbolID).
		Find(&subscriptions).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, &entity.Subscription{
			ID:           subscription.ID,
			UserID:       subscription.UserID,
			Item:         valueobject.SubscriptionType(subscription.FeatureID),
			Active:       subscription.Active,
			ScheduleCron: subscription.ScheduleCron,
			FeatureID:    subscription.FeatureID,
		})
	}
	return entities, nil
}

// Create 建立新訂閱股票關聯
func (r *subscriptionSymbolRepository) Create(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error {
	return r.db.WithContext(ctx).Create(r.toModel(subscriptionSymbol)).Error
}

// BatchCreate 批次建立訂閱股票關聯
func (r *subscriptionSymbolRepository) BatchCreate(ctx context.Context, subscriptionSymbols []*entity.SubscriptionSymbol) error {

	models := make([]*models.SubscriptionSymbol, len(subscriptionSymbols))
	for i, subscriptionSymbol := range subscriptionSymbols {
		models[i] = r.toModel(subscriptionSymbol)
	}
	return r.db.WithContext(ctx).CreateInBatches(models, 100).Error
}

// Update 更新訂閱股票關聯
func (r *subscriptionSymbolRepository) Update(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error {

	dbModel := r.toModel(subscriptionSymbol)

	result := r.db.WithContext(ctx).Model(&models.SubscriptionSymbol{}).
		Where("id = ?", subscriptionSymbol.ID).
		Updates(dbModel)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription symbol not found with id: %d", subscriptionSymbol.ID)
	}

	return nil
}

// Delete 刪除訂閱股票關聯
func (r *subscriptionSymbolRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.SubscriptionSymbol{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription symbol not found with id: %d", id)
	}
	return r.db.Delete(&models.SubscriptionSymbol{}, id).Error
}

// DeleteBySubscriptionID 根據訂閱 ID 刪除所有關聯
func (r *subscriptionSymbolRepository) DeleteBySubscriptionID(ctx context.Context, subscriptionID uint) error {
	result := r.db.WithContext(ctx).Delete(&models.SubscriptionSymbol{}, "subscription_id = ?", subscriptionID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription symbol not found with subscription id: %d", subscriptionID)
	}

	return nil
}

// DeleteBySubscriptionAndSymbol 根據訂閱和股票刪除關聯
func (r *subscriptionSymbolRepository) DeleteBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) error {
	result := r.db.WithContext(ctx).Delete(&models.SubscriptionSymbol{}, "subscription_id = ? AND symbol_id = ?", subscriptionID, symbolID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription symbol not found with subscription id: %d and symbol id: %d", subscriptionID, symbolID)
	}
	return nil
}
