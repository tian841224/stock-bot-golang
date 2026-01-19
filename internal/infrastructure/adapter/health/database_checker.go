package health

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"gorm.io/gorm"
)

type databaseHealthChecker struct {
	db *gorm.DB
}

func NewDatabaseHealthChecker(db *gorm.DB) *databaseHealthChecker {
	return &databaseHealthChecker{
		db: db,
	}
}

func (c *databaseHealthChecker) CheckDatabase(ctx context.Context) port.HealthStatus {
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
			Message:      "Cannot get database stats",
			ResponseTime: responseTime,
		}
	}

	stats := sqlDB.Stats()
	if stats.OpenConnections >= stats.MaxOpenConnections {
		return port.HealthStatus{
			Status:       "degraded",
			Message:      "Connection pool near limit",
			ResponseTime: responseTime,
		}
	}

	return port.HealthStatus{
		Status:       "healthy",
		Message:      "Connected",
		ResponseTime: responseTime,
	}
}
