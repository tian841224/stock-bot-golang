package models

// 使用者觀察清單
type WatchlistItem struct {
	Model
	// 使用者ID
	WatchlistID uint `gorm:"column:watchlist_id;type:bigint;index;uniqueIndex:idx_watchlist_symbol,priority:1" json:"watchlist_id"`
	// 股票代號
	SymbolID uint `gorm:"column:symbol_id;type:bigint;index;uniqueIndex:idx_watchlist_symbol,priority:2" json:"symbol_id"`
	// 關聯資料表
	Watchlist *Watchlist `gorm:"foreignKey:WatchlistID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Symbol    *Symbols   `gorm:"foreignKey:SymbolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (WatchlistItem) TableName() string {
	return "watchlist_items"
}

func init() {
	RegisterModel(&WatchlistItem{})
}
