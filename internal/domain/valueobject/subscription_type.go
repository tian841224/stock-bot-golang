package valueobject

import "errors"

// SubscriptionType types
type SubscriptionType int

const (
	SubscriptionTypeDefault         SubscriptionType = 0
	SubscriptionTypeStockInfo       SubscriptionType = 1
	SubscriptionTypeStockNews       SubscriptionType = 2
	SubscriptionTypeDailyMarketInfo SubscriptionType = 3
	SubscriptionTypeTopVolumeItems  SubscriptionType = 4
)

// NewSubscriptionType 建立並驗證訂閱類型
func NewSubscriptionType(value int) (SubscriptionType, error) {
	st := SubscriptionType(value)
	if !st.IsValid() {
		return SubscriptionTypeDefault, errors.New("invalid subscription type")
	}
	return st, nil
}

// SubscriptionTypeMap mapping table for subscription types
var SubscriptionTypeMap = map[string]SubscriptionType{
	"0": SubscriptionTypeDefault,
	"1": SubscriptionTypeStockInfo,
	"2": SubscriptionTypeStockNews,
	"3": SubscriptionTypeDailyMarketInfo,
	"4": SubscriptionTypeTopVolumeItems,
}

// GetName returns the name of the subscription type
func (s SubscriptionType) GetName() string {
	switch s {
	case SubscriptionTypeDefault:
		return "Default"
	case SubscriptionTypeStockInfo:
		return "股票資訊"
	case SubscriptionTypeStockNews:
		return "股票新聞"
	case SubscriptionTypeDailyMarketInfo:
		return "每日大盤資訊"
	case SubscriptionTypeTopVolumeItems:
		return "交易量前20名"
	default:
		return "Default"
	}
}

// ParseSubscriptionType parses subscription type from input string
func ParseSubscriptionType(input string) (SubscriptionType, bool) {
	item, exists := SubscriptionTypeMap[input]
	return item, exists
}

// IsValid 驗證訂閱類型是否有效
func (s SubscriptionType) IsValid() bool {
	return s >= SubscriptionTypeDefault && s <= SubscriptionTypeTopVolumeItems
}

// Equals 比較兩個訂閱類型是否相等
func (s SubscriptionType) Equals(other SubscriptionType) bool {
	return s == other
}
