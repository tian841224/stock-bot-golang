package entity

import "time"

// SyncMetadata 同步元資料實體
type SyncMetadata struct {
	ID            uint
	Market        string
	LastSyncAt    *time.Time
	LastSuccessAt *time.Time
	LastError     *string
	TotalCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
