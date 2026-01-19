package port

import (
	dto "github.com/tian841224/stock-bot/internal/application/dto"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// Formatter 定義訊息格式化介面，用於將業務資料轉換為使用者可讀的訊息格式。
type FormatterPort interface {
	// FormatDailyMarketInfo 格式化大盤資訊訊息
	FormatDailyMarketInfo(data *[]dto.DailyMarketInfo, userType valueobject.UserType) string

	// FormatStockPerformance 格式化股票績效表現
	FormatStockPerformance(stockName, symbol string, data *[]dto.StockPerformanceData, userType valueobject.UserType) string
	// FormatStockInfo 格式化股票詳細資訊
	FormatStockCompanyInfo(data *dto.StockCompanyInfo, userType valueobject.UserType) string
	// FormatTopVolumeStock 格式化交易量排行資訊
	FormatTopVolumeStock(data *[]dto.TopVolume, userType valueobject.UserType) string

	// FormatStockPriceByDate 格式化指定日期的股價資訊
	FormatStockPrice(data *dto.StockPrice, userType valueobject.UserType) string

	// FormatRevenueMessage 格式化營收資訊
	FormatStockRevenue(data *dto.StockRevenue, userType valueobject.UserType) string

	// FormatChartCaption 格式化圖表標題
	FormatChartCaption(name, symbol, chartType string) string

	// FormatTelegramNewsMessage 格式化 Telegram 股票新聞訊息
	FormatTelegramNewsMessage(news []dto.StockNews, stockName, symbol string) *dto.TgStockNewsMessage

	// FormatLineNewsMessage 格式化 Line 股票新聞訊息
	FormatLineNewsMessage(news []dto.StockNews, stockName, symbol string) *dto.LineStockNewsMessage

	// FormatSubscribed 格式化訂閱股票和項目
	FormatSubscribed(stocks []*dto.UserSubscriptionStock, items []*dto.UserSubscriptionItem) string
}
