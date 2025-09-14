package twse

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 註冊 TWSE 的路由
func RegisterRoutes(r *gin.Engine, handler *TwseHandler) {
	r.GET("/daily_market_info", func(c *gin.Context) {
		handler.GetDailyMarketInfo(nil, c)
	})
	r.GET("/after_trading_volume", handler.GetAfterTradingVolume)
	r.GET("/top_volume_items", handler.GetTopVolumeItems)
}
