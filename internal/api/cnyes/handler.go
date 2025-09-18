package cnyes

import (
	"net/http"
	cnyesService "stock-bot/internal/service/cnyes"

	"github.com/gin-gonic/gin"
)

// CnyesHandler 鉅亨網API處理器
type CnyesHandler struct {
	cnyesService cnyesService.CnyesServiceInterface
}

// NewCnyesHandler 建立新的鉅亨網處理器
func NewCnyesHandler(cnyesService cnyesService.CnyesServiceInterface) *CnyesHandler {
	return &CnyesHandler{
		cnyesService: cnyesService,
	}
}

// GetStockQuote 取得股票報價資訊
// @Summary 取得股票即時報價
// @Description 透過鉅亨網API取得指定股票的即時報價資訊
// @Tags 鉅亨網API
// @Accept json
// @Produce json
// @Param stock_id path string true "股票代碼 (例如: 2330)"
// @Success 200 {object} cnyesService.StockQuoteInfo "成功取得股票報價資訊"
// @Failure 400 {object} map[string]interface{} "請求參數錯誤"
// @Failure 404 {object} map[string]interface{} "查無股票資料"
// @Failure 500 {object} map[string]interface{} "內部伺服器錯誤"
// @Router /cnyes/stock/{stock_id} [get]
func (h *CnyesHandler) GetStockQuote(c *gin.Context) {
	stockID := c.Param("stock_id")

	// 檢查必要參數
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "股票代碼為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	// 呼叫服務取得股票資訊
	stockInfo, err := h.cnyesService.GetStockQuote(stockID)
	if err != nil {
		// 根據錯誤類型回傳不同的HTTP狀態碼
		if err.Error() == "查無股票資料: "+stockID {
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "查無股票資料",
				"code":     "STOCK_NOT_FOUND",
				"stock_id": stockID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得股票資訊失敗",
			"code":   "INTERNAL_ERROR",
			"detail": err.Error(),
		})
		return
	}

	// 回傳成功結果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stockInfo,
		"message": "成功取得股票報價資訊",
	})
}

// GetStockQuoteRaw 取得原始股票報價資料
// @Summary 取得原始股票報價資料
// @Description 直接回傳鉅亨網API的原始回應資料
// @Tags 鉅亨網API
// @Accept json
// @Produce json
// @Param stock_id path string true "股票代碼 (例如: 2330)"
// @Success 200 {object} dto.CnyesStockQuoteResponseDto "原始API回應資料"
// @Failure 400 {object} map[string]interface{} "請求參數錯誤"
// @Failure 500 {object} map[string]interface{} "內部伺服器錯誤"
// @Router /cnyes/stock/{stock_id}/raw [get]
func (h *CnyesHandler) GetStockQuoteRaw(c *gin.Context) {
	stockID := c.Param("stock_id")

	// 檢查必要參數
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "股票代碼為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	// 建構股票符號
	symbol := "TWS:" + stockID + ":STOCK"

	// 透過服務取得格式化資料
	stockInfo, err := h.cnyesService.GetStockQuote(stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得股票資訊失敗",
			"code":   "INTERNAL_ERROR",
			"detail": err.Error(),
		})
		return
	}

	// 回傳格式化後的資料
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stockInfo,
		"message": "取得股票報價資料",
		"symbol":  symbol,
	})
}

// GetHealthCheck 健康檢查
// @Summary 健康檢查
// @Description 檢查鉅亨網API服務是否正常運作
// @Tags 鉅亨網API
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服務正常"
// @Router /cnyes/health [get]
func (h *CnyesHandler) GetHealthCheck(c *gin.Context) {
	// 嘗試取得台積電的資料來測試API連線
	_, err := h.cnyesService.GetStockQuote("2330")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "cnyes-api",
			"error":   err.Error(),
			"message": "鉅亨網API服務異常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "cnyes-api",
		"message": "鉅亨網API服務正常",
	})
}
