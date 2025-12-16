package notification

import (
	"context"
	"fmt"
	"testing"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// mockUserAccountPort 用於測試的 UserAccountPort mock
type mockUserAccountPort struct {
	GetOrCreateFunc func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
}

func (m *mockUserAccountPort) GetOrCreate(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
	if m.GetOrCreateFunc != nil {
		return m.GetOrCreateFunc(ctx, accountID, userType)
	}
	return nil, nil
}

// mockValidationPort 用於測試的 ValidationPort mock
type mockValidationPort struct {
	ValidateSymbolFunc func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
}

func (m *mockValidationPort) ValidateSymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	if m.ValidateSymbolFunc != nil {
		return m.ValidateSymbolFunc(ctx, symbol)
	}
	return nil, nil
}

// mockUserSubscriptionPort 用於測試的 UserSubscriptionPort mock
type mockUserSubscriptionPort struct {
	GetUserSubscriptionListFunc      func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	AddUserSubscriptionStockFunc      func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	DeleteUserSubscriptionStockFunc   func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	GetUserSubscriptionStockListFunc  func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
	GetUserSubscriptionItemsFunc      func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	IsStockSubscribedFunc            func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	UpdateUserSubscriptionItemFunc    func(ctx context.Context, userID uint, item valueobject.SubscriptionType, status bool) (bool, error)
	AddUserSubscriptionItemFunc      func(ctx context.Context, userID uint, item valueobject.SubscriptionType) error
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	if m.GetUserSubscriptionListFunc != nil {
		return m.GetUserSubscriptionListFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.AddUserSubscriptionStockFunc != nil {
		return m.AddUserSubscriptionStockFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.DeleteUserSubscriptionStockFunc != nil {
		return m.DeleteUserSubscriptionStockFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
	if m.GetUserSubscriptionStockListFunc != nil {
		return m.GetUserSubscriptionStockListFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) GetUserSubscriptionItems(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	if m.GetUserSubscriptionItemsFunc != nil {
		return m.GetUserSubscriptionItemsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserSubscriptionPort) IsStockSubscribed(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
	if m.IsStockSubscribedFunc != nil {
		return m.IsStockSubscribedFunc(ctx, userID, stockSymbol)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) UpdateUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType, status bool) (bool, error) {
	if m.UpdateUserSubscriptionItemFunc != nil {
		return m.UpdateUserSubscriptionItemFunc(ctx, userID, item, status)
	}
	return false, nil
}

func (m *mockUserSubscriptionPort) AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error {
	if m.AddUserSubscriptionItemFunc != nil {
		return m.AddUserSubscriptionItemFunc(ctx, userID, item)
	}
	return nil
}

func TestSubscriptionUsecase_GetUserSubscriptionList(t *testing.T) {
	tests := []struct {
		name                 string
		userID               uint
		mockUserFunc         func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
		mockSubscriptionFunc func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
		expectError          bool
	}{
		{
			name:   "成功取得使用者訂閱列表",
			userID: 1,
			mockSubscriptionFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
				return []*dto.UserSubscriptionItem{
					{Item: valueobject.SubscriptionTypeStockInfo, Status: true},
					{Item: valueobject.SubscriptionTypeStockNews, Status: true},
					{Item: valueobject.SubscriptionTypeDailyMarketInfo, Status: true},
					{Item: valueobject.SubscriptionTypeTopVolumeItems, Status: true},
				}, nil
			},
			expectError: false,
		},
		{
			name:   "空訂閱列表",
			userID: 2,
			mockSubscriptionFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
				return []*dto.UserSubscriptionItem{}, nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSub := &mockUserSubscriptionPort{
				GetUserSubscriptionListFunc: tt.mockSubscriptionFunc,
			}

			subscriptions, err := mockSub.GetUserSubscriptionList(context.Background(), tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
				if subscriptions == nil {
					t.Errorf("期望有訂閱列表但為 nil")
				}
			}
		})
	}
}

func TestAddStockSubscription(t *testing.T) {
	tests := []struct {
		name                    string
		input                   SubscriptionInput
		mockUserFunc            func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
		mockValidateFunc        func(ctx context.Context, symbol string) (*entity.StockSymbol, error)
		mockIsSubscribedFunc    func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
		mockAddSubscriptionFunc func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
		expectSuccess           bool
		expectError             bool
		errorContains           string
	}{
		{
			name: "成功新增股票訂閱",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					ID:     1,
					Symbol: "2330",
					Name:   "台積電",
					Market: "TW",
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			mockAddSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "取得使用者失敗",
			input: SubscriptionInput{
				AccountID: "invalid_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return nil, fmt.Errorf("資料庫連線失敗")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "無法取得使用者",
		},
		{
			name: "使用者不存在",
			input: SubscriptionInput{
				AccountID: "nonexistent_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return nil, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "無此帳號",
		},
		{
			name: "驗證股票代號失敗",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "9999",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, fmt.Errorf("查詢失敗")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "驗證股票代號失敗",
		},
		{
			name: "股票代號無效",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "INVALID",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return nil, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "無此股票代號",
		},
		{
			name: "檢查訂閱狀態失敗",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					ID:     1,
					Symbol: "2330",
					Name:   "台積電",
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, fmt.Errorf("資料庫錯誤")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "檢查訂閱狀態失敗",
		},
		{
			name: "已訂閱此股票",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					ID:     1,
					Symbol: "2330",
					Name:   "台積電",
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "已訂閱此股票",
		},
		{
			name: "新增訂閱失敗_返回錯誤",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					ID:     1,
					Symbol: "2330",
					Name:   "台積電",
					Market: "TW",
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			mockAddSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, fmt.Errorf("資料庫錯誤")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "新增訂閱失敗",
		},
		{
			name: "新增訂閱失敗_返回false",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockValidateFunc: func(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
				return &entity.StockSymbol{
					ID:     1,
					Symbol: "2330",
					Name:   "台積電",
					Market: "TW",
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			mockAddSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "新增訂閱失敗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserAccount := &mockUserAccountPort{
				GetOrCreateFunc: tt.mockUserFunc,
			}
			mockValidation := &mockValidationPort{
				ValidateSymbolFunc: tt.mockValidateFunc,
			}
			mockSubscription := &mockUserSubscriptionPort{
				IsStockSubscribedFunc:        tt.mockIsSubscribedFunc,
				AddUserSubscriptionStockFunc: tt.mockAddSubscriptionFunc,
			}

			uc := NewSubscriptionsUsecase(mockUserAccount, mockSubscription, mockValidation)

			success, err := uc.AddStockSubscription(context.Background(), tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
			}

			if success != tt.expectSuccess {
				t.Errorf("返回值不符合期望，期望: %v, 實際: %v", tt.expectSuccess, success)
			}
		})
	}
}

func TestRemoveStockSubscription(t *testing.T) {
	tests := []struct {
		name                       string
		input                      SubscriptionInput
		mockUserFunc               func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
		mockIsSubscribedFunc       func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
		mockDeleteSubscriptionFunc func(ctx context.Context, userID uint, stockSymbol string) (bool, error)
		expectSuccess              bool
		expectError                bool
		errorContains              string
	}{
		{
			name: "成功移除股票訂閱",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			mockDeleteSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "取得使用者失敗",
			input: SubscriptionInput{
				AccountID: "invalid_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return nil, fmt.Errorf("資料庫連線失敗")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "無法取得使用者",
		},
		{
			name: "檢查訂閱狀態失敗",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, fmt.Errorf("資料庫錯誤")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "檢查訂閱狀態失敗",
		},
		{
			name: "尚未訂閱此股票",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "尚未訂閱此股票",
		},
		{
			name: "刪除訂閱失敗_返回錯誤",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			mockDeleteSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, fmt.Errorf("資料庫錯誤")
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "取消訂閱失敗",
		},
		{
			name: "刪除訂閱失敗_返回false",
			input: SubscriptionInput{
				AccountID: "test_user",
				UserType:  valueobject.UserTypeLine,
				Symbol:    "2330",
			},
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockIsSubscribedFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return true, nil
			},
			mockDeleteSubscriptionFunc: func(ctx context.Context, userID uint, stockSymbol string) (bool, error) {
				return false, nil
			},
			expectSuccess: false,
			expectError:   true,
			errorContains: "取消訂閱失敗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserAccount := &mockUserAccountPort{
				GetOrCreateFunc: tt.mockUserFunc,
			}
			mockSubscription := &mockUserSubscriptionPort{
				IsStockSubscribedFunc:           tt.mockIsSubscribedFunc,
				DeleteUserSubscriptionStockFunc: tt.mockDeleteSubscriptionFunc,
			}

			uc := NewSubscriptionsUsecase(mockUserAccount, mockSubscription, nil)

			success, err := uc.RemoveStockSubscription(context.Background(), tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
			}

			if success != tt.expectSuccess {
				t.Errorf("返回值不符合期望，期望: %v, 實際: %v", tt.expectSuccess, success)
			}
		})
	}
}

func TestListSubscriptions(t *testing.T) {
	tests := []struct {
		name                 string
		accountID            string
		userType             valueobject.UserType
		mockUserFunc         func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
		mockGetStockListFunc func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
		mockGetItemsFunc     func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
		expectError          bool
		errorContains        string
	}{
		{
			name:      "成功列出所有訂閱",
			accountID: "test_user",
			userType:  valueobject.UserTypeLine,
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockGetStockListFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
				return []*dto.UserSubscriptionStock{
					{Stock: "2330", Status: true},
					{Stock: "2317", Status: true},
				}, nil
			},
			mockGetItemsFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
				return []*dto.UserSubscriptionItem{
					{Item: valueobject.SubscriptionTypeStockInfo, Status: true},
					{Item: valueobject.SubscriptionTypeStockNews, Status: true},
				}, nil
			},
			expectError: false,
		},
		{
			name:      "取得使用者失敗",
			accountID: "invalid_user",
			userType:  valueobject.UserTypeLine,
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return nil, fmt.Errorf("資料庫連線失敗")
			},
			expectError:   true,
			errorContains: "無法取得使用者",
		},
		{
			name:      "取得訂閱股票列表失敗",
			accountID: "test_user",
			userType:  valueobject.UserTypeLine,
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockGetStockListFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
				return nil, fmt.Errorf("資料庫錯誤")
			},
			expectError:   true,
			errorContains: "無法取得訂閱股票列表",
		},
		{
			name:      "取得訂閱項目列表失敗",
			accountID: "test_user",
			userType:  valueobject.UserTypeLine,
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockGetStockListFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
				return []*dto.UserSubscriptionStock{
					{Stock: "2330", Status: true},
				}, nil
			},
			mockGetItemsFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
				return nil, fmt.Errorf("資料庫錯誤")
			},
			expectError:   true,
			errorContains: "無法取得訂閱設定",
		},
		{
			name:      "空訂閱列表",
			accountID: "test_user",
			userType:  valueobject.UserTypeLine,
			mockUserFunc: func(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
				return &entity.User{
					ID:        1,
					AccountID: "test_user",
					UserType:  valueobject.UserTypeLine,
				}, nil
			},
			mockGetStockListFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
				return []*dto.UserSubscriptionStock{}, nil
			},
			mockGetItemsFunc: func(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
				return []*dto.UserSubscriptionItem{}, nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserAccount := &mockUserAccountPort{
				GetOrCreateFunc: tt.mockUserFunc,
			}
			mockSubscription := &mockUserSubscriptionPort{
				GetUserSubscriptionStockListFunc: tt.mockGetStockListFunc,
				GetUserSubscriptionItemsFunc:     tt.mockGetItemsFunc,
			}

			uc := NewSubscriptionsUsecase(mockUserAccount, mockSubscription, nil)

			result, err := uc.ListSubscriptions(context.Background(), tt.accountID, tt.userType)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望錯誤但沒有發生錯誤")
				}
				if tt.errorContains != "" && err != nil {
					if !containsString(err.Error(), tt.errorContains) {
						t.Errorf("錯誤訊息不符合期望，期望包含: %s, 實際: %s", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望錯誤但發生錯誤: %v", err)
				}
				if result.UserID == "" {
					t.Errorf("期望有 UserID 但為空")
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
