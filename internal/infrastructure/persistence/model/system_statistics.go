package models

import "time"

// SystemStatistics 系統統計資料
type SystemStatistics struct {
	ID        uint      `gorm:"primarykey"`
	Key       string    `gorm:"column:key;type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     int64     `gorm:"column:value;type:bigint;not null;default:0" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemStatistics) TableName() string {
	return "system_statistics"
}

func init() {
	RegisterModel(&SystemStatistics{})
}
