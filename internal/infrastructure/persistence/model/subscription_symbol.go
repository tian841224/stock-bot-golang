package models

// 訂閱股票模型
type SubscriptionSymbol struct {
	Model
	// 使用者ID
	UserID uint `gorm:"column:user_id;type:bigint;index" json:"user_id"`
	// 股票ID
	SymbolID uint `gorm:"column:symbol_id;type:bigint;index;uniqueIndex:idx_user_symbol,priority:1" json:"symbol_id"`
	// 關聯資料表
	StockSymbol *StockSymbol `gorm:"foreignKey:SymbolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User        *User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (SubscriptionSymbol) TableName() string {
	return "subscription_symbols"
}

func init() {
	RegisterModel(&SubscriptionSymbol{})
}
