package models

// 股票代號模型
type Symbol struct {
	Model
	// 股票代號
	Symbol string `gorm:"column:symbol;type:varchar(255);not null;index:idx_symbol_market,priority:1" json:"symbol"`
	// 股票名稱
	Name string `gorm:"column:name;type:varchar(255)" json:"name"`
	// 市場
	Market string `gorm:"column:market;type:varchar(255);not null;index:idx_symbol_market,priority:2" json:"market"`
}

func (Symbol) TableName() string {
	return "symbols"
}

func init() {
	RegisterModel(&Symbol{})
}
