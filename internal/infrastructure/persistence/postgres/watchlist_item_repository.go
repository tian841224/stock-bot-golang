package repository

import (
	"context"
	"fmt"

	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type WatchlistItemRepository interface {
	Create(ctx context.Context, item *models.WatchlistItem) error
	GetByID(ctx context.Context, id uint) (*models.WatchlistItem, error)
	GetByWatchlistID(ctx context.Context, watchlistID uint) ([]*models.WatchlistItem, error)
	GetBySymbolID(ctx context.Context, symbolID uint) ([]*models.WatchlistItem, error)
	GetByWatchlistAndSymbol(ctx context.Context, watchlistID, symbolID uint) (*models.WatchlistItem, error)
	Update(ctx context.Context, item *models.WatchlistItem) error
	Delete(ctx context.Context, id uint) error
	DeleteByWatchlistID(ctx context.Context, watchlistID uint) error
	DeleteByWatchlistAndSymbol(ctx context.Context, watchlistID, symbolID uint) error
	List(ctx context.Context, offset, limit int) ([]*models.WatchlistItem, error)
	BatchCreate(ctx context.Context, items []*models.WatchlistItem) error
	GetSymbolsByWatchlistID(ctx context.Context, watchlistID uint) ([]*models.StockSymbol, error)
	IsSymbolInWatchlist(ctx context.Context, watchlistID, symbolID uint) (bool, error)
}

type watchlistItemRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewWatchlistItemRepository(db *gorm.DB, log logger.Logger) WatchlistItemRepository {
	return &watchlistItemRepository{
		db:     db,
		logger: log,
	}
}

// Create 建立新觀察清單項目
func (r *watchlistItemRepository) Create(ctx context.Context, item *models.WatchlistItem) error {
	r.logger.Info("Creating watchlist item", logger.Any("watchlist_id", item.WatchlistID), logger.Any("symbol_id", item.SymbolID))

	err := r.db.WithContext(ctx).Create(item).Error
	if err != nil {
		r.logger.Error("Failed to create watchlist item", logger.Error(err), logger.Any("watchlist_id", item.WatchlistID))
		return err
	}

	r.logger.Info("Watchlist item created successfully", logger.Any("id", item.ID))
	return nil
}

// GetByID 根據 ID 取得觀察清單項目
func (r *watchlistItemRepository) GetByID(ctx context.Context, id uint) (*models.WatchlistItem, error) {
	var item models.WatchlistItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Watchlist item not found", logger.Any("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get watchlist item", logger.Error(err), logger.Any("id", id))
		return nil, err
	}
	return &item, nil
}

// GetByWatchlistID 根據觀察清單 ID 取得項目
func (r *watchlistItemRepository) GetByWatchlistID(ctx context.Context, watchlistID uint) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	err := r.db.WithContext(ctx).Preload("Symbol").Where("watchlist_id = ?", watchlistID).Find(&items).Error
	return items, err
}

// GetBySymbolID 根據股票 ID 取得觀察清單項目
func (r *watchlistItemRepository) GetBySymbolID(ctx context.Context, symbolID uint) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	err := r.db.WithContext(ctx).Preload("Watchlist").Where("symbol_id = ?", symbolID).Find(&items).Error
	return items, err
}

// GetByWatchlistAndSymbol 根據觀察清單和股票取得項目
func (r *watchlistItemRepository) GetByWatchlistAndSymbol(ctx context.Context, watchlistID, symbolID uint) (*models.WatchlistItem, error) {
	var item models.WatchlistItem
	err := r.db.WithContext(ctx).Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Update 更新觀察清單項目
func (r *watchlistItemRepository) Update(ctx context.Context, item *models.WatchlistItem) error {
	r.logger.Info("Updating watchlist item", logger.Any("id", item.ID))

	err := r.db.WithContext(ctx).Save(item).Error
	if err != nil {
		r.logger.Error("Failed to update watchlist item", logger.Error(err), logger.Any("id", item.ID))
		return err
	}

	r.logger.Info("Watchlist item updated successfully", logger.Any("id", item.ID))
	return nil
}

// Delete 刪除觀察清單項目
func (r *watchlistItemRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting watchlist item", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.WatchlistItem{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete watchlist item", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Watchlist item not found for deletion", logger.Any("id", id))
		return fmt.Errorf("watchlist item not found with id: %d", id)
	}

	r.logger.Info("Watchlist item deleted successfully", logger.Any("id", id))
	return nil
}

// DeleteByWatchlistID 根據觀察清單 ID 刪除所有項目
func (r *watchlistItemRepository) DeleteByWatchlistID(ctx context.Context, watchlistID uint) error {
	r.logger.Info("Deleting watchlist items by watchlist ID", logger.Any("watchlist_id", watchlistID))

	result := r.db.WithContext(ctx).Where("watchlist_id = ?", watchlistID).Delete(&models.WatchlistItem{})
	if result.Error != nil {
		r.logger.Error("Failed to delete watchlist items by watchlist ID", logger.Error(result.Error), logger.Any("watchlist_id", watchlistID))
		return result.Error
	}

	r.logger.Info("Watchlist items deleted successfully", logger.Any("watchlist_id", watchlistID), logger.Int64("rows_affected", result.RowsAffected))
	return nil
}

// DeleteByWatchlistAndSymbol 根據觀察清單和股票刪除項目
func (r *watchlistItemRepository) DeleteByWatchlistAndSymbol(ctx context.Context, watchlistID, symbolID uint) error {
	r.logger.Info("Deleting watchlist item by watchlist and symbol", logger.Any("watchlist_id", watchlistID), logger.Any("symbol_id", symbolID))

	result := r.db.WithContext(ctx).Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).Delete(&models.WatchlistItem{})
	if result.Error != nil {
		r.logger.Error("Failed to delete watchlist item", logger.Error(result.Error), logger.Any("watchlist_id", watchlistID), logger.Any("symbol_id", symbolID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Watchlist item not found for deletion", logger.Any("watchlist_id", watchlistID), logger.Any("symbol_id", symbolID))
		return fmt.Errorf("watchlist item not found with watchlist_id: %d and symbol_id: %d", watchlistID, symbolID)
	}

	r.logger.Info("Watchlist item deleted successfully", logger.Any("watchlist_id", watchlistID), logger.Any("symbol_id", symbolID))
	return nil
}

// List 取得觀察清單項目列表
func (r *watchlistItemRepository) List(ctx context.Context, offset, limit int) ([]*models.WatchlistItem, error) {
	var items []*models.WatchlistItem
	// 優化: 列表查詢移除 Preload
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&items).Error
	return items, err
}

// BatchCreate 批次建立觀察清單項目
func (r *watchlistItemRepository) BatchCreate(ctx context.Context, items []*models.WatchlistItem) error {
	r.logger.Info("Batch creating watchlist items", logger.Int("count", len(items)))

	err := r.db.WithContext(ctx).CreateInBatches(items, 100).Error
	if err != nil {
		r.logger.Error("Failed to batch create watchlist items", logger.Error(err), logger.Int("count", len(items)))
		return err
	}

	r.logger.Info("Batch create watchlist items completed", logger.Int("count", len(items)))
	return nil
}

// GetSymbolsByWatchlistID 根據觀察清單 ID 取得所有股票
func (r *watchlistItemRepository) GetSymbolsByWatchlistID(ctx context.Context, watchlistID uint) ([]*models.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Table("symbols").
		Joins("JOIN watchlist_items ON symbols.id = watchlist_items.symbol_id").
		Where("watchlist_items.watchlist_id = ?", watchlistID).
		Find(&symbols).Error
	return symbols, err
}

// IsSymbolInWatchlist 檢查股票是否在觀察清單中
func (r *watchlistItemRepository) IsSymbolInWatchlist(ctx context.Context, watchlistID, symbolID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WatchlistItem{}).
		Where("watchlist_id = ? AND symbol_id = ?", watchlistID, symbolID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
