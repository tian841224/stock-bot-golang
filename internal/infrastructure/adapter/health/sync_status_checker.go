package health

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
)

type syncStatusChecker struct {
	syncMetadataRepo port.SyncMetadataRepository
}

func NewSyncStatusChecker(syncMetadataRepo port.SyncMetadataRepository) *syncStatusChecker {
	return &syncStatusChecker{
		syncMetadataRepo: syncMetadataRepo,
	}
}

func (c *syncStatusChecker) CheckSyncStatus(ctx context.Context) port.SyncHealthStatus {
	status := "healthy"

	twMetadata, _ := c.syncMetadataRepo.GetByMarket(ctx, "TW")
	usMetadata, _ := c.syncMetadataRepo.GetByMarket(ctx, "US")

	twInfo := port.SyncInfo{}
	if twMetadata != nil {
		twInfo.LastSuccess = twMetadata.LastSuccessAt
		twInfo.TotalCount = twMetadata.TotalCount
		twInfo.LastError = twMetadata.LastError

		if twMetadata.LastSuccessAt != nil {
			hoursSinceSync := time.Since(*twMetadata.LastSuccessAt).Hours()
			if hoursSinceSync > 25 {
				status = "degraded"
			}
		}

		if twMetadata.LastError != nil && *twMetadata.LastError != "" {
			status = "degraded"
		}
	}

	usInfo := port.SyncInfo{}
	if usMetadata != nil {
		usInfo.LastSuccess = usMetadata.LastSuccessAt
		usInfo.TotalCount = usMetadata.TotalCount
		usInfo.LastError = usMetadata.LastError

		if usMetadata.LastSuccessAt != nil {
			hoursSinceSync := time.Since(*usMetadata.LastSuccessAt).Hours()
			if hoursSinceSync > 25 {
				status = "degraded"
			}
		}

		if usMetadata.LastError != nil && *usMetadata.LastError != "" {
			status = "degraded"
		}
	}

	return port.SyncHealthStatus{
		Status:       status,
		TaiwanStocks: twInfo,
		USStocks:     usInfo,
	}
}
