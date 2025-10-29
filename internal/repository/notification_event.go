package repository

import (
	"time"

	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type NotificationEventRepository interface {
	Create(event *models.NotificationEvent) error
	GetByID(id uint) (*models.NotificationEvent, error)
	GetByUserID(userID uint) ([]*models.NotificationEvent, error)
	GetByFeatureID(featureID uint) ([]*models.NotificationEvent, error)
	GetBySubscriptionID(subscriptionID uint) ([]*models.NotificationEvent, error)
	GetBySymbolID(symbolID uint) ([]*models.NotificationEvent, error)
	Update(event *models.NotificationEvent) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.NotificationEvent, error)
	GetByDateRange(startDate, endDate time.Time) ([]*models.NotificationEvent, error)
	GetByUserAndFeature(userID, featureID uint) ([]*models.NotificationEvent, error)
	GetRecentEvents(userID uint, limit int) ([]*models.NotificationEvent, error)
	BatchCreate(events []*models.NotificationEvent) error
}

type notificationEventRepository struct {
	db *gorm.DB
}

func NewNotificationEventRepository(db *gorm.DB) NotificationEventRepository {
	return &notificationEventRepository{db: db}
}

// Create 建立新通知事件
func (r *notificationEventRepository) Create(event *models.NotificationEvent) error {
	return r.db.Create(event).Error
}

// GetByID 根據 ID 取得通知事件
func (r *notificationEventRepository) GetByID(id uint) (*models.NotificationEvent, error) {
	var event models.NotificationEvent
	err := r.db.Preload("User").Preload("Feature").Preload("Subscription").Preload("Symbol").First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByUserID 根據使用者 ID 取得通知事件
func (r *notificationEventRepository) GetByUserID(userID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("Feature").Preload("Symbol").Where("user_id = ?", userID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByFeatureID 根據功能 ID 取得通知事件
func (r *notificationEventRepository) GetByFeatureID(featureID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("User").Preload("Symbol").Where("feature_id = ?", featureID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetBySubscriptionID 根據訂閱 ID 取得通知事件
func (r *notificationEventRepository) GetBySubscriptionID(subscriptionID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("User").Preload("Feature").Preload("Symbol").Where("subscription_id = ?", subscriptionID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetBySymbolID 根據股票 ID 取得通知事件
func (r *notificationEventRepository) GetBySymbolID(symbolID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("User").Preload("Feature").Where("symbol_id = ?", symbolID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// Update 更新通知事件
func (r *notificationEventRepository) Update(event *models.NotificationEvent) error {
	return r.db.Save(event).Error
}

// Delete 刪除通知事件
func (r *notificationEventRepository) Delete(id uint) error {
	return r.db.Delete(&models.NotificationEvent{}, id).Error
}

// List 取得通知事件列表
func (r *notificationEventRepository) List(offset, limit int) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("User").Preload("Feature").Preload("Symbol").Offset(offset).Limit(limit).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByDateRange 根據日期範圍取得通知事件
func (r *notificationEventRepository) GetByDateRange(startDate, endDate time.Time) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Where("occurred_at BETWEEN ? AND ?", startDate, endDate).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByUserAndFeature 根據使用者和功能取得通知事件
func (r *notificationEventRepository) GetByUserAndFeature(userID, featureID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("Symbol").Where("user_id = ? AND feature_id = ?", userID, featureID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetRecentEvents 取得使用者最近的通知事件
func (r *notificationEventRepository) GetRecentEvents(userID uint, limit int) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.Preload("Feature").Preload("Symbol").Where("user_id = ?", userID).Order("occurred_at DESC").Limit(limit).Find(&events).Error
	return events, err
}

// BatchCreate 批次建立通知事件
func (r *notificationEventRepository) BatchCreate(events []*models.NotificationEvent) error {
	return r.db.CreateInBatches(events, 100).Error
}
