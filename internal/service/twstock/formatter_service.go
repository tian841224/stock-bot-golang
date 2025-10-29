package twstock

import (
	"fmt"
	"strings"

	"stock-bot/internal/domain/stock"
	stockDto "stock-bot/internal/service/twstock/dto"
)

// FormatterService 格式化服務
type FormatterService struct {
	domainService *DomainService
}

// NewFormatterService 建立格式化服務
func NewFormatterService(domainService *DomainService) *FormatterService {
	return &FormatterService{
		domainService: domainService,
	}
}

// FormatStockInfoForDisplay 格式化股票資訊用於顯示
func (f *FormatterService) FormatStockInfoForDisplay(stock *stock.Stock) string {
	if stock == nil {
		return "無股票資料"
	}

	var lines []string

	// 基本資訊
	lines = append(lines, fmt.Sprintf("📊 %s (%s)", stock.Name, stock.ID))
	lines = append(lines, fmt.Sprintf("🏢 %s | %s", stock.Industry, stock.Market))
	lines = append(lines, "")

	// 價格資訊
	if stock.CurrentInfo != nil {
		price := f.domainService.GetStockDomainService().FormatPrice(stock.CurrentInfo.CurrentPrice)
		changeRate := f.domainService.GetStockDomainService().FormatPercentage(stock.CurrentInfo.ChangeRate)
		volume := f.domainService.GetStockDomainService().FormatVolume(stock.CurrentInfo.Volume)

		lines = append(lines, fmt.Sprintf("💰 現價: %s", price))
		lines = append(lines, fmt.Sprintf("📈 漲跌: %s", changeRate))
		lines = append(lines, fmt.Sprintf("📊 成交量: %s", volume))

		if stock.CurrentInfo.Turnover > 0 {
			turnover := stock.GetTurnoverInBillions()
			lines = append(lines, fmt.Sprintf("💵 成交額: %.2f億", turnover))
		}
		lines = append(lines, "")
	}

	// 財務指標
	if stock.Financials != nil {
		lines = append(lines, "📋 財務指標:")
		lines = append(lines, fmt.Sprintf("  本益比: %.2f", stock.Financials.PE))
		lines = append(lines, fmt.Sprintf("  本淨比: %.2f", stock.Financials.PB))
		lines = append(lines, fmt.Sprintf("  EPS: %.2f", stock.Financials.EPS))
		lines = append(lines, fmt.Sprintf("  殖利率: %.2f%%", stock.Financials.DividendRate))

		if stock.Financials.MarketCap > 0 {
			marketCap := stock.GetMarketCapInTrillions()
			lines = append(lines, fmt.Sprintf("  市值: %.2f兆", marketCap))
		}
		lines = append(lines, "")
	}

	// 狀態資訊
	status := stock.GetPriceChangeStatus()
	lines = append(lines, fmt.Sprintf("📈 狀態: %s", status))

	if stock.IsTradingDay() {
		lines = append(lines, "🟢 今日交易")
	} else {
		lines = append(lines, "🔴 非交易日")
	}

	return strings.Join(lines, "\n")
}

// FormatPerformanceTable 格式化績效表格
func (f *FormatterService) FormatPerformanceTable(stockName, symbol string, performanceData []stockDto.StockPerformanceData) string {
	if len(performanceData) == 0 {
		return fmt.Sprintf("📊 %s (%s)\n無績效資料", stockName, symbol)
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("📊 %s (%s) 績效分析", stockName, symbol))
	lines = append(lines, "="+strings.Repeat("=", 30))

	for _, data := range performanceData {
		lines = append(lines, fmt.Sprintf("%s: %s", data.PeriodName, data.Performance))
	}

	return strings.Join(lines, "\n")
}

// FormatRevenueInfo 格式化營收資訊
func (f *FormatterService) FormatRevenueInfo(revenue *stock.Revenue) string {
	if revenue == nil || len(revenue.Time) == 0 {
		return "無營收資料"
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("📈 股票代號: %s", revenue.StockID))
	lines = append(lines, "")

	// 顯示最新幾期營收
	displayCount := 6
	if len(revenue.Time) < displayCount {
		displayCount = len(revenue.Time)
	}

	lines = append(lines, "📊 最近營收資料:")
	for i := len(revenue.Time) - displayCount; i < len(revenue.Time); i++ {
		period := fmt.Sprintf("%d/%02d",
			revenue.Time[i]/10000,
			(revenue.Time[i]%10000)/100)

		revenueValue := int64(0)
		yoy := 0.0

		if i < len(revenue.SaleMonth) {
			revenueValue = revenue.SaleMonth[i]
		}
		if i < len(revenue.YoY) {
			yoy = revenue.YoY[i]
		}

		lines = append(lines, fmt.Sprintf("  %s: %s (年增率: %.1f%%)",
			period,
			f.domainService.GetStockDomainService().FormatVolume(revenueValue),
			yoy))
	}

	return strings.Join(lines, "\n")
}

// FormatMarketStatus 格式化市場狀態
func (f *FormatterService) FormatMarketStatus() string {
	status := f.domainService.GetStockDomainService().GetMarketStatus()

	switch status {
	case "交易中":
		return "🟢 市場交易中"
	case "休市":
		return "🔴 市場休市"
	case "收盤":
		return "🟡 市場收盤"
	default:
		return "❓ 市場狀態未知"
	}
}
