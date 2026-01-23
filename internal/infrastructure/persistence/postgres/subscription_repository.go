package repository

import (
	"context"
	"fmt"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.SubscriptionReader = (*subscriptionRepository)(nil)
var _ repo.SubscriptionWriter = (*subscriptionRepository)(nil)

func NewSubscriptionRepository(db *gorm.DB, log logger.Logger) *subscriptionRepository {
	return &subscriptionRepository{
		db:     db,
		logger: log,
	}
}

func (r *subscriptionRepository) toEntity(model *models.Subscription) *entity.Subscription {
	subscription := &entity.Subscription{
		ID:           model.ID,
		UserID:       model.UserID,
		FeatureID:    model.FeatureID,
		Item:         valueobject.SubscriptionType(model.FeatureID),
		Active:       model.Status,
		ScheduleCron: model.ScheduleCron,
	}

	// 轉換 Feature 關聯
	if model.Feature != nil {
		subscription.Feature = &entity.Feature{
			ID:          model.Feature.ID,
			Name:        model.Feature.Name,
			Code:        model.Feature.Code,
			Description: model.Feature.Description,
		}
	}

	// 轉換 User 關聯
	if model.User != nil {
		subscription.User = &entity.User{
			ID:        model.User.ID,
			AccountID: model.User.AccountID,
			UserType:  model.User.UserType,
			Status:    model.User.Status,
		}
	}

	return subscription
}

func (r *subscriptionRepository) toModel(entity *entity.Subscription) *models.Subscription {
	return &models.Subscription{
		UserID:       entity.UserID,
		FeatureID:    entity.FeatureID,
		Status:       entity.Active,
		ScheduleCron: entity.ScheduleCron,
	}
}

// GetByID 根據 ID 取得訂閱
func (r *subscriptionRepository) GetByID(ctx context.Context, id uint) (*entity.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").First(&subscription, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&subscription), nil
}

// GetByUserID 根據使用者 ID 取得訂閱
func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("Feature").Where("user_id = ?", userID).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// GetByFeatureID 根據功能 ID 取得訂閱
func (r *subscriptionRepository) GetByFeatureID(ctx context.Context, featureID uint) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Where("feature_id = ?", featureID).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// GetUserSubscriptionList 取得使用者訂閱項目列表
func (r *subscriptionRepository) GetUserSubscriptionList(ctx context.Context, userID uint) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("Feature").Where("user_id = ?", userID).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// GetByUserAndFeature 根據使用者和功能取得訂閱
func (r *subscriptionRepository) GetByUserAndFeature(ctx context.Context, userID, featureID uint) (*entity.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Where("user_id = ? AND feature_id = ?", userID, featureID).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&subscription), nil
}

// GetByStatus 根據狀態取得訂閱
func (r *subscriptionRepository) GetByStatus(ctx context.Context, status bool) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Where("status = ?", status).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// List 取得訂閱列表
func (r *subscriptionRepository) List(ctx context.Context, offset, limit int) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// GetActiveSubscriptions 取得啟用的訂閱
func (r *subscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Where("status = ?", true).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// GetBySchedule 根據排程取得訂閱
func (r *subscriptionRepository) GetBySchedule(ctx context.Context, scheduleCron string) ([]*entity.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Feature").Where("schedule_cron = ?", scheduleCron).Find(&subscriptions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.Subscription
	for _, subscription := range subscriptions {
		entities = append(entities, r.toEntity(subscription))
	}
	return entities, nil
}

// Create 建立新訂閱
func (r *subscriptionRepository) Create(ctx context.Context, subscription *entity.Subscription) error {
	r.logger.Info("Creating subscription", logger.Any("user_id", subscription.UserID), logger.Any("feature_id", subscription.FeatureID))

	dbModel := r.toModel(subscription)
	err := r.db.WithContext(ctx).Create(dbModel).Error
	if err != nil {
		r.logger.Error("Failed to create subscription", logger.Error(err), logger.Any("user_id", subscription.UserID))
		return err
	}

	r.logger.Info("Subscription created successfully", logger.Any("user_id", subscription.UserID))
	return nil
}

// Update 更新訂閱
func (r *subscriptionRepository) Update(ctx context.Context, subscription *entity.Subscription) error {
	r.logger.Info("Updating subscription", logger.Any("id", subscription.ID))

	dbModel := r.toModel(subscription)

	result := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("id = ?", subscription.ID).
		Updates(dbModel)

	if result.Error != nil {
		r.logger.Error("Failed to update subscription", logger.Error(result.Error), logger.Any("id", subscription.ID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Subscription not found for update", logger.Any("id", subscription.ID))
		return fmt.Errorf("subscription not found with id: %d", subscription.ID)
	}

	r.logger.Info("Subscription updated successfully", logger.Any("id", subscription.ID))
	return nil
}

// UpdateStatus 更新訂閱狀態
func (r *subscriptionRepository) UpdateStatus(ctx context.Context, id uint, status bool) error {
	r.logger.Info("Updating subscription status", logger.Any("id", id), logger.Bool("status", status))

	result := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		r.logger.Error("Failed to update subscription status", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Subscription not found for status update", logger.Any("id", id))
		return fmt.Errorf("subscription not found with id: %d", id)
	}

	r.logger.Info("Subscription status updated successfully", logger.Any("id", id), logger.Bool("status", status))
	return nil
}

// Delete 刪除訂閱
func (r *subscriptionRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting subscription", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.Subscription{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete subscription", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Subscription not found for deletion", logger.Any("id", id))
		return fmt.Errorf("subscription not found with id: %d", id)
	}

	r.logger.Info("Subscription deleted successfully", logger.Any("id", id))
	return nil
}
