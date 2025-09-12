package models

// 訂閱模型
type Subscription struct {
	Model
	// 使用者ID
	UserID uint `gorm:"column:user_id;type:bigint;index:idx_subscriptions_user_feature,priority:1" json:"user_id"`
	// 功能ID
	FeatureID uint `gorm:"column:feature_id;type:bigint;index:idx_subscriptions_user_feature,priority:2" json:"feature_id"`
	// 狀態
	Status string `gorm:"column:status;type:varchar(255)" json:"status"`
	// 排程
	ScheduleCron string `gorm:"column:schedule_cron;type:varchar(255)" json:"schedule_cron"`
	// 關聯資料表
	User    *User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Feature *Feature `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

func init() {
	RegisterModel(&Subscription{})
}
