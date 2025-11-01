package repository

import (
	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uint) (*models.Subscription, error)
	GetByUserID(userID uint) ([]*models.Subscription, error)
	GetByFeatureID(featureID uint) ([]*models.Subscription, error)
	GetByUserAndFeature(userID, featureID uint) (*models.Subscription, error)
	GetByStatus(status bool) ([]*models.Subscription, error)
	Update(subscription *models.Subscription) error
	UpdateStatus(id uint, status bool) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.Subscription, error)
	GetActiveSubscriptions() ([]*models.Subscription, error)
	GetBySchedule(scheduleCron string) ([]*models.Subscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create 建立新訂閱
func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

// GetByID 根據 ID 取得訂閱
func (r *subscriptionRepository) GetByID(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("User").Preload("Feature").First(&subscription, id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetByUserID 根據使用者 ID 取得訂閱
func (r *subscriptionRepository) GetByUserID(userID uint) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("Feature").Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

// GetByFeatureID 根據功能 ID 取得訂閱
func (r *subscriptionRepository) GetByFeatureID(featureID uint) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("User").Where("feature_id = ?", featureID).Find(&subscriptions).Error
	return subscriptions, err
}

// GetByUserAndFeature 根據使用者和功能取得訂閱
func (r *subscriptionRepository) GetByUserAndFeature(userID, featureID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Where("user_id = ? AND feature_id = ?", userID, featureID).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetByStatus 根據狀態取得訂閱
func (r *subscriptionRepository) GetByStatus(status bool) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("User").Preload("Feature").Where("status = ?", status).Find(&subscriptions).Error
	return subscriptions, err
}

// Update 更新訂閱
func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

// UpdateStatus 更新訂閱狀態
func (r *subscriptionRepository) UpdateStatus(id uint, status bool) error {
	return r.db.Model(&models.Subscription{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 刪除訂閱
func (r *subscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}

// List 取得訂閱列表
func (r *subscriptionRepository) List(offset, limit int) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("User").Preload("Feature").Offset(offset).Limit(limit).Find(&subscriptions).Error
	return subscriptions, err
}

// GetActiveSubscriptions 取得啟用的訂閱
func (r *subscriptionRepository) GetActiveSubscriptions() ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("User").Preload("Feature").Where("status = ?", true).Find(&subscriptions).Error
	return subscriptions, err
}

// GetBySchedule 根據排程取得訂閱
func (r *subscriptionRepository) GetBySchedule(scheduleCron string) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("User").Preload("Feature").Where("schedule_cron = ?", scheduleCron).Find(&subscriptions).Error
	return subscriptions, err
}
