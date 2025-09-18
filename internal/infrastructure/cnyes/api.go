package cnyes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-bot/internal/infrastructure/cnyes/dto"
	"time"
)

// CnyesAPIInterface 定義鉅亨網 API 的介面
type CnyesAPIInterface interface {
	GetStockQuote(symbol string) (dto.CnyesStockQuoteResponseDto, error)
}

// CnyesAPI 鉅亨網 API 客戶端
type CnyesAPI struct {
	baseURL string
	client  *http.Client
}

// NewCnyesAPI 建立新的鉅亨網 API 客戶端
func NewCnyesAPI() *CnyesAPI {
	return &CnyesAPI{
		baseURL: "https://ws.api.cnyes.com/ws/api/v1/quote/quotes",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetStockQuote 取得股票報價資訊
func (c *CnyesAPI) GetStockQuote(symbol string) (dto.CnyesStockQuoteResponseDto, error) {
	var response dto.CnyesStockQuoteResponseDto

	// 建構完整的URL
	url := fmt.Sprintf("https://ws.api.cnyes.com/ws/api/v1/quote/quotes/TWS:%s:STOCK?column=K,E,KEY,M,AI", symbol)

	// 建立請求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, fmt.Errorf("建立請求失敗: %v", err)
	}

	// 設定標頭
	req.Header.Set("Accept", "application/json")

	// 發送請求
	resp, err := c.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("無法連接到鉅亨網 API: %v", err)
	}
	defer resp.Body.Close()

	// 檢查狀態碼
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("鉅亨網 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	// 解析 JSON 回應
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
	}

	return response, nil
}
