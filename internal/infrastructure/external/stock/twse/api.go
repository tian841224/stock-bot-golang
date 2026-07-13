package twse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/twse/dto"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type TwseAPI struct {
	baseURL string
	client  *http.Client
	logger  logger.Logger
}

func NewTwseAPI(log logger.Logger) *TwseAPI {
	return &TwseAPI{
		baseURL: "https://www.twse.com.tw/rwd/zh",
		client:  &http.Client{Timeout: 10 * time.Second},
		logger:  log,
	}
}

// GetTopVolumeItems 成交量前 20 股票
func (t *TwseAPI) GetTopVolumeItems() (dto.TopVolumeItemsResponseDto, error) {
	urlStr := t.baseURL + "/afterTrading/MI_INDEX20"
	t.logger.Debug("twse API request", logger.String("url", urlStr))

	req, err := t.getRequest(urlStr)
	if err != nil {
		return dto.TopVolumeItemsResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.Error("twse API request failed", logger.String("url", urlStr), logger.Error(err))
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("failed to connect to TWSE API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.logger.Error("twse API returned non-200",
			logger.String("url", urlStr),
			logger.Int("status_code", resp.StatusCode))
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("TWSE API error, status: %d", resp.StatusCode)
	}

	var response dto.TopVolumeItemsResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("failed to parse response JSON: %v", err)
	}
	return response, nil
}

// GetDailyMarketInfo 取得大盤每日成交資訊
func (t *TwseAPI) GetDailyMarketInfo() (dto.DailyMarketInfoResponseDto, error) {
	urlStr := t.baseURL + "/afterTrading/FMTQIK"
	t.logger.Debug("twse API request", logger.String("url", urlStr))

	req, err := t.getRequest(urlStr)
	if err != nil {
		return dto.DailyMarketInfoResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.Error("twse API request failed", logger.String("url", urlStr), logger.Error(err))
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("failed to connect to TWSE API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.logger.Error("twse API returned non-200",
			logger.String("url", urlStr),
			logger.Int("status_code", resp.StatusCode))
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("TWSE API error, status: %d", resp.StatusCode)
	}

	var response dto.DailyMarketInfoResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("failed to parse response JSON: %v", err)
	}
	return response, nil
}

// 設定Request參數
func (f *TwseAPI) getRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
