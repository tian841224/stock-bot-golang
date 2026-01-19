package health

import (
	"context"
	"runtime"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"

	"gorm.io/gorm"
)

type healthChecker struct {
	db               *gorm.DB
	finmindAPI       *finmindtrade.FinmindTradeAPI
	fugleAPI         *fugle.FugleAPI
	syncMetadataRepo port.SyncMetadataRepository
}

var _ port.HealthChecker = (*healthChecker)(nil)

func NewHealthChecker(
	db *gorm.DB,
	finmindAPI *finmindtrade.FinmindTradeAPI,
	fugleAPI *fugle.FugleAPI,
	syncMetadataRepo port.SyncMetadataRepository,
) port.HealthChecker {
	return &healthChecker{
		db:               db,
		finmindAPI:       finmindAPI,
		fugleAPI:         fugleAPI,
		syncMetadataRepo: syncMetadataRepo,
	}
}

func (c *healthChecker) CheckDatabase(ctx context.Context) port.HealthStatus {
	start := time.Now()

	var result int
	err := c.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error

	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		return port.HealthStatus{
			Status:       "unhealthy",
			Message:      err.Error(),
			ResponseTime: responseTime,
		}
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return port.HealthStatus{
			Status:       "degraded",
			Message:      "無法取得資料庫統計資訊",
			ResponseTime: responseTime,
		}
	}

	stats := sqlDB.Stats()
	if stats.OpenConnections >= stats.MaxOpenConnections-5 {
		return port.HealthStatus{
			Status:       "degraded",
			Message:      "連線池接近上限",
			ResponseTime: responseTime,
		}
	}

	return port.HealthStatus{
		Status:       "healthy",
		Message:      "Connected",
		ResponseTime: responseTime,
	}
}

func (c *healthChecker) CheckAPI(ctx context.Context, apiName string) port.HealthStatus {
	start := time.Now()

	switch apiName {
	case "finmind":
		if c.finmindAPI == nil {
			return port.HealthStatus{
				Status:       "degraded",
				Message:      "API client not initialized",
				ResponseTime: 0,
			}
		}

		_, err := c.finmindAPI.GetTodayInfo()
		responseTime := time.Since(start).Milliseconds()

		if err != nil {
			return port.HealthStatus{
				Status:       "degraded",
				Message:      err.Error(),
				ResponseTime: responseTime,
			}
		}

		return port.HealthStatus{
			Status:       "healthy",
			Message:      "API responsive",
			ResponseTime: responseTime,
		}

	case "fugle":
		if c.fugleAPI == nil {
			return port.HealthStatus{
				Status:       "degraded",
				Message:      "API client not initialized",
				ResponseTime: 0,
			}
		}

		responseTime := time.Since(start).Milliseconds()

		return port.HealthStatus{
			Status:       "healthy",
			Message:      "API client ready",
			ResponseTime: responseTime,
		}

	default:
		return port.HealthStatus{
			Status:       "unhealthy",
			Message:      "Unknown API",
			ResponseTime: 0,
		}
	}
}

func (c *healthChecker) CheckSyncStatus(ctx context.Context) port.SyncHealthStatus {
	if c.syncMetadataRepo == nil {
		return port.SyncHealthStatus{
			Status: "degraded",
			TaiwanStocks: port.SyncInfo{
				LastSuccess: nil,
				TotalCount:  0,
				LastError:   stringPtr("Sync metadata repository not initialized"),
			},
			USStocks: port.SyncInfo{
				LastSuccess: nil,
				TotalCount:  0,
				LastError:   stringPtr("Sync metadata repository not initialized"),
			},
		}
	}

	twMetadata, _ := c.syncMetadataRepo.GetByMarket(ctx, "TW")
	usMetadata, _ := c.syncMetadataRepo.GetByMarket(ctx, "US")

	status := "healthy"

	if twMetadata != nil && twMetadata.LastSuccessAt != nil {
		if time.Since(*twMetadata.LastSuccessAt) > 25*time.Hour {
			status = "degraded"
		}
	}

	if usMetadata != nil && usMetadata.LastSuccessAt != nil {
		if time.Since(*usMetadata.LastSuccessAt) > 25*time.Hour {
			status = "degraded"
		}
	}

	twInfo := port.SyncInfo{
		LastSuccess: nil,
		TotalCount:  0,
		LastError:   nil,
	}
	if twMetadata != nil {
		twInfo.LastSuccess = twMetadata.LastSuccessAt
		twInfo.TotalCount = twMetadata.TotalCount
		twInfo.LastError = twMetadata.LastError
	}

	usInfo := port.SyncInfo{
		LastSuccess: nil,
		TotalCount:  0,
		LastError:   nil,
	}
	if usMetadata != nil {
		usInfo.LastSuccess = usMetadata.LastSuccessAt
		usInfo.TotalCount = usMetadata.TotalCount
		usInfo.LastError = usMetadata.LastError
	}

	return port.SyncHealthStatus{
		Status:       status,
		TaiwanStocks: twInfo,
		USStocks:     usInfo,
	}
}

func (c *healthChecker) CheckResources(ctx context.Context) port.ResourceHealthStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memoryMB := int64(mem.Alloc / 1024 / 1024)
	goroutines := runtime.NumGoroutine()
	cpuCores := runtime.NumCPU()

	status := "healthy"
	if memoryMB > 1024 {
		status = "degraded"
	}
	if goroutines > 1000 {
		status = "degraded"
	}

	return port.ResourceHealthStatus{
		Status:        status,
		MemoryUsageMB: memoryMB,
		Goroutines:    goroutines,
		CPUCores:      cpuCores,
	}
}

func stringPtr(s string) *string {
	return &s
}
