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

// SubscriptionItemMap mapping table for subscription items
var SubscriptionItemMap = map[string]SubscriptionItem{
	"0": SubscriptionItemDefault,
	"1": SubscriptionItemStockInfo,
	"2": SubscriptionItemStockNews,
	"3": SubscriptionItemDailyMarketInfo,
	"4": SubscriptionItemTopVolumeItems,
}

// GetName returns the name of the subscription item
func (s SubscriptionItem) GetName() string {
	switch s {
	case SubscriptionItemStockInfo:
		return "Stock Info"
	case SubscriptionItemStockNews:
		return "Stock News"
	case SubscriptionItemDailyMarketInfo:
		return "Daily Market Info"
	case SubscriptionItemTopVolumeItems:
		return "Top Volume Items"
	default:
		return "Default"
	}
}

// ParseSubscriptionItem parses subscription item from input string
func ParseSubscriptionItem(input string) (SubscriptionItem, bool) {
	item, exists := SubscriptionItemMap[input]
	return item, exists
}
