package cnyes

import (
	"fmt"
	"stock-bot/internal/infrastructure/cnyes"
	"stock-bot/internal/infrastructure/cnyes/dto"
	stockDto "stock-bot/internal/service/cnyes/dto"
)

// CnyesServiceInterface 定義鉅亨網服務介面
type CnyesServiceInterface interface {
	GetStockQuote(stockID string) (*stockDto.StockQuoteInfo, error)
	FormatStockQuote(data dto.CnyesStockQuoteDataDto) *stockDto.StockQuoteInfo
}

// CnyesService 鉅亨網服務
type CnyesService struct {
	cnyesAPI cnyes.CnyesAPIInterface
}

// NewCnyesService 建立新的鉅亨網服務
func NewCnyesService() *CnyesService {
	return &CnyesService{
		cnyesAPI: cnyes.NewCnyesAPI(),
	}
}

// GetStockQuote 取得股票報價資訊
func (s *CnyesService) GetStockQuote(stockID string) (*stockDto.StockQuoteInfo, error) {
	// 建構股票符號 (格式: TWS:2330:STOCK)
	symbol := fmt.Sprintf("TWS:%s:STOCK", stockID)

	// 呼叫API
	response, err := s.cnyesAPI.GetStockQuote(symbol)
	if err != nil {
		return nil, fmt.Errorf("取得股票報價失敗: %v", err)
	}

	// 檢查回應
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("API回應錯誤: %s", response.Message)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("查無股票資料: %s", stockID)
	}

	// 格式化資料（取第一筆）
	stockInfo := s.FormatStockQuote(response.Data[0])
	return stockInfo, nil
}

// FormatStockQuote 格式化股票報價資料
func (s *CnyesService) FormatStockQuote(data dto.CnyesStockQuoteDataDto) *stockDto.StockQuoteInfo {
	return &stockDto.StockQuoteInfo{
		// 基本資訊
		StockID:   data.StockID,
		StockName: data.StockName,
		Industry:  data.Industry,
		Market:    data.Market,

		// 價格資訊
		CurrentPrice: data.CurrentPrice,
		Change:       data.Change,
		ChangeRate:   data.ChangeRate,
		OpenPrice:    data.OpenPrice,
		HighPrice:    data.HighPrice,
		LowPrice:     data.LowPrice,
		PrevClose:    data.PrevClose,

		// 成交量資訊 (轉換單位)
		Volume:      int64(data.Volume),
		Turnover:    data.Turnover / 1e8,    // 轉換為億元
		VolumeRatio: data.VolumeRatio * 100, // 轉換為百分比
		Amplitude:   data.Amplitude,

		// 財務指標
		PE:           data.PE,
		PB:           data.PB,
		MarketCap:    data.MarketCap / 1e12, // 轉換為兆元
		BookValue:    data.BookValue,
		EPS:          data.EPS,
		QuarterEPS:   data.QuarterEPS,
		Dividend:     data.Dividend,
		DividendRate: data.DividendRate,
		GrossMargin:  data.GrossMargin,
		OperMargin:   data.OperMargin,
		NetMargin:    data.NetMargin,

		// 價位區間
		UpperLimit:  data.UpperLimit,
		LowerLimit:  data.LowerLimit,
		High52W:     data.High52W,
		Low52W:      data.Low52W,
		High52WDate: data.High52WDate,
		Low52WDate:  data.Low52WDate,

		// 五檔資訊
		BidPrices: []float64{
			data.BidPrice1, data.BidPrice2, data.BidPrice3, data.BidPrice4, data.BidPrice5,
		},
		AskPrices: []float64{
			data.AskPrice1, data.AskPrice2, data.AskPrice3, data.AskPrice4, data.AskPrice5,
		},

		// 內外盤資訊
		OutVolume: int64(data.OutVolume),
		InVolume:  int64(data.InVolume),
		OutRatio:  data.OutRatio,
		InRatio:   data.InRatio,
	}
}
