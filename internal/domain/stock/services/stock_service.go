package services

import (
	"fmt"
	"time"

	"github.com/tian841224/stock-bot/internal/domain/stock"
)

// StockPerformance 股票績效
type StockPerformance struct {
	StockID     string
	Period      string
	PeriodName  string
	Performance string
}

// StockDomainService 股票領域服務
type StockDomainService struct{}

// NewStockDomainService 建立股票領域服務
func NewStockDomainService() *StockDomainService {
	return &StockDomainService{}
}

// CalculatePriceChange 計算價格變動
func (s *StockDomainService) CalculatePriceChange(current, previous float64) (change float64, changeRate float64) {
	change = current - previous
	if previous != 0 {
		changeRate = (change / previous) * 100
	}
	return change, changeRate
}

// CalculateAmplitude 計算振幅
func (s *StockDomainService) CalculateAmplitude(high, low, prevClose float64) float64 {
	if prevClose == 0 {
		return 0
	}

	amplitude := ((high - low) / prevClose) * 100
	return amplitude
}

// CalculateVolumeRatio 計算週轉率
func (s *StockDomainService) CalculateVolumeRatio(volume, avgVolume float64) float64 {
	if avgVolume == 0 {
		return 0
	}
	return (volume / avgVolume) * 100
}

// ValidateStockData 驗證股票資料
func (s *StockDomainService) ValidateStockData(stock *stock.Stock) error {
	if stock == nil {
		return fmt.Errorf("股票資料不能為空")
	}

	if !stock.IsValid() {
		return fmt.Errorf("股票資料無效")
	}

	if stock.CurrentInfo == nil {
		return fmt.Errorf("缺少價格資訊")
	}

	// 驗證價格資料
	if stock.CurrentInfo.CurrentPrice < 0 {
		return fmt.Errorf("現價不能為負數")
	}

	if stock.CurrentInfo.Volume < 0 {
		return fmt.Errorf("成交量不能為負數")
	}

	return nil
}

// IsTradingTime 檢查是否為交易時間
func (s *StockDomainService) IsTradingTime() bool {
	now := time.Now()

	// 檢查是否為週末
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	// 檢查是否為交易時間 (9:00-13:30)
	hour := now.Hour()
	minute := now.Minute()

	// 上午盤：9:00-12:00
	if hour >= 9 && hour < 12 {
		return true
	}

	// 下午盤：13:30-14:30
	if hour == 13 && minute >= 30 {
		return true
	}
	if hour == 14 && minute <= 30 {
		return true
	}

	return false
}

// GetMarketStatus 取得市場狀態
func (s *StockDomainService) GetMarketStatus() string {
	if s.IsTradingTime() {
		return "交易中"
	}

	now := time.Now()
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return "休市"
	}

	return "收盤"
}

// CalculatePerformance 計算績效
func (s *StockDomainService) CalculatePerformance(stock *stock.Stock, period string) (*StockPerformance, error) {
	if stock == nil {
		return nil, fmt.Errorf("股票資料不能為空")
	}

	// 這裡可以根據不同期間計算績效
	// 實際實作需要根據具體的業務需求
	performance := &StockPerformance{
		StockID:     stock.ID,
		Period:      period,
		PeriodName:  s.getPeriodName(period),
		Performance: s.calculatePerformanceValue(stock, period),
	}

	return performance, nil
}

// getPeriodName 取得期間名稱
func (s *StockDomainService) getPeriodName(period string) string {
	periodMap := map[string]string{
		"1d":  "1日",
		"1w":  "1週",
		"1m":  "1月",
		"3m":  "3月",
		"6m":  "6月",
		"1y":  "1年",
		"ytd": "年初至今",
	}

	if name, exists := periodMap[period]; exists {
		return name
	}
	return period
}

// calculatePerformanceValue 計算績效值
func (s *StockDomainService) calculatePerformanceValue(stock *stock.Stock, period string) string {
	if stock.CurrentInfo == nil {
		return "無資料"
	}

	// 根據期間和價格變動計算績效
	changeRate := stock.CurrentInfo.ChangeRate

	if changeRate > 0 {
		return fmt.Sprintf("上漲 %.2f%%", changeRate)
	} else if changeRate < 0 {
		return fmt.Sprintf("下跌 %.2f%%", changeRate)
	}

	return "持平"
}

// FormatPrice 格式化價格
func (s *StockDomainService) FormatPrice(price float64) string {
	if price >= 1000 {
		return fmt.Sprintf("%.0f", price)
	} else if price >= 100 {
		return fmt.Sprintf("%.1f", price)
	} else if price >= 10 {
		return fmt.Sprintf("%.2f", price)
	} else {
		return fmt.Sprintf("%.3f", price)
	}
}

// FormatVolume 格式化成交量
func (s *StockDomainService) FormatVolume(volume int64) string {
	if volume >= 100000000 {
		return fmt.Sprintf("%.2f億", float64(volume)/100000000)
	} else if volume >= 10000 {
		return fmt.Sprintf("%.2f萬", float64(volume)/10000)
	}
	return fmt.Sprintf("%d", volume)
}

// FormatPercentage 格式化百分比
func (s *StockDomainService) FormatPercentage(value float64) string {
	if value > 0 {
		return fmt.Sprintf("+%.2f%%", value)
	} else if value < 0 {
		return fmt.Sprintf("%.2f%%", value)
	}
	return "0.00%"
}
