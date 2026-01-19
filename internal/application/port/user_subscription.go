package port

import (
	"context"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type UserSubscriptionPort interface {
	// 取得使用者訂閱項目列表
	GetUserSubscriptionItemList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	// 取得使用者訂閱股票列表
	GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
	// 取得使用者訂閱詳細資料
	GetUserSubscriptionDetail(ctx context.Context, userID uint) (*dto.UserSubscriptionDetail, error)
	// 新增使用者訂閱項目
	AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error
	// 新增使用者訂閱股票
	AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) error
	// 刪除使用者訂閱股票
	DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) error
	// 刪除使用者訂閱項目
	DeleteUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) error
}
