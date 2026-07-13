package models

// SubscriptionItem types
type SubscriptionItem int

const (
	SubscriptionItemDefault         SubscriptionItem = 0
	SubscriptionItemStockInfo       SubscriptionItem = 1
	SubscriptionItemStockNews       SubscriptionItem = 2
	SubscriptionItemDailyMarketInfo SubscriptionItem = 3
	SubscriptionItemTopVolumeItems  SubscriptionItem = 4
)

// GetName returns the name of the subscription item
func (s SubscriptionItem) GetName() string {
	switch s {
	case SubscriptionItemDefault:
		return "Default"
	case SubscriptionItemStockInfo:
		return "股票資訊"
	case SubscriptionItemStockNews:
		return "股票新聞"
	case SubscriptionItemDailyMarketInfo:
		return "每日大盤資訊"
	case SubscriptionItemTopVolumeItems:
		return "交易量前20名"
	default:
		return "Default"
	}
}
