package health

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type HealthCheckUsecase interface {
	GetHealthStatus(ctx context.Context) (*HealthCheckResponse, error)
}

type healthCheckUsecase struct {
	checker port.HealthChecker
	service string
	version string
	logger  logger.Logger
}

func NewHealthCheckUsecase(
	checker port.HealthChecker,
	service string,
	version string,
	log logger.Logger,
) HealthCheckUsecase {
	return &healthCheckUsecase{
		checker: checker,
		service: service,
		version: version,
		logger:  log,
	}
}

type HealthCheckResponse struct {
	Status         string                 `json:"status"`
	Timestamp      string                 `json:"timestamp"`
	Service        string                 `json:"service"`
	Version        string                 `json:"version"`
	Checks         map[string]interface{} `json:"checks"`
	OverallHealthy bool                   `json:"overall_healthy"`
}

func (h *healthCheckUsecase) GetHealthStatus(ctx context.Context) (*HealthCheckResponse, error) {
	checks := make(map[string]interface{})

	dbStatus := h.checker.CheckDatabase(ctx)
	checks["database"] = map[string]interface{}{
		"status":           dbStatus.Status,
		"message":          dbStatus.Message,
		"response_time_ms": dbStatus.ResponseTime,
	}

	finmindStatus := h.checker.CheckAPI(ctx, "finmind")
	checks["finmind_api"] = map[string]interface{}{
		"status":           finmindStatus.Status,
		"message":          finmindStatus.Message,
		"response_time_ms": finmindStatus.ResponseTime,
	}

	fugleStatus := h.checker.CheckAPI(ctx, "fugle")
	checks["fugle_api"] = map[string]interface{}{
		"status":           fugleStatus.Status,
		"message":          fugleStatus.Message,
		"response_time_ms": fugleStatus.ResponseTime,
	}

	syncStatus := h.checker.CheckSyncStatus(ctx)
	checks["last_sync"] = map[string]interface{}{
		"status": syncStatus.Status,
		"taiwan_stocks": map[string]interface{}{
			"last_success": syncStatus.TaiwanStocks.LastSuccess,
			"total_count":  syncStatus.TaiwanStocks.TotalCount,
			"last_error":   syncStatus.TaiwanStocks.LastError,
		},
		"us_stocks": map[string]interface{}{
			"last_success": syncStatus.USStocks.LastSuccess,
			"total_count":  syncStatus.USStocks.TotalCount,
			"last_error":   syncStatus.USStocks.LastError,
		},
	}

	resourceStatus := h.checker.CheckResources(ctx)
	checks["resources"] = map[string]interface{}{
		"status":          resourceStatus.Status,
		"memory_usage_mb": resourceStatus.MemoryUsageMB,
		"goroutines":      resourceStatus.Goroutines,
		"cpu_cores":       resourceStatus.CPUCores,
	}

	overallStatus, overallHealthy := h.determineOverallStatus(dbStatus, finmindStatus, fugleStatus, syncStatus, resourceStatus)

	response := &HealthCheckResponse{
		Status:         overallStatus,
		Timestamp:      time.Now().Format(time.RFC3339),
		Service:        h.service,
		Version:        h.version,
		Checks:         checks,
		OverallHealthy: overallHealthy,
	}

	return response, nil
}

func (h *healthCheckUsecase) determineOverallStatus(
	dbStatus port.HealthStatus,
	finmindStatus port.HealthStatus,
	fugleStatus port.HealthStatus,
	syncStatus port.SyncHealthStatus,
	resourceStatus port.ResourceHealthStatus,
) (string, bool) {
	if dbStatus.Status == "unhealthy" {
		return "unhealthy", false
	}

	degradedCount := 0
	if finmindStatus.Status == "degraded" || finmindStatus.Status == "unhealthy" {
		degradedCount++
	}
	if fugleStatus.Status == "degraded" || fugleStatus.Status == "unhealthy" {
		degradedCount++
	}
	if syncStatus.Status == "degraded" || syncStatus.Status == "unhealthy" {
		degradedCount++
	}
	if resourceStatus.Status == "degraded" {
		degradedCount++
	}

	if degradedCount > 0 {
		return "degraded", true
	}

	return "healthy", true
}
