package twse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"stock-bot/internal/infrastructure/twse/dto"
	"strings"
	"time"
)

type TwseAPI struct {
	baseURL string
	client  *http.Client
}

func NewTwseAPI() *TwseAPI {
	return &TwseAPI{
		baseURL: "https://www.twse.com.tw/rwd/zh",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// 成交量前 20 股票
func (t *TwseAPI) GetTopVolumeItems() (dto.TopVolumeItemsResponseDto, error) {
	urlStr := t.baseURL + "/afterTrading/MI_INDEX20"
	req, err := t.getRequest(urlStr)
	if err != nil {
		return dto.TopVolumeItemsResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	var response dto.TopVolumeItemsResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.TopVolumeItemsResponseDto{}, fmt.Errorf("無法解析回應 JSON: %v", err)
	}
	return response, nil
}

// 盤後資訊 - 依股票代碼查詢
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

	req, err := t.getRequest(u.String())
	if err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	var response dto.AfterTradingVolumeRawResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.AfterTradingVolumeRawResponseDto{}, fmt.Errorf("無法解析回應 JSON: %v", err)
	}
	return response, nil
}

// 取得大盤每日成交資訊
func (t *TwseAPI) GetDailyMarketInfo() (dto.DailyMarketInfoResponseDto, error) {
	urlStr := t.baseURL + "/afterTrading/FMTQIK"
	req, err := t.getRequest(urlStr)
	if err != nil {
		return dto.DailyMarketInfoResponseDto{}, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	var response dto.DailyMarketInfoResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.DailyMarketInfoResponseDto{}, fmt.Errorf("無法解析回應 JSON: %v", err)
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
