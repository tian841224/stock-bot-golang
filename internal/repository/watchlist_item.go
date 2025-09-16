package repository

import (
	"stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type WatchlistItemRepository interface {
	Create(item *models.WatchlistItem) error
	GetByID(id uint) (*models.WatchlistItem, error)
	GetByWatchlistID(watchlistID uint) ([]*models.WatchlistItem, error)
	GetBySymbolID(symbolID uint) ([]*models.WatchlistItem, error)
	GetByWatchlistAndSymbol(watchlistID, symbolID uint) (*models.WatchlistItem, error)
	Update(item *models.WatchlistItem) error
	Delete(id uint) error
	DeleteByWatchlistID(watchlistID uint) error
	DeleteByWatchlistAndSymbol(watchlistID, symbolID uint) error
	List(offset, limit int) ([]*models.WatchlistItem, error)
	BatchCreate(items []*models.WatchlistItem) error
	GetSymbolsByWatchlistID(watchlistID uint) ([]*models.Symbol, error)
	IsSymbolInWatchlist(watchlistID, symbolID uint) (bool, error)
}

type watchlistItemRepository struct {
	db *gorm.DB
}

func NewWatchlistItemRepository(db *gorm.DB) WatchlistItemRepository {
	return &watchlistItemRepository{db: db}
}

// Create 建立新觀察清單項目
func (r *watchlistItemRepository) Create(item *models.WatchlistItem) error {
	return r.db.Create(item).Error
}

// GetByID 根據 ID 取得觀察清單項目
func (r *watchlistItemRepository) GetByID(id uint) (*models.WatchlistItem, error) {
	var item models.WatchlistItem
	err := r.db.Preload("Watchlist").Preload("Symbol").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByWatchlistID 根據觀察清單 ID 取得項目
func (r *watchlistItemRepository) GetByWatchlistID(watchlistID uint) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	err := r.db.Preload("Symbol").Where("watchlist_id = ?", watchlistID).Find(&items).Error
	return items, err
}

// GetBySymbolID 根據股票 ID 取得觀察清單項目
func (r *watchlistItemRepository) GetBySymbolID(symbolID uint) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	err := r.db.Preload("Watchlist").Where("symbol_id = ?", symbolID).Find(&items).Error
	return items, err
}

// GetByWatchlistAndSymbol 根據觀察清單和股票取得項目
func (r *watchlistItemRepository) GetByWatchlistAndSymbol(watchlistID, symbolID uint) (*models.WatchlistItem, error) {
	var item models.WatchlistItem
	err := r.db.Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新觀察清單項目
func (r *watchlistItemRepository) Update(item *models.WatchlistItem) error {
	return r.db.Save(item).Error
}

// Delete 刪除觀察清單項目
func (r *watchlistItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.WatchlistItem{}, id).Error
}

// DeleteByWatchlistID 根據觀察清單 ID 刪除所有項目
func (r *watchlistItemRepository) DeleteByWatchlistID(watchlistID uint) error {
	return r.db.Where("watchlist_id = ?", watchlistID).Delete(&models.WatchlistItem{}).Error
}

// DeleteByWatchlistAndSymbol 根據觀察清單和股票刪除項目
func (r *watchlistItemRepository) DeleteByWatchlistAndSymbol(watchlistID, symbolID uint) error {
	return r.db.Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).Delete(&models.WatchlistItem{}).Error
}

// List 取得觀察清單項目列表
func (r *watchlistItemRepository) List(offset, limit int) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	err := r.db.Preload("Watchlist").Preload("Symbol").Offset(offset).Limit(limit).Find(&items).Error
	return items, err
}

// BatchCreate 批次建立觀察清單項目
func (r *watchlistItemRepository) BatchCreate(items []*models.WatchlistItem) error {
	return r.db.CreateInBatches(items, 100).Error
}

// GetSymbolsByWatchlistID 根據觀察清單 ID 取得所有股票
func (r *watchlistItemRepository) GetSymbolsByWatchlistID(watchlistID uint) ([]*models.Symbol, error) {
	var symbols []*models.Symbol
	err := r.db.Table("symbols").
		Joins("JOIN watchlist_items ON symbols.id = watchlist_items.symbol_id").
		Where("watchlist_items.watchlist_id = ?", watchlistID).
		Find(&symbols).Error
	return symbols, err
}

// IsSymbolInWatchlist 檢查股票是否在觀察清單中
func (r *watchlistItemRepository) IsSymbolInWatchlist(watchlistID, symbolID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.WatchlistItem{}).
		Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
