package cnyes

import (
	twstockService "stock-bot/internal/service/twstock"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 註冊鉅亨網API的路由
func RegisterRoutes(r *gin.Engine, twstockService *twstockService.StockService) {
	// 建立處理器
	handler := NewCnyesHandler(twstockService)

	// 建立路由群組
	cnyesGroup := r.Group("/cnyes")
	{
		// 健康檢查
		cnyesGroup.GET("/health", handler.GetHealthCheck)

		// 股票相關路由
		stockGroup := cnyesGroup.Group("/stock")
		{
			// 取得格式化的股票報價資訊
			stockGroup.GET("/:stock_id", handler.GetStockQuote)

			// 取得原始股票報價資料
			stockGroup.GET("/:stock_id/raw", handler.GetStockQuoteRaw)

			// 取得股票營收資料
			stockGroup.GET("/:stock_id/revenue", handler.GetStockRevenue)

			// 取得股票營收圖表
			stockGroup.GET("/:stock_id/revenue/chart", handler.GetStockRevenueChart)
		}
	}
}
