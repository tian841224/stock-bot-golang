package repository

import (
	"context"
	"fmt"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type subscriptionSymbolRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.SubscriptionSymbolReader = (*subscriptionSymbolRepository)(nil)
var _ repo.SubscriptionSymbolWriter = (*subscriptionSymbolRepository)(nil)

func NewSubscriptionSymbolRepository(db *gorm.DB, log logger.Logger) *subscriptionSymbolRepository {
	return &subscriptionSymbolRepository{
		db:     db,
		logger: log,
	}
}

func (r *subscriptionSymbolRepository) toEntity(model *models.SubscriptionSymbol) *entity.SubscriptionSymbol {
	result := &entity.SubscriptionSymbol{
		ID:       model.ID,
		UserID:   model.UserID,
		SymbolID: model.SymbolID,
	}

	if model.StockSymbol != nil {
		result.StockSymbol = &entity.StockSymbol{
			ID:     model.StockSymbol.ID,
			Symbol: model.StockSymbol.Symbol,
			Market: model.StockSymbol.Market,
			Name:   model.StockSymbol.Name,
		}
	}

	if model.User != nil {
		result.User = &entity.User{
			ID:        model.User.ID,
			AccountID: model.User.AccountID,
			UserType:  model.User.UserType,
			Status:    model.User.Status,
		}
	}

	return result
}

func (r *subscriptionSymbolRepository) toModel(entity *entity.SubscriptionSymbol) *models.SubscriptionSymbol {
	model := &models.SubscriptionSymbol{
		UserID:   entity.UserID,
		SymbolID: entity.SymbolID,
	}

	if entity.StockSymbol != nil {
		model.StockSymbol = &models.StockSymbol{
			Symbol: entity.StockSymbol.Symbol,
			Market: entity.StockSymbol.Market,
			Name:   entity.StockSymbol.Name,
		}
	}

	return model
}

// GetByID 根據 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetByID(ctx context.Context, id uint) (*entity.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.WithContext(ctx).First(&subscriptionSymbol, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&subscriptionSymbol), nil
}

// GetBySymbolID 根據股票 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("StockSymbol").Where("symbol_id = ?", symbolID).Find(&subscriptionSymbols).Error

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

// GetBySubscriptionID 根據訂閱 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).
		Table("subscription_symbols").
		Preload("StockSymbol").
		Joins("JOIN subscriptions ON subscriptions.user_id = subscription_symbols.user_id").
		Where("subscriptions.id = ?", subscriptionID).
		Find(&subscriptionSymbols).Error

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

// GetBySubscriptionAndSymbol 根據訂閱 ID 和股票 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.WithContext(ctx).
		Table("subscription_symbols").
		Preload("StockSymbol").
		Joins("JOIN subscriptions ON subscriptions.user_id = subscription_symbols.user_id").
		Where("subscriptions.id = ? AND subscription_symbols.symbol_id = ?", subscriptionID, symbolID).
		First(&subscriptionSymbol).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&subscriptionSymbol), nil
}

// GetUserSubscriptionStockList 取得使用者訂閱股票列表
func (r *subscriptionSymbolRepository) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Preload("StockSymbol").Where("user_id = ?", userID).Find(&subscriptionSymbols).Error
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

// GetUserSubscribedSymbols 取得使用者訂閱的所有股票
func (r *subscriptionSymbolRepository) GetUserSubscribedSymbols(ctx context.Context, userID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Table("stock_symbol").
		Joins("JOIN subscription_symbols ON stock_symbol.id = subscription_symbols.symbol_id").
		Where("subscription_symbols.user_id = ?", userID).
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

// GetUserSubscribedSymbolsByFeature 根據使用者和功能取得訂閱股票
func (r *subscriptionSymbolRepository) GetUserSubscribedSymbolsByFeature(ctx context.Context, userID, featureID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Table("stock_symbol").
		Joins("JOIN subscription_symbols ON stock_symbol.id = subscription_symbols.symbol_id").
		Where("subscription_symbols.user_id = ?", userID).
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

// List 取得訂閱股票關聯列表
func (r *subscriptionSymbolRepository) List(ctx context.Context, offset, limit int) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&subscriptionSymbols).Error

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
	err := r.db.WithContext(ctx).Order(order).Find(&subscriptionSymbols).Error

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

// GetAllWithDetails 取得所有訂閱股票關聯 (包含關聯資料)
func (r *subscriptionSymbolRepository) GetAllWithDetails(ctx context.Context) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	// Preload StockSymbol and User to avoid N+1 queries
	err := r.db.WithContext(ctx).
		Preload("StockSymbol").
		Preload("User").
		Find(&subscriptionSymbols).Error

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

// GetByFeature 取得所有訂閱特定功能的使用者及其關注的股票
func (r *subscriptionSymbolRepository) GetByFeature(ctx context.Context, feature valueobject.SubscriptionType) ([]*entity.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol

	err := r.db.WithContext(ctx).
		Preload("StockSymbol").
		Preload("User").
		Table("subscription_symbols").
		Joins("JOIN subscriptions ON subscriptions.user_id = subscription_symbols.user_id").
		Where("subscriptions.feature_id = ? AND subscriptions.status = ?", feature, true).
		Find(&subscriptionSymbols).Error

	if err != nil {
		return nil, err
	}

	var entities []*entity.SubscriptionSymbol
	for _, s := range subscriptionSymbols {
		entities = append(entities, r.toEntity(s))
	}
	return entities, nil
}

// Create 建立新訂閱股票關聯
func (r *subscriptionSymbolRepository) Create(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error {
	r.logger.Info("Creating subscription symbol", logger.Any("user_id", subscriptionSymbol.UserID), logger.Any("symbol_id", subscriptionSymbol.SymbolID))

	err := r.db.WithContext(ctx).Create(r.toModel(subscriptionSymbol)).Error
	if err != nil {
		r.logger.Error("Failed to create subscription symbol", logger.Error(err), logger.Any("user_id", subscriptionSymbol.UserID), logger.Any("symbol_id", subscriptionSymbol.SymbolID))
		return err
	}

	r.logger.Info("Subscription symbol created successfully", logger.Any("user_id", subscriptionSymbol.UserID), logger.Any("symbol_id", subscriptionSymbol.SymbolID))
	return nil
}

// BatchCreate 批次建立訂閱股票關聯
func (r *subscriptionSymbolRepository) BatchCreate(ctx context.Context, subscriptionSymbols []*entity.SubscriptionSymbol) error {
	r.logger.Info("Batch creating subscription symbols", logger.Int("count", len(subscriptionSymbols)))

	models := make([]*models.SubscriptionSymbol, len(subscriptionSymbols))
	for i, subscriptionSymbol := range subscriptionSymbols {
		models[i] = r.toModel(subscriptionSymbol)
	}

	err := r.db.WithContext(ctx).CreateInBatches(models, 100).Error
	if err != nil {
		r.logger.Error("Failed to batch create subscription symbols", logger.Error(err), logger.Int("count", len(subscriptionSymbols)))
		return err
	}

	r.logger.Info("Batch create subscription symbols completed", logger.Int("count", len(subscriptionSymbols)))
	return nil
}

// Update 更新訂閱股票關聯
func (r *subscriptionSymbolRepository) Update(ctx context.Context, subscriptionSymbol *entity.SubscriptionSymbol) error {
	r.logger.Info("Updating subscription symbol", logger.Any("id", subscriptionSymbol.ID))

	dbModel := r.toModel(subscriptionSymbol)

	result := r.db.WithContext(ctx).Model(&models.SubscriptionSymbol{}).
		Where("id = ?", subscriptionSymbol.ID).
		Updates(dbModel)

	if result.Error != nil {
		r.logger.Error("Failed to update subscription symbol", logger.Error(result.Error), logger.Any("id", subscriptionSymbol.ID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Subscription symbol not found for update", logger.Any("id", subscriptionSymbol.ID))
		return fmt.Errorf("subscription symbol not found with id: %d", subscriptionSymbol.ID)
	}

	r.logger.Info("Subscription symbol updated successfully", logger.Any("id", subscriptionSymbol.ID))
	return nil
}

// Delete 刪除訂閱股票關聯
func (r *subscriptionSymbolRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting subscription symbol", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.SubscriptionSymbol{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete subscription symbol", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Subscription symbol not found for deletion", logger.Any("id", id))
		return fmt.Errorf("subscription symbol not found with id: %d", id)
	}

	r.logger.Info("Subscription symbol deleted successfully", logger.Any("id", id))
	return nil
}
