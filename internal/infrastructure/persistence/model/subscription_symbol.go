package models

// 訂閱股票模型
type SubscriptionSymbol struct {
	Model
	// 訂閱ID
	SubscriptionID uint `gorm:"column:subscription_id;type:bigint;index;uniqueIndex:idx_sub_symbol,priority:1" json:"subscription_id"`
	// 股票ID
	SymbolID uint `gorm:"column:symbol_id;type:bigint;index;uniqueIndex:idx_sub_symbol,priority:2" json:"symbol_id"`
	// 關聯資料表
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StockSymbol  *StockSymbol  `gorm:"foreignKey:SymbolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (SubscriptionSymbol) TableName() string {
	return "subscription_symbols"
}

func init() {
	RegisterModel(&SubscriptionSymbol{})
}
