package repository

import (
	"stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type SubscriptionSymbolRepository interface {
	Create(subscriptionSymbol *models.SubscriptionSymbol) error
	GetByID(id uint) (*models.SubscriptionSymbol, error)
	GetBySubscriptionID(subscriptionID uint) ([]*models.SubscriptionSymbol, error)
	GetBySymbolID(symbolID uint) ([]*models.SubscriptionSymbol, error)
	GetBySubscriptionAndSymbol(subscriptionID, symbolID uint) (*models.SubscriptionSymbol, error)
	Update(subscriptionSymbol *models.SubscriptionSymbol) error
	Delete(id uint) error
	DeleteBySubscriptionID(subscriptionID uint) error
	DeleteBySubscriptionAndSymbol(subscriptionID, symbolID uint) error
	List(offset, limit int) ([]*models.SubscriptionSymbol, error)
	BatchCreate(subscriptionSymbols []*models.SubscriptionSymbol) error
	GetSymbolsBySubscriptionID(subscriptionID uint) ([]*models.Symbol, error)
	GetSubscriptionsBySymbolID(symbolID uint) ([]*models.Subscription, error)
}

type subscriptionSymbolRepository struct {
	db *gorm.DB
}

func NewSubscriptionSymbolRepository(db *gorm.DB) SubscriptionSymbolRepository {
	return &subscriptionSymbolRepository{db: db}
}

// Create 建立新訂閱股票關聯
func (r *subscriptionSymbolRepository) Create(subscriptionSymbol *models.SubscriptionSymbol) error {
	return r.db.Create(subscriptionSymbol).Error
}

// GetByID 根據 ID 取得訂閱股票關聯
func (r *subscriptionSymbolRepository) GetByID(id uint) (*models.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.Preload("Subscription").Preload("Symbol").First(&subscriptionSymbol, id).Error
	if err != nil {
		return nil, err
	}
	return &subscriptionSymbol, nil
}

// GetBySubscriptionID 根據訂閱 ID 取得股票關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionID(subscriptionID uint) ([]*models.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.Preload("Symbol").Where("subscription_id = ?", subscriptionID).Find(&subscriptionSymbols).Error
	return subscriptionSymbols, err
}

// GetBySymbolID 根據股票 ID 取得訂閱關聯
func (r *subscriptionSymbolRepository) GetBySymbolID(symbolID uint) ([]*models.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.Preload("Subscription").Where("symbol_id = ?", symbolID).Find(&subscriptionSymbols).Error
	return subscriptionSymbols, err
}

// GetBySubscriptionAndSymbol 根據訂閱和股票取得關聯
func (r *subscriptionSymbolRepository) GetBySubscriptionAndSymbol(subscriptionID, symbolID uint) (*models.SubscriptionSymbol, error) {
	var subscriptionSymbol models.SubscriptionSymbol
	err := r.db.Where("subscription_id = ? AND symbol_id = ?", subscriptionID, symbolID).First(&subscriptionSymbol).Error
	if err != nil {
		return nil, err
	}
	return &subscriptionSymbol, nil
}

// Update 更新訂閱股票關聯
func (r *subscriptionSymbolRepository) Update(subscriptionSymbol *models.SubscriptionSymbol) error {
	return r.db.Save(subscriptionSymbol).Error
}

// Delete 刪除訂閱股票關聯
func (r *subscriptionSymbolRepository) Delete(id uint) error {
	return r.db.Delete(&models.SubscriptionSymbol{}, id).Error
}

// DeleteBySubscriptionID 根據訂閱 ID 刪除所有關聯
func (r *subscriptionSymbolRepository) DeleteBySubscriptionID(subscriptionID uint) error {
	return r.db.Where("subscription_id = ?", subscriptionID).Delete(&models.SubscriptionSymbol{}).Error
}

// DeleteBySubscriptionAndSymbol 根據訂閱和股票刪除關聯
func (r *subscriptionSymbolRepository) DeleteBySubscriptionAndSymbol(subscriptionID, symbolID uint) error {
	return r.db.Where("subscription_id = ? AND symbol_id = ?", subscriptionID, symbolID).Delete(&models.SubscriptionSymbol{}).Error
}

// List 取得訂閱股票關聯列表
func (r *subscriptionSymbolRepository) List(offset, limit int) ([]*models.SubscriptionSymbol, error) {
	var subscriptionSymbols []*models.SubscriptionSymbol
	err := r.db.Preload("Subscription").Preload("Symbol").Offset(offset).Limit(limit).Find(&subscriptionSymbols).Error
	return subscriptionSymbols, err
}

// BatchCreate 批次建立訂閱股票關聯
func (r *subscriptionSymbolRepository) BatchCreate(subscriptionSymbols []*models.SubscriptionSymbol) error {
	return r.db.CreateInBatches(subscriptionSymbols, 100).Error
}

// GetSymbolsBySubscriptionID 根據訂閱 ID 取得所有股票
func (r *subscriptionSymbolRepository) GetSymbolsBySubscriptionID(subscriptionID uint) ([]*models.Symbol, error) {
	var symbols []*models.Symbol
	err := r.db.Table("symbols").
		Joins("JOIN subscription_symbols ON symbols.id = subscription_symbols.symbol_id").
		Where("subscription_symbols.subscription_id = ?", subscriptionID).
		Find(&symbols).Error
	return symbols, err
}

// GetSubscriptionsBySymbolID 根據股票 ID 取得所有訂閱
func (r *subscriptionSymbolRepository) GetSubscriptionsBySymbolID(symbolID uint) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Table("subscriptions").
		Joins("JOIN subscription_symbols ON subscriptions.id = subscription_symbols.subscription_id").
		Where("subscription_symbols.symbol_id = ?", symbolID).
		Find(&subscriptions).Error
	return subscriptions, err
}
