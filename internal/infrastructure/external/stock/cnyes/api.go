// Package cnyes 提供鉅亨網 API 的實作
package cnyes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/cnyes/dto"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// CnyesAPI 鉅亨網 API 客戶端
type CnyesAPI struct {
	baseURL string
	client  *http.Client
	logger  logger.Logger
}

// NewCnyesAPI 建立新的鉅亨網 API 客戶端
func NewCnyesAPI(log logger.Logger) *CnyesAPI {
	return &CnyesAPI{
		baseURL: "https://ws.api.cnyes.com/ws/api/v1/quote/quotes",
		client:  &http.Client{Timeout: 10 * time.Second},
		logger:  log,
	}
}

// GetStockQuote 取得股票報價資訊
func (c *CnyesAPI) GetStockQuote(symbol string) (dto.CnyesStockQuoteResponseDto, error) {
	url := fmt.Sprintf("https://ws.api.cnyes.com/ws/api/v1/quote/quotes/TWS:%s:STOCK?column=K,E,KEY,M,AI", symbol)
	return getResponse[dto.CnyesStockQuoteResponseDto](c, url)
}

// GetRevenue 取得財報
func (c *CnyesAPI) GetRevenue(symbol string, months int) (response dto.CnyesRevenueResponseDto, err error) {
	url := fmt.Sprintf("https://marketinfo.api.cnyes.com/mi/api/v1/TWS:%s:STOCK/revenue?months=%d", symbol, months)
	return getResponse[dto.CnyesRevenueResponseDto](c, url)
}

func getResponse[T any](c *CnyesAPI, url string) (response T, err error) {
	c.logger.Debug("cnyes API request", logger.String("url", url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}

	// 發送請求
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("cnyes API request failed", logger.String("url", url), logger.Error(err))
		return response, fmt.Errorf("failed to connect to cnyes API: %v", err)
	}
	defer resp.Body.Close()

	// 檢查狀態碼
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("cnyes API returned non-200",
			logger.String("url", url),
			logger.Int("status_code", resp.StatusCode))
		return response, fmt.Errorf("cnyes API error, status: %d", resp.StatusCode)
	}

	// 解析 JSON 回應
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return response, nil
}
