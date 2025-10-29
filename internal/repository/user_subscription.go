package repository

import (
	"fmt"

	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

// UserSubscriptionRepository 使用者訂閱相關的 repository
type UserSubscriptionRepository interface {
	// 訂閱項目相關
	GetUserSubscriptionByItem(userID uint, item models.SubscriptionItem) (*models.Subscription, error)
	AddUserSubscriptionItem(userID uint, item models.SubscriptionItem) error
	UpdateUserSubscriptionItem(userID uint, item models.SubscriptionItem, status string) error
	GetUserSubscriptionList(userID uint) ([]*models.Subscription, error)

	// 訂閱股票相關
	AddUserSubscriptionStock(userID uint, stockSymbol string) (bool, error)
	DeleteUserSubscriptionStock(userID uint, stockSymbol string) (bool, error)
	GetUserSubscriptionStockList(userID uint) ([]*UserSubscriptionStock, error)
}

type userSubscriptionRepository struct {
	db *gorm.DB
}

// UserSubscriptionStock 使用者訂閱股票資訊
type UserSubscriptionStock struct {
	Stock  string `json:"stock"`
	Status int    `json:"status"`
}

func NewUserSubscriptionRepository(db *gorm.DB) UserSubscriptionRepository {
	return &userSubscriptionRepository{db: db}
}

// GetUserSubscriptionByItem 根據使用者 ID 和訂閱項目取得訂閱資料
func (r *userSubscriptionRepository) GetUserSubscriptionByItem(userID uint, item models.SubscriptionItem) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Joins("JOIN features ON features.id = subscriptions.feature_id").
		Where("subscriptions.user_id = ? AND features.code = ?", userID, fmt.Sprintf("%d", int(item))).
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// AddUserSubscriptionItem 新增使用者訂閱項目
func (r *userSubscriptionRepository) AddUserSubscriptionItem(userID uint, item models.SubscriptionItem) error {
	// 先檢查 feature 是否存在，不存在則建立
	var feature models.Feature
	itemCode := fmt.Sprintf("%d", int(item))
	err := r.db.Where("code = ?", itemCode).First(&feature).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			feature = models.Feature{
				Name:        item.GetName(),
				Code:        itemCode,
				Description: item.GetName(),
			}
			if err := r.db.Create(&feature).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 檢查是否已經訂閱
	var existingSubscription models.Subscription
	err = r.db.Where("user_id = ? AND feature_id = ?", userID, feature.ID).First(&existingSubscription).Error
	if err == nil {
		// 已存在，更新狀態為啟用
		existingSubscription.Status = "active"
		return r.db.Save(&existingSubscription).Error
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// 建立新的訂閱
	subscription := models.Subscription{
		UserID:    userID,
		FeatureID: feature.ID,
		Status:    "active",
	}
	return r.db.Create(&subscription).Error
}

// UpdateUserSubscriptionItem 更新使用者訂閱項目狀態
func (r *userSubscriptionRepository) UpdateUserSubscriptionItem(userID uint, item models.SubscriptionItem, status string) error {
	return r.db.Joins("JOIN features ON features.id = subscriptions.feature_id").
		Where("subscriptions.user_id = ? AND features.code = ?", userID, fmt.Sprintf("%d", int(item))).
		Update("status", status).Error
}

// GetUserSubscriptionList 取得使用者訂閱項目列表
func (r *userSubscriptionRepository) GetUserSubscriptionList(userID uint) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Preload("Feature").Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

// AddUserSubscriptionStock 新增使用者訂閱股票
func (r *userSubscriptionRepository) AddUserSubscriptionStock(userID uint, stockSymbol string) (bool, error) {
	// 先取得股票資訊
	var symbol models.Symbol
	err := r.db.Where("symbol = ? ", stockSymbol).First(&symbol).Error
	if err != nil {
		return false, err
	}

	// 取得股票功能
	var stockFeature models.Feature
	err = r.db.Where("code = ?", "1").First(&stockFeature).Error
	if err != nil {
		return false, err
	}

	// 檢查使用者是否已有股票訂閱
	var subscription models.Subscription
	err = r.db.Where("user_id = ? AND feature_id = ?", userID, stockFeature.ID).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			subscription = models.Subscription{
				UserID:    userID,
				FeatureID: stockFeature.ID,
				Status:    "active",
			}
			if err := r.db.Create(&subscription).Error; err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	// 檢查是否已經訂閱此股票
	var existingSubscriptionSymbol models.SubscriptionSymbol
	err = r.db.Where("subscription_id = ? AND symbol_id = ?", subscription.ID, symbol.ID).First(&existingSubscriptionSymbol).Error
	if err == nil {
		return false, nil // 已經訂閱過
	} else if err != gorm.ErrRecordNotFound {
		return false, err
	}

	// 建立新的股票訂閱
	subscriptionSymbol := models.SubscriptionSymbol{
		SubscriptionID: subscription.ID,
		SymbolID:       symbol.ID,
	}
	if err := r.db.Create(&subscriptionSymbol).Error; err != nil {
		return false, err
	}

	return true, nil
}

// DeleteUserSubscriptionStock 刪除使用者訂閱股票
func (r *userSubscriptionRepository) DeleteUserSubscriptionStock(userID uint, stockSymbol string) (bool, error) {
	// 先取得股票資訊
	var symbol models.Symbol
	err := r.db.Where("symbol = ?  ?", stockSymbol).First(&symbol).Error
	if err != nil {
		return false, err
	}

	// 取得股票功能
	var stockFeature models.Feature
	err = r.db.Where("code = ?", "股票資訊").First(&stockFeature).Error
	if err != nil {
		return false, err
	}

	// 取得使用者訂閱
	var subscription models.Subscription
	err = r.db.Where("user_id = ? AND feature_id = ?", userID, stockFeature.ID).First(&subscription).Error
	if err != nil {
		return false, err
	}

	// 刪除股票訂閱
	result := r.db.Where("subscription_id = ? AND symbol_id = ?", subscription.ID, symbol.ID).Delete(&models.SubscriptionSymbol{})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

// GetUserSubscriptionStockList 取得使用者訂閱股票列表
func (r *userSubscriptionRepository) GetUserSubscriptionStockList(userID uint) ([]*UserSubscriptionStock, error) {
	var results []*UserSubscriptionStock

	query := `
		SELECT s.symbol as stock, 1 as status
		FROM subscription_symbols ss
		JOIN subscriptions sub ON sub.id = ss.subscription_id
		JOIN symbols s ON s.id = ss.symbol_id
		JOIN features f ON f.id = sub.feature_id
		WHERE sub.user_id = ? AND f.code = '股票資訊' AND sub.status = 'active'
	`

	err := r.db.Raw(query, userID).Scan(&results).Error
	return results, err
}
