package user

import (
	"context"

	"github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type UserSubscriptionUsecase interface {
	GetUserSubscriptionItemList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error)
	GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error)
	GetUserSubscriptionDetail(ctx context.Context, userID uint) (*dto.UserSubscriptionDetail, error)
	AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) (string, error)
	AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (string, error)
	DeleteUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) (string, error)
	DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (string, error)
}

type userSubscriptionUsecase struct {
	userAccountPort      port.UserAccountPort
	userSubscriptionPort port.UserSubscriptionPort
	validationPort       port.ValidationPort
}

var _ UserSubscriptionUsecase = (*userSubscriptionUsecase)(nil)

func NewUserSubscriptionUsecase(
	userAccountPort port.UserAccountPort,
	userSubscriptionPort port.UserSubscriptionPort,
	validationPort port.ValidationPort,
) UserSubscriptionUsecase {
	return &userSubscriptionUsecase{
		userAccountPort:      userAccountPort,
		userSubscriptionPort: userSubscriptionPort,
		validationPort:       validationPort,
	}
}

func (u *userSubscriptionUsecase) GetUserSubscriptionItemList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionItem, error) {
	return u.userSubscriptionPort.GetUserSubscriptionItemList(ctx, userID)
}

func (u *userSubscriptionUsecase) GetUserSubscriptionStockList(ctx context.Context, userID uint) ([]*dto.UserSubscriptionStock, error) {
	return u.userSubscriptionPort.GetUserSubscriptionStockList(ctx, userID)
}

func (u *userSubscriptionUsecase) GetUserSubscriptionDetail(ctx context.Context, userID uint) (*dto.UserSubscriptionDetail, error) {
	userSubscriptionItemList, err := u.userSubscriptionPort.GetUserSubscriptionItemList(ctx, userID)
	if err != nil {
		return nil, err
	}

	userSubscriptionStockList, err := u.userSubscriptionPort.GetUserSubscriptionStockList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserSubscriptionDetail{
		Items:  userSubscriptionItemList,
		Stocks: userSubscriptionStockList,
	}, nil
}
func (u *userSubscriptionUsecase) AddUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) (string, error) {
	userSubscriptionList, err := u.userSubscriptionPort.GetUserSubscriptionItemList(ctx, userID)
	if err != nil {
		return "", err
	}

	for _, userSubscription := range userSubscriptionList {
		if userSubscription.Item == item {
			return "已經訂閱過此項目:" + item.GetName(), nil
		}
	}

	return "訂閱成功", u.userSubscriptionPort.AddUserSubscriptionItem(ctx, userID, item)
}

func (u *userSubscriptionUsecase) AddUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (string, error) {
	stockSymbolEntity, err := u.validationPort.ValidateSymbol(ctx, stockSymbol)
	if err != nil {
		return "", err
	}

	userSubscriptionStockList, err := u.userSubscriptionPort.GetUserSubscriptionStockList(ctx, userID)
	if err != nil {
		return "", err
	}

	for _, userSubscriptionStock := range userSubscriptionStockList {
		if userSubscriptionStock.Symbol == stockSymbolEntity.Symbol {
			return "已經訂閱過此股票:" + stockSymbolEntity.Symbol, nil
		}
	}

	err = u.userSubscriptionPort.AddUserSubscriptionStock(ctx, userID, stockSymbolEntity.Symbol)
	if err != nil {
		return "", err
	}

	return "訂閱成功", nil
}

func (u *userSubscriptionUsecase) DeleteUserSubscriptionItem(ctx context.Context, userID uint, item valueobject.SubscriptionType) (string, error) {

	userSubscriptionItemList, err := u.userSubscriptionPort.GetUserSubscriptionItemList(ctx, userID)
	if err != nil {
		return "", err
	}

	found := false
	for _, userSubscriptionItem := range userSubscriptionItemList {
		if userSubscriptionItem.Item == item {
			found = true
			break
		}
	}

	if !found {
		return "未訂閱此項目:" + item.GetName(), nil
	}

	err = u.userSubscriptionPort.DeleteUserSubscriptionItem(ctx, userID, item)
	if err != nil {
		return "", err
	}

	return "已取消訂閱項目", nil
}

func (u *userSubscriptionUsecase) DeleteUserSubscriptionStock(ctx context.Context, userID uint, stockSymbol string) (string, error) {

	userSubscriptionStockList, err := u.userSubscriptionPort.GetUserSubscriptionStockList(ctx, userID)
	if err != nil {
		return "", err
	}

	found := false
	for _, userSubscriptionStock := range userSubscriptionStockList {
		if userSubscriptionStock.Symbol == stockSymbol {
			found = true
			break
		}
	}

	if !found {
		return "未訂閱此股票:" + stockSymbol, nil
	}

	err = u.userSubscriptionPort.DeleteUserSubscriptionStock(ctx, userID, stockSymbol)
	if err != nil {
		return "", err
	}

	return "已取消訂閱股票", nil
}
