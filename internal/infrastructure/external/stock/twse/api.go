package twse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// GetAfterTradingVolume 盤後資訊 - 依股票代碼查詢
func (t *TwseAPI) GetAfterTradingVolume(symbol string, date string) (dto.AfterTradingVolumeRawResponseDto, error) {
	u, err := url.Parse(t.baseURL + "/afterTrading/MI_INDEX")
	if err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, err
	}
	q := u.Query()
	if strings.TrimSpace(date) != "" {
		q.Set("date", date)
	}
	q.Set("type", "ALLBUT0999")
	u.RawQuery = q.Encode()

	t.logger.Debug("twse API request", logger.String("url", u.String()), logger.String("symbol", symbol))

	req, err := t.getRequest(u.String())
	if err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.Error("twse API request failed", logger.String("url", u.String()), logger.Error(err))
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("failed to connect to TWSE API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.logger.Error("twse API returned non-200",
			logger.String("url", u.String()),
			logger.Int("status_code", resp.StatusCode))
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("TWSE API error, status: %d", resp.StatusCode)
	}

	var response dto.AfterTradingVolumeRawResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("failed to parse response JSON: %v", err)
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
