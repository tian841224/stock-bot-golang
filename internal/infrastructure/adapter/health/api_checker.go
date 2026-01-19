package health

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle"
)

type apiHealthChecker struct {
	finmindAPI *finmindtrade.FinmindTradeAPI
	fugleAPI   *fugle.FugleAPI
}

func NewAPIHealthChecker(finmindAPI *finmindtrade.FinmindTradeAPI, fugleAPI *fugle.FugleAPI) *apiHealthChecker {
	return &apiHealthChecker{
		finmindAPI: finmindAPI,
		fugleAPI:   fugleAPI,
	}
}

func (c *apiHealthChecker) CheckAPI(ctx context.Context, apiName string) port.HealthStatus {
	start := time.Now()

	switch apiName {
	case "finmind":
		if c.finmindAPI == nil {
			return port.HealthStatus{
				Status:       "degraded",
				Message:      "FinMind API not initialized",
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
				Message:      "Fugle API not initialized",
				ResponseTime: 0,
			}
		}

		responseTime := time.Since(start).Milliseconds()

		return port.HealthStatus{
			Status:       "healthy",
			Message:      "API available",
			ResponseTime: responseTime,
		}

	default:
		return port.HealthStatus{
			Status:       "unhealthy",
			Message:      "Unknown API: " + apiName,
			ResponseTime: 0,
		}
	}
}
