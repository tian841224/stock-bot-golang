// Package middleware 提供 Gin HTTP 中介層
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// RequestID 產生唯一的 request ID 並注入 context 與 response header
// 優先讀取 X-Request-Id header（Load Balancer/前端帶入），若無則自動產生
func RequestID(baseLogger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}

		// 建立帶有 request_id 的子 logger，後續 handler 可透過 context 取得
		reqLogger := baseLogger.With(logger.String("request_id", requestID))

		// 注入 context 讓下游函式使用
		ctx := logger.WithLogger(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// 回應 header 方便前端和 Load Balancer 對照
		c.Writer.Header().Set("X-Request-Id", requestID)

		c.Next()
	}
}
