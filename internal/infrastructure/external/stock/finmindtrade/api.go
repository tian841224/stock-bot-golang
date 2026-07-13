package finmindtrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade/dto"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type FinmindTradeAPI struct {
	baseURL    string
	httpHeader http.Header
	client     *http.Client
	logger     logger.Logger
}

func NewFinmindTradeAPI(cfg config.Config, log logger.Logger) *FinmindTradeAPI {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	if cfg.FINMIND_TOKEN != "" {
		header.Set("Authorization", "Bearer "+cfg.FINMIND_TOKEN)
	}
	return &FinmindTradeAPI{
		baseURL:    "https://api.finmindtrade.com/api/v4/data",
		client:     &http.Client{Timeout: 10 * time.Second},
		httpHeader: header,
		logger:     log,
	}
}

// GetTaiwanStockInfo 取得台灣股票資訊
func (f *FinmindTradeAPI) GetTaiwanStockInfo() (response dto.TaiwanStockInfoResponseDto, err error) {
	requestDto := dto.FinmindtradeRequestDto{
		DataSet: "TaiwanStockInfo",
	}
	return doRequest[dto.TaiwanStockInfoResponseDto](f, requestDto)
}

// GetTaiwanStockPrice 取得台灣股票價格
func (f *FinmindTradeAPI) GetTaiwanStockPrice(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockPriceResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockPrice"
	return doRequest[dto.TaiwanStockPriceResponseDto](f, requestDto)
}

// GetTaiwanStockTradingDate 台股交易日
func (f *FinmindTradeAPI) GetTaiwanStockTradingDate(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockTradingDateResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockTradingDate"
	return doRequest[dto.TaiwanStockTradingDateResponseDto](f, requestDto)
}

// GetTaiwanStockSplitPrice 台股分割股價
func (f *FinmindTradeAPI) GetTaiwanStockSplitPrice(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockSplitPriceResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockSplitPrice"
	return doRequest[dto.TaiwanStockSplitPriceResponseDto](f, requestDto)
}

func (f *FinmindTradeAPI) GetTaiwanStockNews(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanNewsResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockNews"
	return doRequest[dto.TaiwanNewsResponseDto](f, requestDto)
}

// GetUSStockInfo 美股股票清單
func (f *FinmindTradeAPI) GetUSStockInfo() (response dto.USStockInfoResponseDto, err error) {
	requestDto := dto.FinmindtradeRequestDto{
		DataSet: "USStockInfo",
	}
	return doRequest[dto.USStockInfoResponseDto](f, requestDto)
}

// GetTodayInfo 大盤資訊(法人/資券/美股大盤)
func (f *FinmindTradeAPI) GetTodayInfo() (response dto.TodayInfoResponseDto, err error) {
	baseURL := "https://api.web.finmindtrade.com/v2/today_info"
	f.logger.Debug("finmind API request", logger.String("url", baseURL))

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return response, err
	}
	req.Header = f.httpHeader

	resp, err := f.client.Do(req)
	if err != nil {
		f.logger.Error("finmind API request failed", logger.String("url", baseURL), logger.Error(err))
		return response, fmt.Errorf("failed to connect to finmind API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		f.logger.Error("finmind API returned non-200",
			logger.String("url", baseURL),
			logger.Int("status_code", resp.StatusCode))
		return response, fmt.Errorf("finmind API error, status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to parse response JSON: %v", err)
	}
	return response, nil
}

// doRequest 共用方法：送出請求並解析 JSON 至指定型別
func doRequest[T any](f *FinmindTradeAPI, requestDto dto.FinmindtradeRequestDto) (response T, err error) {
	req, err := f.getRequest()
	if err != nil {
		return response, fmt.Errorf("failed to build request: %v", err)
	}
	query := req.URL.Query()
	if requestDto.DataSet != "" {
		query.Add("dataset", requestDto.DataSet)
	}
	if requestDto.StockID != "" {
		query.Add("stock_id", requestDto.StockID)
	}
	if requestDto.DataID != "" {
		query.Add("data_id", requestDto.DataID)
	}
	if requestDto.StartDate != "" {
		query.Add("start_date", requestDto.StartDate)
	}
	if requestDto.EndDate != "" {
		query.Add("end_date", requestDto.EndDate)
	}
	req.URL.RawQuery = query.Encode()

	f.logger.Debug("finmind API request",
		logger.String("dataset", requestDto.DataSet),
		logger.String("stock_id", requestDto.StockID))

	resp, err := f.client.Do(req)
	if err != nil {
		f.logger.Error("finmind API request failed",
			logger.String("dataset", requestDto.DataSet),
			logger.Error(err))
		return response, fmt.Errorf("failed to connect to finmind API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		f.logger.Error("finmind API returned non-200",
			logger.String("dataset", requestDto.DataSet),
			logger.Int("status_code", resp.StatusCode))
		return response, fmt.Errorf("finmind API error, status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return response, nil
}

// 設定Request參數
func (f *FinmindTradeAPI) getRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", f.baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = f.httpHeader

	return req, nil
}
