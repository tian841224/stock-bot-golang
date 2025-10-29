package repository

import (
	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type WatchlistRepository interface {
	Create(watchlist *models.Watchlist) error
	GetByID(id uint) (*models.Watchlist, error)
	GetByUserID(userID uint) (*models.Watchlist, error)
	GetByUserIDWithItems(userID uint) (*models.Watchlist, error)
	Update(watchlist *models.Watchlist) error
	Delete(id uint) error
	DeleteByUserID(userID uint) error
	List(offset, limit int) ([]*models.Watchlist, error)
	GetAllWithItems() ([]*models.Watchlist, error)
}

type watchlistRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository(db *gorm.DB) WatchlistRepository {
	return &watchlistRepository{db: db}
}

// Create 建立新觀察清單
func (r *watchlistRepository) Create(watchlist *models.Watchlist) error {
	return r.db.Create(watchlist).Error
}

// GetByID 根據 ID 取得觀察清單
func (r *watchlistRepository) GetByID(id uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.Preload("User").First(&watchlist, id).Error
	if err != nil {
		return nil, err
	}
	return &watchlist, nil
}

// GetByUserID 根據使用者 ID 取得觀察清單
func (r *watchlistRepository) GetByUserID(userID uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.Where("user_id = ?", userID).First(&watchlist).Error
	if err != nil {
		return nil, err
	}
	return &watchlist, nil
}

// GetByUserIDWithItems 根據使用者 ID 取得觀察清單（包含項目）
func (r *watchlistRepository) GetByUserIDWithItems(userID uint) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	err := r.db.Preload("WatchlistItems").Preload("WatchlistItems.Symbol").Where("user_id = ?", userID).First(&watchlist).Error
	if err != nil {
		return nil, err
	}
	return &watchlist, nil
}

// Update 更新觀察清單
func (r *watchlistRepository) Update(watchlist *models.Watchlist) error {
	return r.db.Save(watchlist).Error
}

// Delete 刪除觀察清單
func (r *watchlistRepository) Delete(id uint) error {
	return r.db.Delete(&models.Watchlist{}, id).Error
}

// DeleteByUserID 根據使用者 ID 刪除觀察清單
func (r *watchlistRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Watchlist{}).Error
}

// List 取得觀察清單列表
func (r *watchlistRepository) List(offset, limit int) ([]*models.Watchlist, error) {
	var watchlists []*models.Watchlist
	err := r.db.Preload("User").Offset(offset).Limit(limit).Find(&watchlists).Error
	return watchlists, err
}

// GetAllWithItems 取得所有觀察清單（包含項目）
func (r *watchlistRepository) GetAllWithItems() ([]*models.Watchlist, error) {
	var watchlists []*models.Watchlist
	err := r.db.Preload("User").Preload("WatchlistItems").Preload("WatchlistItems.Symbol").Find(&watchlists).Error
	return watchlists, err
}
