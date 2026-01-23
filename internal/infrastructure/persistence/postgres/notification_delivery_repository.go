package repository

import (
	"context"
	"fmt"
	"time"

	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type NotificationDeliveryRepository interface {
	Create(ctx context.Context, delivery *models.NotificationDelivery) error
	GetByID(ctx context.Context, id uint) (*models.NotificationDelivery, error)
	GetByEventID(ctx context.Context, eventID uint) ([]*models.NotificationDelivery, error)
	GetByChannelID(ctx context.Context, channelID uint) ([]*models.NotificationDelivery, error)
	GetByStatus(ctx context.Context, status models.NotificationDeliveryStatus) ([]*models.NotificationDelivery, error)
	Update(ctx context.Context, delivery *models.NotificationDelivery) error
	UpdateStatus(ctx context.Context, id uint, status models.NotificationDeliveryStatus) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.NotificationDelivery, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NotificationDelivery, error)
	GetFailedDeliveries(ctx context.Context) ([]*models.NotificationDelivery, error)
	BatchCreate(ctx context.Context, deliveries []*models.NotificationDelivery) error
}

type notificationDeliveryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewNotificationDeliveryRepository(db *gorm.DB, log logger.Logger) NotificationDeliveryRepository {
	return &notificationDeliveryRepository{
		db:     db,
		logger: log,
	}
}

// Create 建立新投遞紀錄
func (r *notificationDeliveryRepository) Create(ctx context.Context, delivery *models.NotificationDelivery) error {
	r.logger.Info("Creating notification delivery", logger.Any("event_id", delivery.EventID), logger.Any("channel_id", delivery.ChannelID))

	err := r.db.WithContext(ctx).Create(delivery).Error
	if err != nil {
		r.logger.Error("Failed to create notification delivery", logger.Error(err), logger.Any("event_id", delivery.EventID))
		return err
	}

	r.logger.Info("Notification delivery created successfully", logger.Any("id", delivery.ID))
	return nil
}

// GetByID 根據 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByID(ctx context.Context, id uint) (*models.NotificationDelivery, error) {
	var delivery models.NotificationDelivery
	err := r.db.WithContext(ctx).Preload("Event").First(&delivery, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Notification delivery not found", logger.Any("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get notification delivery", logger.Error(err), logger.Any("id", id))
		return nil, err
	}
	return &delivery, nil
}

// GetByEventID 根據事件 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByEventID(ctx context.Context, eventID uint) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&deliveries).Error
	return deliveries, err
}

// GetByChannelID 根據管道 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByChannelID(ctx context.Context, channelID uint) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Where("channel_id = ?", channelID).Find(&deliveries).Error
	return deliveries, err
}

// GetByStatus 根據狀態取得投遞紀錄
func (r *notificationDeliveryRepository) GetByStatus(ctx context.Context, status models.NotificationDeliveryStatus) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&deliveries).Error
	return deliveries, err
}

// Update 更新投遞紀錄
func (r *notificationDeliveryRepository) Update(ctx context.Context, delivery *models.NotificationDelivery) error {
	r.logger.Info("Updating notification delivery", logger.Any("id", delivery.ID))

	err := r.db.WithContext(ctx).Save(delivery).Error
	if err != nil {
		r.logger.Error("Failed to update notification delivery", logger.Error(err), logger.Any("id", delivery.ID))
		return err
	}

	r.logger.Info("Notification delivery updated successfully", logger.Any("id", delivery.ID))
	return nil
}

// UpdateStatus 更新投遞狀態
func (r *notificationDeliveryRepository) UpdateStatus(ctx context.Context, id uint, status models.NotificationDeliveryStatus) error {
	r.logger.Info("Updating notification delivery status", logger.Any("id", id), logger.String("status", string(status)))

	err := r.db.WithContext(ctx).Model(&models.NotificationDelivery{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		r.logger.Error("Failed to update notification delivery status", logger.Error(err), logger.Any("id", id))
		return err
	}

	r.logger.Info("Notification delivery status updated successfully", logger.Any("id", id))
	return nil
}

// Delete 刪除投遞紀錄
func (r *notificationDeliveryRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting notification delivery", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.NotificationDelivery{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete notification delivery", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Notification delivery not found for deletion", logger.Any("id", id))
		return fmt.Errorf("notification delivery not found with id: %d", id)
	}

	r.logger.Info("Notification delivery deleted successfully", logger.Any("id", id))
	return nil
}

// List 取得投遞紀錄列表
func (r *notificationDeliveryRepository) List(ctx context.Context, offset, limit int) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&deliveries).Error
	return deliveries, err
}

// GetByDateRange 根據日期範圍取得投遞紀錄
func (r *notificationDeliveryRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Where("sent_at BETWEEN ? AND ?", startDate, endDate).Find(&deliveries).Error
	return deliveries, err
}

// GetFailedDeliveries 取得失敗的投遞紀錄
func (r *notificationDeliveryRepository) GetFailedDeliveries(ctx context.Context) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.WithContext(ctx).Where("status = ?", models.NotificationDeliveryStatusFailed).Find(&deliveries).Error
	return deliveries, err
}

// BatchCreate 批次建立投遞紀錄
func (r *notificationDeliveryRepository) BatchCreate(ctx context.Context, deliveries []*models.NotificationDelivery) error {
	r.logger.Info("Batch creating notification deliveries", logger.Int("count", len(deliveries)))

	err := r.db.WithContext(ctx).CreateInBatches(deliveries, 100).Error
	if err != nil {
		r.logger.Error("Failed to batch create notification deliveries", logger.Error(err), logger.Int("count", len(deliveries)))
		return err
	}

	r.logger.Info("Batch create notification deliveries completed", logger.Int("count", len(deliveries)))
	return nil
}
