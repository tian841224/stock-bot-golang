package repository

import (
	"context"
	"fmt"

	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type WatchlistRepository interface {
	Create(ctx context.Context, watchlist *models.Watchlist) error
	GetByID(ctx context.Context, id uint) (*models.Watchlist, error)
	GetByUserID(ctx context.Context, userID uint) (*models.Watchlist, error)
	GetByUserIDWithItems(ctx context.Context, userID uint) (*models.Watchlist, error)
	Update(ctx context.Context, watchlist *models.Watchlist) error
	Delete(ctx context.Context, id uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Watchlist, error)
	GetAllWithItems(ctx context.Context) ([]*models.Watchlist, error)
}

type watchlistRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewWatchlistRepository(db *gorm.DB, log logger.Logger) WatchlistRepository {
	return &watchlistRepository{
		db:     db,
		logger: log,
	}
}

// Create 建立新觀察清單
func (r *watchlistRepository) Create(ctx context.Context, watchlist *models.Watchlist) error {
	r.logger.Info("Creating watchlist", logger.Any("user_id", watchlist.UserID))

	err := r.db.WithContext(ctx).Create(watchlist).Error
	if err != nil {
		r.logger.Error("Failed to create watchlist", logger.Error(err), logger.Any("user_id", watchlist.UserID))
		return err
	}

	r.logger.Info("Watchlist created successfully", logger.Any("id", watchlist.ID))
	return nil
}

// GetByID 根據 ID 取得觀察清單
func (r *watchlistRepository) GetByID(ctx context.Context, id uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.WithContext(ctx).First(&watchlist, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Watchlist not found", logger.Any("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get watchlist", logger.Error(err), logger.Any("id", id))
		return nil, err
	}
	return &watchlist, nil
}

// GetByUserID 根據使用者 ID 取得觀察清單
func (r *watchlistRepository) GetByUserID(ctx context.Context, userID uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&watchlist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &watchlist, nil
}

// GetByUserIDWithItems 根據使用者 ID 取得觀察清單（包含項目）
func (r *watchlistRepository) GetByUserIDWithItems(ctx context.Context, userID uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.WithContext(ctx).Preload("WatchlistItems").Preload("WatchlistItems.Symbol").Where("user_id = ?", userID).First(&watchlist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &watchlist, nil
}

// Update 更新觀察清單
func (r *watchlistRepository) Update(ctx context.Context, watchlist *models.Watchlist) error {
	r.logger.Info("Updating watchlist", logger.Any("id", watchlist.ID))

	err := r.db.WithContext(ctx).Save(watchlist).Error
	if err != nil {
		r.logger.Error("Failed to update watchlist", logger.Error(err), logger.Any("id", watchlist.ID))
		return err
	}

	r.logger.Info("Watchlist updated successfully", logger.Any("id", watchlist.ID))
	return nil
}

// Delete 刪除觀察清單
func (r *watchlistRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting watchlist", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.Watchlist{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete watchlist", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Watchlist not found for deletion", logger.Any("id", id))
		return fmt.Errorf("watchlist not found with id: %d", id)
	}

	r.logger.Info("Watchlist deleted successfully", logger.Any("id", id))
	return nil
}

// DeleteByUserID 根據使用者 ID 刪除觀察清單
func (r *watchlistRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	r.logger.Info("Deleting watchlist by user ID", logger.Any("user_id", userID))

	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Watchlist{})
	if result.Error != nil {
		r.logger.Error("Failed to delete watchlist by user ID", logger.Error(result.Error), logger.Any("user_id", userID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Watchlist not found for user", logger.Any("user_id", userID))
		return fmt.Errorf("watchlist not found for user id: %d", userID)
	}

	r.logger.Info("Watchlist deleted successfully", logger.Any("user_id", userID), logger.Int64("rows_affected", result.RowsAffected))
	return nil
}

// List 取得觀察清單列表
func (r *watchlistRepository) List(ctx context.Context, offset, limit int) ([]*models.Watchlist, error) {
	var watchlists []*models.Watchlist
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&watchlists).Error
	return watchlists, err
}

// GetAllWithItems 取得所有觀察清單（包含項目）
func (r *watchlistRepository) GetAllWithItems(ctx context.Context) ([]*models.Watchlist, error) {
	var watchlists []*models.Watchlist
	err := r.db.WithContext(ctx).Preload("User").Preload("WatchlistItems").Preload("WatchlistItems.Symbol").Find(&watchlists).Error
	return watchlists, err
}
