package port

import (
	"context"
	"time"
)

// HealthChecker 定義健康檢查介面
type HealthChecker interface {
	CheckDatabase(ctx context.Context) HealthStatus
	CheckAPI(ctx context.Context, apiName string) HealthStatus
	CheckSyncStatus(ctx context.Context) SyncHealthStatus
	CheckResources(ctx context.Context) ResourceHealthStatus
}

// HealthStatus 健康狀態
type HealthStatus struct {
	Status       string
	Message      string
	ResponseTime int64
}

// SyncHealthStatus 同步健康狀態
type SyncHealthStatus struct {
	Status       string
	TaiwanStocks SyncInfo
	USStocks     SyncInfo
}

// SyncInfo 同步資訊
type SyncInfo struct {
	LastSuccess *time.Time
	TotalCount  int
	LastError   *string
}

// ResourceHealthStatus 資源健康狀態
type ResourceHealthStatus struct {
	Status        string
	MemoryUsageMB int64
	Goroutines    int
	CPUCores      int
}
