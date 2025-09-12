package models

import "time"

// 通知事件模型
type NotificationEvent struct {
	Model
	// 使用者ID
	UserID uint `gorm:"column:user_id;type:bigint;index" json:"user_id"`
	// 功能ID
	FeatureID uint `gorm:"column:feature_id;type:bigint;index" json:"feature_id"`
	// 訂閱ID
	SubscriptionID *uint `gorm:"column:subscription_id;type:bigint;index" json:"subscription_id"`
	// 股票ID
	SymbolID *uint `gorm:"column:symbol_id;type:bigint;index" json:"symbol_id"`
	// 事件內容
	Payload string `gorm:"column:payload;type:jsonb" json:"payload"`
	// 發生時間
	OccurredAt time.Time `gorm:"column:occurred_at;type:timestamptz" json:"occurred_at"`
	// 關聯資料表
	User         *User         `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Feature      *Feature      `gorm:"foreignKey:FeatureID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Symbol       *Symbols      `gorm:"foreignKey:SymbolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (NotificationEvent) TableName() string {
	return "notification_events"
}

func init() {
	RegisterModel(&NotificationEvent{})
}
