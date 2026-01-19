package repository

import (
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type NotificationDeliveryRepository interface {
	Create(delivery *models.NotificationDelivery) error
	GetByID(id uint) (*models.NotificationDelivery, error)
	GetByEventID(eventID uint) ([]*models.NotificationDelivery, error)
	GetByChannelID(channelID uint) ([]*models.NotificationDelivery, error)
	GetByStatus(status models.NotificationDeliveryStatus) ([]*models.NotificationDelivery, error)
	Update(delivery *models.NotificationDelivery) error
	UpdateStatus(id uint, status models.NotificationDeliveryStatus) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.NotificationDelivery, error)
	GetByDateRange(startDate, endDate time.Time) ([]*models.NotificationDelivery, error)
	GetFailedDeliveries() ([]*models.NotificationDelivery, error)
	BatchCreate(deliveries []*models.NotificationDelivery) error
}

type notificationDeliveryRepository struct {
	db *gorm.DB
}

func NewNotificationDeliveryRepository(db *gorm.DB) NotificationDeliveryRepository {
	return &notificationDeliveryRepository{db: db}
}

// Create 建立新投遞紀錄
func (r *notificationDeliveryRepository) Create(delivery *models.NotificationDelivery) error {
	return r.db.Create(delivery).Error
}

// GetByID 根據 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByID(id uint) (*models.NotificationDelivery, error) {
	var delivery models.NotificationDelivery
	err := r.db.Preload("Event").First(&delivery, id).Error
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

// GetByEventID 根據事件 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByEventID(eventID uint) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Where("event_id = ?", eventID).Find(&deliveries).Error
	return deliveries, err
}

// GetByChannelID 根據管道 ID 取得投遞紀錄
func (r *notificationDeliveryRepository) GetByChannelID(channelID uint) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Where("channel_id = ?", channelID).Find(&deliveries).Error
	return deliveries, err
}

// GetByStatus 根據狀態取得投遞紀錄
func (r *notificationDeliveryRepository) GetByStatus(status models.NotificationDeliveryStatus) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Where("status = ?", status).Find(&deliveries).Error
	return deliveries, err
}

// Update 更新投遞紀錄
func (r *notificationDeliveryRepository) Update(delivery *models.NotificationDelivery) error {
	return r.db.Save(delivery).Error
}

// UpdateStatus 更新投遞狀態
func (r *notificationDeliveryRepository) UpdateStatus(id uint, status models.NotificationDeliveryStatus) error {
	return r.db.Model(&models.NotificationDelivery{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 刪除投遞紀錄
func (r *notificationDeliveryRepository) Delete(id uint) error {
	return r.db.Delete(&models.NotificationDelivery{}, id).Error
}

// List 取得投遞紀錄列表
func (r *notificationDeliveryRepository) List(offset, limit int) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Preload("Event").Offset(offset).Limit(limit).Order("created_at DESC").Find(&deliveries).Error
	return deliveries, err
}

// GetByDateRange 根據日期範圍取得投遞紀錄
func (r *notificationDeliveryRepository) GetByDateRange(startDate, endDate time.Time) ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Where("sent_at BETWEEN ? AND ?", startDate, endDate).Find(&deliveries).Error
	return deliveries, err
}

// GetFailedDeliveries 取得失敗的投遞紀錄
func (r *notificationDeliveryRepository) GetFailedDeliveries() ([]*models.NotificationDelivery, error) {
	var deliveries []*models.NotificationDelivery
	err := r.db.Where("status = ?", models.NotificationDeliveryStatusFailed).Find(&deliveries).Error
	return deliveries, err
}

// BatchCreate 批次建立投遞紀錄
func (r *notificationDeliveryRepository) BatchCreate(deliveries []*models.NotificationDelivery) error {
	return r.db.CreateInBatches(deliveries, 100).Error
}
