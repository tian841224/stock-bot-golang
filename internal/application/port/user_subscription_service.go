package port

import (
	"context"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type UserSubscriptionPort interface {

	// 更新使用者訂閱項目狀態
	// UpdateUserSubscriptionItem(userID uint, item valueobject.SubscriptionType, status bool) (bool, error)
	// 取得使用者訂閱項目列表
	GetUserSubscriptionList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	// 取得使用者訂閱股票列表
	GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
	// 取得使用者訂閱偏好設定
	GetUserSubscriptionItems(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	// 新增使用者訂閱項目
	AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error
	// 新增使用者訂閱股票
	AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	// 刪除使用者訂閱股票
	DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (bool, error)
	// 檢查使用者是否訂閱股票
	IsStockSubscribed(ctx context.Context, userID uint, stockSymbol string) (bool, error)
}
