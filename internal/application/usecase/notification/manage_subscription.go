package notification

import (
	"context"
	"fmt"

	port "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// ========== Input/Output 結構定義 ==========

type SubscriptionInput struct {
	AccountID string
	UserType  valueobject.UserType
	Symbol    string
}

type UpdatePreferenceInput struct {
	AccountID string
	UserType  valueobject.UserType
	Item      string
	Status    bool
}

type ListSubscriptionsOutput struct {
	UserID            string
	Stocks            []*dto.UserSubscriptionStock
	SubscriptionTypes []*dto.UserSubscriptionItem
}

// ========== Usecase 介面定義 ==========

type SubscriptionsUsecase interface {
	AddStockSubscription(ctx context.Context, input SubscriptionInput) (bool, error)
	RemoveStockSubscription(ctx context.Context, input SubscriptionInput) (bool, error)
	ListSubscriptions(ctx context.Context, accountID string, userType valueobject.UserType) (ListSubscriptionsOutput, error)
}

// ========== Usecase 實作 ==========

type subscriptionsUsecase struct {
	userAccountPort      port.UserAccountPort
	userSubscriptionPort port.UserSubscriptionPort
	validationPort       port.ValidationPort
}

func NewSubscriptionsUsecase(
	userAccountPort port.UserAccountPort,
	userSubscriptionPort port.UserSubscriptionPort,
	validationPort port.ValidationPort,
) SubscriptionsUsecase {
	return &subscriptionsUsecase{
		userAccountPort:      userAccountPort,
		userSubscriptionPort: userSubscriptionPort,
		validationPort:       validationPort,
	}
}

// AddStockSubscription 新增股票訂閱
func (uc *subscriptionsUsecase) AddStockSubscription(ctx context.Context, input SubscriptionInput) (bool, error) {
	// 1. 取得或建立使用者
	user, err := uc.getUserOrFail(ctx, input.AccountID, input.UserType)
	if err != nil {
		return false, err
	}

	// 2. 驗證股票代號
	if err := uc.validateStock(ctx, input.Symbol); err != nil {
		return false, err
	}

	// 3. 檢查是否已訂閱
	exists, err := uc.userSubscriptionPort.IsStockSubscribed(ctx, user.ID, input.Symbol)
	if err != nil {
		return false, fmt.Errorf("檢查訂閱狀態失敗: %w", err)
	}
	if exists {
		return false, fmt.Errorf("已訂閱此股票: %s", input.Symbol)
	}

	// 4. 新增訂閱
	success, err := uc.userSubscriptionPort.AddUserSubscriptionStock(ctx, user.ID, input.Symbol)
	if err != nil || !success {
		return false, fmt.Errorf("新增訂閱失敗: %w", err)
	}

	return true, nil
}

// RemoveStockSubscription 移除股票訂閱
func (uc *subscriptionsUsecase) RemoveStockSubscription(ctx context.Context, input SubscriptionInput) (bool, error) {
	// 1. 取得使用者
	user, err := uc.getUserOrFail(ctx, input.AccountID, input.UserType)
	if err != nil {
		return false, err
	}

	// 2. 檢查訂閱是否存在
	exists, err := uc.userSubscriptionPort.IsStockSubscribed(ctx, user.ID, input.Symbol)
	if err != nil {
		return false, fmt.Errorf("檢查訂閱狀態失敗: %w", err)
	}
	if !exists {
		return false, fmt.Errorf("尚未訂閱此股票: %s", input.Symbol)
	}

	// 3. 刪除訂閱
	success, err := uc.userSubscriptionPort.DeleteUserSubscriptionStock(ctx, user.ID, input.Symbol)
	if err != nil || !success {
		return false, fmt.Errorf("取消訂閱失敗: %w", err)
	}

	return true, nil
}

// ListSubscriptions 列出使用者的所有訂閱
func (uc *subscriptionsUsecase) ListSubscriptions(ctx context.Context, accountID string, userType valueobject.UserType) (ListSubscriptionsOutput, error) {
	// 1. 取得使用者
	user, err := uc.getUserOrFail(ctx, accountID, userType)
	if err != nil {
		return ListSubscriptionsOutput{}, err
	}

	// 2. 取得訂閱股票列表
	stocks, err := uc.userSubscriptionPort.GetUserSubscriptionStockList(ctx, user.ID)
	if err != nil {
		return ListSubscriptionsOutput{}, fmt.Errorf("無法取得訂閱股票列表: %w", err)
	}

	// 3. 取得訂閱項目列表
	items, err := uc.userSubscriptionPort.GetUserSubscriptionItems(ctx, user.ID)
	if err != nil {
		return ListSubscriptionsOutput{}, fmt.Errorf("無法取得訂閱設定: %w", err)
	}

	return ListSubscriptionsOutput{
		UserID:            user.AccountID,
		Stocks:            stocks,
		SubscriptionTypes: items,
	}, nil
}

// ========== 私有輔助方法 ==========

// getUserOrFail 取得使用者，失敗則返回錯誤
func (uc *subscriptionsUsecase) getUserOrFail(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
	user, err := uc.userAccountPort.GetOrCreate(ctx, accountID, userType)
	if err != nil {
		return nil, fmt.Errorf("無法取得使用者: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("無此帳號: %s", accountID)
	}
	return user, nil
}

// validateStock 驗證股票代號是否有效
func (uc *subscriptionsUsecase) validateStock(ctx context.Context, symbol string) error {
	validated, err := uc.validationPort.ValidateSymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("驗證股票代號失敗: %w", err)
	}
	if validated == nil {
		return fmt.Errorf("無此股票代號: %s", symbol)
	}
	return nil
}
