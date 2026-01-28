package health

import (
	"github.com/gin-gonic/gin"
	"github.com/tian841224/stock-bot/internal/application/usecase/health"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type HealthHandler struct {
	healthUsecase health.HealthCheckUsecase
	logger        logger.Logger
}

func NewHealthHandler(healthUsecase health.HealthCheckUsecase, log logger.Logger) *HealthHandler {
	return &HealthHandler{
		healthUsecase: healthUsecase,
		logger:        log,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	status, err := h.healthUsecase.GetHealthStatus(ctx)
	if err != nil {
		h.logger.Error("健康檢查失敗", logger.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	httpCode := 200
	if status.Status == "unhealthy" {
		httpCode = 503
	} else if status.Status == "degraded" {
		httpCode = 200
	}

	c.JSON(httpCode, status)
}
