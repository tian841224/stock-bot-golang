package models

// 使用者觀察清單
type Watchlist struct {
	Model
	// 使用者ID
	UserID uint `gorm:"column:user_id;type:bigint" json:"user_id"`
	// 關聯資料表
	User           *User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	WatchlistItems []*WatchlistItem `gorm:"foreignKey:WatchlistID;references:ID" json:"-"`
}

func (Watchlist) TableName() string {
	return "watchlists"
}

func init() {
	RegisterModel(&Watchlist{})
}
