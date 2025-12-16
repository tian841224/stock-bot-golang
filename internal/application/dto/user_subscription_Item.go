package dto

import "github.com/tian841224/stock-bot/internal/domain/valueobject"

type UserSubscriptionItem struct {
	Item   valueobject.SubscriptionType
	Status bool
}
