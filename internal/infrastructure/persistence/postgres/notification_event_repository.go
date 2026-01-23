package repository

import (
	"context"
	"fmt"
	"time"

	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type NotificationEventRepository interface {
	Create(ctx context.Context, event *models.NotificationEvent) error
	GetByID(ctx context.Context, id uint) (*models.NotificationEvent, error)
	GetByUserID(ctx context.Context, userID uint) ([]*models.NotificationEvent, error)
	GetByFeatureID(ctx context.Context, featureID uint) ([]*models.NotificationEvent, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*models.NotificationEvent, error)
	GetBySymbolID(ctx context.Context, symbolID uint) ([]*models.NotificationEvent, error)
	Update(ctx context.Context, event *models.NotificationEvent) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.NotificationEvent, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NotificationEvent, error)
	GetByUserAndFeature(ctx context.Context, userID, featureID uint) ([]*models.NotificationEvent, error)
	GetRecentEvents(ctx context.Context, userID uint, limit int) ([]*models.NotificationEvent, error)
	BatchCreate(ctx context.Context, events []*models.NotificationEvent) error
}

type notificationEventRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewNotificationEventRepository(db *gorm.DB, log logger.Logger) NotificationEventRepository {
	return &notificationEventRepository{
		db:     db,
		logger: log,
	}
}

// Create 建立新通知事件
func (r *notificationEventRepository) Create(ctx context.Context, event *models.NotificationEvent) error {
	r.logger.Info("Creating notification event", logger.Any("user_id", event.UserID), logger.Any("feature_id", event.FeatureID))

	err := r.db.WithContext(ctx).Create(event).Error
	if err != nil {
		r.logger.Error("Failed to create notification event", logger.Error(err), logger.Any("user_id", event.UserID))
		return err
	}

	r.logger.Info("Notification event created successfully", logger.Any("id", event.ID))
	return nil
}

// GetByID 根據 ID 取得通知事件
func (r *notificationEventRepository) GetByID(ctx context.Context, id uint) (*models.NotificationEvent, error) {
	var event models.NotificationEvent
	// 優化: 移除過度的 Preload，4個關聯太重，改為按需加載
	err := r.db.WithContext(ctx).First(&event, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Notification event not found", logger.Any("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get notification event", logger.Error(err), logger.Any("id", id))
		return nil, err
	}
	return &event, nil
}

// GetByUserID 根據使用者 ID 取得通知事件
func (r *notificationEventRepository) GetByUserID(ctx context.Context, userID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	// 優化: 只保留顯示所需的關聯，移除過度加載
	err := r.db.WithContext(ctx).Preload("Feature").Where("user_id = ?", userID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByFeatureID 根據功能 ID 取得通知事件
func (r *notificationEventRepository) GetByFeatureID(ctx context.Context, featureID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Preload("User").Preload("Symbol").Where("feature_id = ?", featureID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetBySubscriptionID 根據訂閱 ID 取得通知事件
func (r *notificationEventRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Preload("Symbol").Where("subscription_id = ?", subscriptionID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetBySymbolID 根據股票 ID 取得通知事件
func (r *notificationEventRepository) GetBySymbolID(ctx context.Context, symbolID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Where("symbol_id = ?", symbolID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// Update 更新通知事件
func (r *notificationEventRepository) Update(ctx context.Context, event *models.NotificationEvent) error {
	r.logger.Info("Updating notification event", logger.Any("id", event.ID))

	err := r.db.WithContext(ctx).Save(event).Error
	if err != nil {
		r.logger.Error("Failed to update notification event", logger.Error(err), logger.Any("id", event.ID))
		return err
	}

	r.logger.Info("Notification event updated successfully", logger.Any("id", event.ID))
	return nil
}

// Delete 刪除通知事件
func (r *notificationEventRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting notification event", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.NotificationEvent{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete notification event", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Notification event not found for deletion", logger.Any("id", id))
		return fmt.Errorf("notification event not found with id: %d", id)
	}

	r.logger.Info("Notification event deleted successfully", logger.Any("id", id))
	return nil
}

// List 取得通知事件列表
func (r *notificationEventRepository) List(ctx context.Context, offset, limit int) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	// 優化: 列表查詢移除所有 Preload，大幅減少資料傳輸
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByDateRange 根據日期範圍取得通知事件
func (r *notificationEventRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Where("occurred_at BETWEEN ? AND ?", startDate, endDate).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetByUserAndFeature 根據使用者和功能取得通知事件
func (r *notificationEventRepository) GetByUserAndFeature(ctx context.Context, userID, featureID uint) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Preload("Symbol").Where("user_id = ? AND feature_id = ?", userID, featureID).Order("occurred_at DESC").Find(&events).Error
	return events, err
}

// GetRecentEvents 取得使用者最近的通知事件
func (r *notificationEventRepository) GetRecentEvents(ctx context.Context, userID uint, limit int) ([]*models.NotificationEvent, error) {
	var events []*models.NotificationEvent
	err := r.db.WithContext(ctx).Preload("Feature").Preload("Symbol").Where("user_id = ?", userID).Order("occurred_at DESC").Limit(limit).Find(&events).Error
	return events, err
}

// BatchCreate 批次建立通知事件
func (r *notificationEventRepository) BatchCreate(ctx context.Context, events []*models.NotificationEvent) error {
	r.logger.Info("Batch creating notification events", logger.Int("count", len(events)))

	err := r.db.WithContext(ctx).CreateInBatches(events, 100).Error
	if err != nil {
		r.logger.Error("Failed to batch create notification events", logger.Error(err), logger.Int("count", len(events)))
		return err
	}

	r.logger.Info("Batch create notification events completed", logger.Int("count", len(events)))
	return nil
}
