package models

import "time"

// SyncMetadata 同步元資料模型
type SyncMetadata struct {
	Model
	Market        string     `gorm:"column:market;type:varchar(10);not null;uniqueIndex" json:"market"`
	LastSyncAt    *time.Time `gorm:"column:last_sync_at" json:"last_sync_at"`
	LastSuccessAt *time.Time `gorm:"column:last_success_at" json:"last_success_at"`
	LastError     *string    `gorm:"column:last_error;type:text" json:"last_error"`
	TotalCount    int        `gorm:"column:total_count;default:0" json:"total_count"`
}

func (SyncMetadata) TableName() string {
	return "sync_metadata"
}

func init() {
	RegisterModel(&SyncMetadata{})
}
