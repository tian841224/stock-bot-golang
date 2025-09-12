package models

import "time"

// 投遞紀錄模型
type NotificationDelivery struct {
	Model
	// 事件ID
	EventID uint `gorm:"column:event_id;type:bigint;index" json:"event_id"`
	// 管道ID
	ChannelID uint `gorm:"column:channel_id;type:bigint;index" json:"channel_id"`
	// 投遞狀態
	Status NotificationDeliveryStatus `gorm:"column:status;type:varchar(255)" json:"status"`
	// 回應
	Response string `gorm:"column:response;type:jsonb" json:"response"`
	// 發送時間
	SentAt time.Time `gorm:"column:sent_at;type:timestamptz" json:"sent_at"`
	// 關聯資料表
	Event *NotificationEvent `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// 投遞狀態
type NotificationDeliveryStatus string

const (
	NotificationDeliveryStatusQueued NotificationDeliveryStatus = "queued"
	NotificationDeliveryStatusSent   NotificationDeliveryStatus = "sent"
	NotificationDeliveryStatusFailed NotificationDeliveryStatus = "failed"
)

func (NotificationDelivery) TableName() string {
	return "notification_deliveries"
}

func init() {
	RegisterModel(&NotificationDelivery{})
}
