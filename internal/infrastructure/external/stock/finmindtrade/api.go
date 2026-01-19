package finmindtrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/stock/finmindtrade/dto"
)

type FinmindTradeAPI struct {
	baseURL    string
	httpHeader http.Header
	client     *http.Client
}

func NewFinmindTradeAPI(cfg config.Config) *FinmindTradeAPI {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	if cfg.FINMIND_TOKEN != "" {
		header.Set("Authorization", "Bearer "+cfg.FINMIND_TOKEN)
	}
	return &FinmindTradeAPI{
		baseURL:    "https://api.finmindtrade.com/api/v4/data",
		client:     &http.Client{Timeout: 10 * time.Second},
		httpHeader: header,
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

// GetTaiwanExchangeRate 取得台灣匯率
func (f *FinmindTradeAPI) GetTaiwanExchangeRate(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanExchangeRateResponseDto, err error) {
	requestDto.DataSet = "TaiwanExchangeRate"
	return doRequest[dto.TaiwanExchangeRateResponseDto](f, requestDto)
}

// GetTaiwanStockDividend 取得台灣股票股利
func (f *FinmindTradeAPI) GetTaiwanStockDividend(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockDividendResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockDividend"
	return doRequest[dto.TaiwanStockDividendResponseDto](f, requestDto)
}

// GetTaiwanStockFinancialStatements 綜合損益表
func (f *FinmindTradeAPI) GetTaiwanStockFinancialStatements(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockFinancialStatementsResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockFinancialStatements"
	return doRequest[dto.TaiwanStockFinancialStatementsResponseDto](f, requestDto)
}

// GetTaiwanStockMonthRevenue 月營收表
func (f *FinmindTradeAPI) GetTaiwanStockMonthRevenue(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockMonthRevenueResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockMonthRevenue"
	return doRequest[dto.TaiwanStockMonthRevenueResponseDto](f, requestDto)
}

// GetTaiwanStockTradingDate 台股交易日
func (f *FinmindTradeAPI) GetTaiwanStockTradingDate(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockTradingDateResponseDto, err error) {
	requestDto.DataSet = "TaiwanStockTradingDate"
	return doRequest[dto.TaiwanStockTradingDateResponseDto](f, requestDto)
}

// GetTaiwanVariousIndicators 台股各種指標(每5秒)
func (f *FinmindTradeAPI) GetTaiwanVariousIndicators(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanVariousIndicatorsResponseDto, err error) {
	requestDto.DataSet = "TaiwanVariousIndicators5Seconds"
	return doRequest[dto.TaiwanVariousIndicatorsResponseDto](f, requestDto)
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

// GetUSStockPrice 美股盤後股價
func (f *FinmindTradeAPI) GetUSStockPrice(requestDto dto.FinmindtradeRequestDto) (response dto.USStockPriceResponseDto, err error) {
	requestDto.DataSet = "USStockPrice"
	return doRequest[dto.USStockPriceResponseDto](f, requestDto)
}

// GetTodayInfo 大盤資訊(法人/資券/美股大盤)
func (f *FinmindTradeAPI) GetTodayInfo() (response dto.TodayInfoResponseDto, err error) {
	baseURL := "https://api.web.finmindtrade.com/v2/today_info"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return response, err
	}
	req.Header = f.httpHeader

	resp, err := f.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
	}
	return response, nil
}

// GetTaiwanStockAnalysis 取得台灣股票分析
func (f *FinmindTradeAPI) GetTaiwanStockAnalysis(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockAnalysisResponseDto, err error) {
	baseURL := "https://api.web.finmindtrade.com/v2/taiwan_stock_analysis"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return response, err
	}
	req.Header = f.httpHeader

	query := req.URL.Query()
	if requestDto.StockID != "" {
		query.Add("stock_id", requestDto.StockID)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
	}
	return response, nil
}

func (f *FinmindTradeAPI) GetTaiwanStockAnalysisPlot(requestDto dto.FinmindtradeRequestDto) (response dto.TaiwanStockAnalysisPlotResponseDto, err error) {
	baseURL := "https://api.web.finmindtrade.com/v2/taiwan_stock_analysis_plot"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return response, err
	}
	req.Header = f.httpHeader

	query := req.URL.Query()
	if requestDto.StockID != "" {
		query.Add("stock_id", requestDto.StockID)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
	}
	return response, nil
}

// doRequest 共用方法：送出請求並解析 JSON 至指定型別
func doRequest[T any](f *FinmindTradeAPI, requestDto dto.FinmindtradeRequestDto) (response T, err error) {
	req, err := f.getRequest()
	if err != nil {
		return response, fmt.Errorf("無法建立Request: %v", err)
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
	resp, err := f.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("無法連接到外部 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
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
