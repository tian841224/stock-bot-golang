package fugle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/fugle/dto"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// FugleAPI 定義 Fugle API 的實作
type FugleAPI struct {
	baseURL    string
	client     *http.Client
	httpHeader http.Header
	logger     logger.Logger
}

// NewFugleAPI 建立新的 FugleAPI 實例
func NewFugleAPI(cfg config.Config, log logger.Logger) *FugleAPI {
	return &FugleAPI{
		baseURL: "https://api.fugle.tw/marketdata/v1.0/stock/",
		client:  &http.Client{Timeout: 10 * time.Second},
		httpHeader: http.Header{
			"X-API-KEY": []string{cfg.FUGLE_API_KEY},
		},
		logger: log,
	}
}

// GetStockHistoricalCandles 取得股票歷史Ｋ線
func (f *FugleAPI) GetStockHistoricalCandles(requestDto dto.FugleCandlesRequestDto) (dto.FugleCandlesResponseDto, error) {
	apiURL := f.baseURL + "/historical/candles/" + requestDto.Symbol
	params := url.Values{}
	if requestDto.Timeframe != "" {
		params.Add("timeframe", requestDto.Timeframe)
	}
	if requestDto.From != "" {
		params.Add("from", requestDto.From)
	}
	if requestDto.To != "" {
		params.Add("to", requestDto.To)
	}
	if requestDto.Fields != "" {
		params.Add("fields", requestDto.Fields)
	}
	if requestDto.Sort != "" {
		params.Add("sort", requestDto.Sort)
	}
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}
	return getResponse[dto.FugleCandlesResponseDto](f, apiURL)
}

func getResponse[T any](c *FugleAPI, url string) (response T, err error) {
	c.logger.Debug("fugle API request", logger.String("url", url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}

	req.Header = c.httpHeader

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("fugle API request failed", logger.String("url", url), logger.Error(err))
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("fugle API returned non-200",
			logger.String("url", url),
			logger.Int("status_code", resp.StatusCode))
		return response, fmt.Errorf("fugle API error, status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read API response: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return response, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return response, nil
}
