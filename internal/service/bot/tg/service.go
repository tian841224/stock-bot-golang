package tg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	fugleDto "github.com/tian841224/stock-bot/internal/infrastructure/fugle/dto"
	twseDto "github.com/tian841224/stock-bot/internal/infrastructure/twse/dto"
	"github.com/tian841224/stock-bot/internal/repository"
	tgDto "github.com/tian841224/stock-bot/internal/service/bot/tg/dto"
	"github.com/tian841224/stock-bot/internal/service/twstock"
	stockDto "github.com/tian841224/stock-bot/internal/service/twstock/dto"
	"github.com/tian841224/stock-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgService struct {
	stockService         *twstock.StockService
	userSubscriptionRepo repository.UserSubscriptionRepository
}

func NewTgService(
	stockService *twstock.StockService,
	userSubscriptionRepo repository.UserSubscriptionRepository,
) *TgService {
	return &TgService{
		stockService:         stockService,
		userSubscriptionRepo: userSubscriptionRepo,
	}
}

// GetDailyMarketInfo 取得大盤資訊
func (s *TgService) GetDailyMarketInfo(count int) (string, error) {
	marketInfo, err := s.stockService.GetDailyMarketInfo(count)
	if err != nil {
		logger.Log.Error("取得大盤資訊失敗", zap.Error(err))
		return "", fmt.Errorf("查無資料，請確認後再試")
	}
	return s.formatDailyMarketInfoMessage(marketInfo), nil
}

// GetStockPerformance 取得股票績效
func (s *TgService) GetStockPerformance(symbol string) (string, error) {
	// 驗證股票代號並取得基本資訊
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return "", fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 取得績效
	performanceData, err := s.stockService.GetStockPerformance(symbol)
	if err != nil {
		logger.Log.Error("取得股票績效失敗", zap.Error(err))
		return "", fmt.Errorf("取得績效資料失敗，請稍後再試")
	}

	// 格式化績效資料為文字表格
	formattedText := s.formatPerformanceTable(stockName, symbol, performanceData)

	return formattedText, nil
}

// GetStockPerformanceWithChart 取得股票績效並生成圖表
func (s *TgService) GetStockPerformanceWithChart(symbol string, chartType string) ([]byte, string, error) {
	// 驗證股票代號並取得基本資訊
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return nil, "", fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 取得績效和圖表
	performanceChartData, err := s.stockService.GetStockPerformanceWithChart(symbol, chartType)
	if err != nil {
		logger.Log.Error("取得股票績效失敗", zap.Error(err))
		return nil, "", fmt.Errorf("取得績效資料失敗，請稍後再試")
	}

	// 取得績效
	performanceData, err := s.stockService.GetStockPerformance(symbol)
	if err != nil {
		logger.Log.Error("取得股票績效失敗", zap.Error(err))
		return nil, "", fmt.Errorf("取得績效資料失敗，請稍後再試")
	}

	// 格式化績效資料為文字表格
	formattedText := s.formatPerformanceTable(stockName, symbol, performanceData)

	return performanceChartData.ChartData, formattedText, nil
}

// GetTopVolumeItemsFormatted 取得格式化的交易量前20名
func (s *TgService) GetTopVolumeItemsFormatted() (string, error) {
	topItems, err := s.stockService.GetTopVolumeItems()
	if err != nil {
		logger.Log.Error("取得交易量前20名失敗", zap.Error(err))
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	if len(topItems) == 0 {
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	messageText := "🔝<b>今日交易量前二十</b>\n\n"

	for _, item := range topItems {
		emoji := ""
		switch item.UpDownSign {
		case "+":
			emoji = "📈"
		case "-":
			emoji = "📉"
		default:
			emoji = ""
		}

		messageText += fmt.Sprintf("%s<b>%s (%s)</b>\n<code>", emoji, item.StockName, item.StockID)
		messageText += fmt.Sprintf("成交股數：%s\n", item.Volume)
		messageText += fmt.Sprintf("成交筆數：%s\n", item.Transaction)
		messageText += fmt.Sprintf("開盤價：%.2f\n", item.OpenPrice)
		messageText += fmt.Sprintf("收盤價：%.2f\n", item.ClosePrice)
		messageText += fmt.Sprintf("漲跌幅：%s%.2f (%s)\n", item.UpDownSign, item.ChangeAmount, item.PercentageChange)
		messageText += fmt.Sprintf("最高價：%.2f\n", item.HighPrice)
		messageText += fmt.Sprintf("最低價：%.2f\n", item.LowPrice)
		messageText += "</code>\n"
	}

	return messageText, nil
}

// GetStockPriceByDate 取得指定日期的股價資訊
func (s *TgService) GetStockPriceByDate(symbol, date string) (string, error) {
	// 取得指定日期股價資訊
	stockInfo, err := s.stockService.GetStockPrice(symbol, date)
	if err != nil {
		logger.Log.Error("取得股價資訊失敗", zap.Error(err))
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	// 格式化日期顯示
	var displayDate string
	if date != "" && len(date) == 8 {
		displayDate = fmt.Sprintf("%s/%s/%s", date[0:4], date[4:6], date[6:8])
	} else {
		t, _ := time.Parse("2006-01-02", stockInfo.Date)
		displayDate = t.Format("2006/01/02")
	}

	emoji := ""
	switch stockInfo.UpDownSign {
	case "+":
		emoji = "📈"
	case "-":
		emoji = "📉"
	default:
		emoji = ""
	}

	message := fmt.Sprintf(`<b>%s</b>
<b>─── %s (%s) %s ───</b>
<code>開盤價：%.2f
收盤價：%.2f
漲跌幅：%s%.2f (%s)
最高價：%.2f
最低價：%.2f
成交股數：%s
成交筆數：%s</code>`,
		displayDate,
		stockInfo.StockName, stockInfo.StockID, emoji,
		stockInfo.OpenPrice,
		stockInfo.ClosePrice,
		stockInfo.UpDownSign, stockInfo.ChangeAmount, stockInfo.PercentageChange,
		stockInfo.HighPrice,
		stockInfo.LowPrice,
		stockInfo.Volume,
		stockInfo.Transaction)

	return message, nil
}

// GetStockInfo 取得股票詳細資訊
func (s *TgService) GetStockInfo(symbol string) (string, error) {
	stockInfo, err := s.stockService.GetStockInfo(symbol)
	if err != nil {
		logger.Log.Error("取得股票詳細資訊失敗", zap.Error(err))
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	message := s.formatStockInfoMessage(stockInfo)
	return message, nil
}

// GetStockRevenueWithChart 取得股票財報和圖表
func (s *TgService) GetStockRevenueWithChart(symbol string) ([]byte, string, error) {
	revenue, err := s.stockService.GetStockRevenue(symbol)
	if err != nil {
		logger.Log.Error("取得股票財報失敗", zap.Error(err))
		return nil, "", fmt.Errorf("查無資料，請確認後再試")
	}

	chart, err := s.stockService.GetStockRevenueChart(symbol)
	if err != nil {
		logger.Log.Error("取得股票財報圖表失敗", zap.Error(err))
		return nil, "", fmt.Errorf("查無資料，請確認後再試")
	}

	message := s.formatRevenueMessage(revenue)
	return chart, message, nil
}

// GetStockHistoricalCandlesChart 取得股票歷史K線圖
func (s *TgService) GetStockHistoricalCandlesChart(symbol string) ([]byte, string, error) {
	dto := fugleDto.FugleCandlesRequestDto{
		Symbol: symbol,
		From:   time.Now().AddDate(-1, 0, 1).Format("2006-01-02"),
		// To:        time.Now().Format("2006-01-02"),
		Timeframe: "D",
		Fields:    "open,high,low,close,volume",
	}

	chart, stockName, err := s.stockService.GetStockHistoricalCandlesChart(dto)
	if err != nil {
		logger.Log.Error("取得股票歷史K線圖失敗", zap.Error(err))
		return nil, "", fmt.Errorf("查無資料，請確認後再試")
	}

	caption := fmt.Sprintf("⚡️%s(%s)-歷史K線圖", stockName, symbol)
	return chart, caption, nil
}

// GetTaiwanStockNews 取得股票新聞
func (s *TgService) GetTaiwanStockNews(symbol string) (*tgDto.StockNewsMessage, error) {
	// 驗證股票代號
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return nil, fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 取得新聞
	news, err := s.stockService.GetStockNews(symbol)
	if err != nil {
		logger.Log.Error("取得股票新聞失敗", zap.Error(err))
		return nil, fmt.Errorf("取得新聞失敗，請稍後再試")
	}

	if len(news) == 0 {
		return &tgDto.StockNewsMessage{
			Text: fmt.Sprintf("⚡️%s(%s)-即時新聞\n\n暫無新聞資料", stockName, symbol),
		}, nil
	}

	// 建立新聞按鈕
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, n := range news {
		btn := tgbotapi.NewInlineKeyboardButtonURL(n.Title, n.Link)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// 組合訊息
	message := &tgDto.StockNewsMessage{
		Text:                 fmt.Sprintf("⚡️%s(%s)-即時新聞", stockName, symbol),
		InlineKeyboardMarkup: &keyboard,
	}

	return message, nil
}

// AddUserStockSubscription 新增使用者股票訂閱
func (s *TgService) AddUserStockSubscription(userID uint, symbol string) (string, error) {
	// 驗證股票代號
	valid, _, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return "", fmt.Errorf("無此股票代號，請重新確認")
	}

	// 新增股票訂閱
	success, err := s.userSubscriptionRepo.AddUserSubscriptionStock(userID, symbol)
	if err != nil {
		logger.Log.Error("新增股票訂閱失敗", zap.Error(err))
		return "", fmt.Errorf("訂閱失敗，請稍後再試")
	}

	if !success {
		return "已訂閱過此股票", nil
	}

	return "訂閱成功", nil
}

// DeleteUserStockSubscription 刪除使用者股票訂閱
func (s *TgService) DeleteUserStockSubscription(userID uint, symbol string) (string, error) {
	// 刪除股票訂閱
	success, err := s.userSubscriptionRepo.DeleteUserSubscriptionStock(userID, symbol)
	if err != nil {
		logger.Log.Error("刪除股票訂閱失敗", zap.Error(err))
		return "", fmt.Errorf("取消訂閱失敗，請稍後再試")
	}

	if !success {
		return "取消訂閱失敗，請檢查是否已訂閱", nil
	}

	return "取消訂閱成功", nil
}

// GetUserSubscriptionList 取得使用者訂閱清單
func (s *TgService) GetUserSubscriptionList(userID uint) (string, error) {
	// 取得使用者訂閱項目
	subscriptions, err := s.userSubscriptionRepo.GetUserSubscriptionList(userID)
	if err != nil {
		logger.Log.Error("取得使用者訂閱項目失敗", zap.Error(err))
		return "", fmt.Errorf("取得訂閱清單失敗")
	}

	// 取得使用者訂閱股票
	subscriptionStocks, err := s.userSubscriptionRepo.GetUserSubscriptionStockList(userID)
	if err != nil {
		logger.Log.Error("取得使用者訂閱股票失敗", zap.Error(err))
		return "", fmt.Errorf("取得訂閱清單失敗")
	}

	// 組合訊息
	messageText := "📋 <b>您目前的訂閱項目</b>\n\n"

	// 訂閱功能清單
	messageText += "🔔 <b>已訂閱功能：</b>\n"
	hasActiveSubscriptions := false
	for _, sub := range subscriptions {
		if sub.Status == "active" && sub.Feature != nil {
			messageText += fmt.Sprintf("• %s\n", sub.Feature.Name)
			hasActiveSubscriptions = true
		}
	}
	if !hasActiveSubscriptions {
		messageText += "• 尚未訂閱任何功能\n"
	}

	// 訂閱股票清單
	messageText += "\n📈 <b>已訂閱股票：</b>\n"
	if len(subscriptionStocks) > 0 {
		for _, stock := range subscriptionStocks {
			if stock.Status == 1 {
				messageText += fmt.Sprintf("• %s\n", stock.Stock)
			}
		}
	} else {
		messageText += "• 尚未訂閱任何股票\n"
	}

	return messageText, nil
}

// GetDailyMarketInfo 取得大盤資訊
// func (s *TgService) GetDailyMarketInfo(count int) (string, error) {
// 	marketInfoList, err := s.stockService.GetDailyMarketInfo(count)
// 	if err != nil {
// 		logger.Log.Error("取得大盤資訊失敗", zap.Error(err))
// 		return "", fmt.Errorf("查無資料，請確認後再試")
// 	}

// 	if len(marketInfoList) == 0 {
// 		return "", fmt.Errorf("查無資料，請確認後再試")
// 	}

// 	messageText := "<b>台灣股市大盤資訊</b>\n\n"
// 	for _, row := range marketInfoList {
// 		messageText += fmt.Sprintf("<b>%s</b>\n", row.Date)
// 		messageText += "<code>"
// 		messageText += fmt.Sprintf("成交股數：%s\n", row.Volume)
// 		messageText += fmt.Sprintf("成交金額：%s\n", row.Amount)
// 		messageText += fmt.Sprintf("成交筆數：%s\n", row.Transaction)
// 		messageText += fmt.Sprintf("發行量加權股價指數：%s\n", row.Index)
// 		messageText += fmt.Sprintf("漲跌點數：%s\n", row.Change)
// 		messageText += "</code>\n"
// 	}

// 	return messageText, nil
// }

// formatRevenueMessage 格式化股票財報訊息
func (s *TgService) formatRevenueMessage(revenue *stockDto.RevenueDto) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("<b>📊 %s(%s) 月營收</b>\n\n", revenue.Name, revenue.Code))

	// 檢查是否有資料
	if len(revenue.SaleMonth) == 0 || len(revenue.YoY) == 0 {
		message.WriteString("❌ 暫無營收資料")
		return message.String()
	}

	message.WriteString("<pre>")
	// 顯示所有資料
	for i := 0; i < len(revenue.Time); i++ {
		timeStr := s.formatTimeFromTimestamp(revenue.Time[i])

		// 營收(千元) -> 億元
		monthRevenueE := float64(revenue.SaleMonth[i]) / 100000.0

		// 年增率
		yoy := revenue.YoY[i]

		// 累計營收(千元) -> 億元
		accumulatedRevenueE := float64(revenue.SaleAccumulated[i]) / 100000.0

		// 累計年增率
		accumulatedYoY := revenue.YoYAccumulated[i]

		message.WriteString(fmt.Sprintf("---%s---\n", timeStr))
		message.WriteString(fmt.Sprintf("營收(億元): %.2f\n", monthRevenueE))
		message.WriteString(fmt.Sprintf("年增率: %.2f%%\n", yoy))
		message.WriteString(fmt.Sprintf("累計營收(億元): %.2f\n", accumulatedRevenueE))
		message.WriteString(fmt.Sprintf("累計年增率: %.2f%%\n\n", accumulatedYoY))
	}
	message.WriteString("</pre>")

	return message.String()
}

// formatTimeFromTimestamp 將時間戳記格式化為 YYYY/MM 格式
func (s *TgService) formatTimeFromTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006/01")
}

// 格式化股票績效
func (s *TgService) formatPerformanceTable(stockName, symbol string, performanceData *stockDto.StockPerformanceResponseDto) string {
	result := "<pre>"
	// 使用手機友善的格式，避免複雜表格
	result += fmt.Sprintf("📊 <b>%s (%s) 績效表現</b>\n\n", stockName, symbol)

	// 為每個績效期間添加表情符號和格式化
	for _, data := range performanceData.Data {
		// 解析績效數值來決定表情符號
		performanceStr := strings.TrimSuffix(data.Performance, "%")
		performance, err := strconv.ParseFloat(performanceStr, 64)
		var emoji string
		if err == nil {
			if performance >= 0 {
				emoji = "📈" // 上升用📈
			} else {
				emoji = "📉" // 下降用📉
			}
		} else {
			emoji = "📊" // 無法解析用📊
		}

		// 格式化顯示
		result += fmt.Sprintf("%s <b>%s</b>: %s\n", emoji, data.Period, data.Performance)
	}

	result += "</pre>"
	return result
}

// 格式化大盤資訊
func (s *TgService) formatDailyMarketInfoMessage(marketInfo twseDto.DailyMarketInfoResponseDto) string {
	messageText := "<b>台灣股市大盤資訊</b>\n\n"

	// 檢查欄位名稱和資料是否匹配
	if len(marketInfo.Fields) == 0 {
		return messageText + "查無資料"
	}

	for _, row := range marketInfo.Data {
		if len(row) < 6 {
			continue // 跳過資料不完整的行
		}

		// 根據欄位順序解析資料
		// 通常 TWSE 的欄位順序是: ["日期", "成交股數", "成交金額", "成交筆數", "發行量加權股價指數", "漲跌點數"]
		date := row[0]
		volume := row[1]
		amount := row[2]
		transaction := row[3]
		index := row[4]
		change := row[5]

		messageText += fmt.Sprintf("<b>%s</b>\n", date)
		messageText += "<code>"
		messageText += fmt.Sprintf("成交股數：%s\n", volume)
		messageText += fmt.Sprintf("成交金額：%s\n", amount)
		messageText += fmt.Sprintf("成交筆數：%s\n", transaction)
		messageText += fmt.Sprintf("發行量加權股價指數：%s\n", index)
		messageText += fmt.Sprintf("漲跌點數：%s\n", change)
		messageText += "</code>\n"
	}
	return messageText
}

// 格式化股票詳細資訊
func (s *TgService) formatStockInfoMessage(stockInfo *stockDto.StockQuoteInfo) string {
	var message strings.Builder

	message.WriteString("<pre>")
	// 股票基本資訊
	message.WriteString("🏢" + stockInfo.StockName)
	message.WriteString(" (")
	message.WriteString(stockInfo.StockID)
	message.WriteString(")")
	message.WriteString(" | ")
	message.WriteString(stockInfo.Industry)
	message.WriteString(" | ")
	message.WriteString(stockInfo.Market)
	message.WriteString("\n\n")

	// // 價格資訊
	// message.WriteString("現價: ")
	// message.WriteString(fmt.Sprintf("%.2f", stockInfo.CurrentPrice))
	// message.WriteString("\n漲跌: ")
	// message.WriteString(fmt.Sprintf("%+.2f (%.2f%%)", stockInfo.Change, stockInfo.ChangeRate))
	// message.WriteString("\n開盤: ")
	// message.WriteString(fmt.Sprintf("%.2f", stockInfo.OpenPrice))
	// message.WriteString("\n昨收: ")
	// message.WriteString(fmt.Sprintf("%.2f", stockInfo.PrevClose))
	// message.WriteString("\n最高: ")
	// message.WriteString(fmt.Sprintf("%.2f", stockInfo.HighPrice))
	// message.WriteString("\n最低: ")
	// message.WriteString(fmt.Sprintf("%.2f", stockInfo.LowPrice))
	// message.WriteString("\n振幅: ")
	// message.WriteString(fmt.Sprintf("%.2f%%", stockInfo.Amplitude))
	// message.WriteString("\n\n")

	// // 成交資訊
	// volumeStr := fmt.Sprintf("%d", stockInfo.Volume)
	// if stockInfo.Volume >= 1000 {
	// 	volumeStr = fmt.Sprintf("%.1fK", float64(stockInfo.Volume)/1000)
	// }
	// turnoverStr := fmt.Sprintf("%.2f", stockInfo.Turnover/100000000) // 轉換為億元
	// message.WriteString("成交量: ")
	// message.WriteString(volumeStr)
	// message.WriteString(" 張\n成交額: ")
	// message.WriteString(turnoverStr)
	// message.WriteString(" 億\n週轉率: ")
	// message.WriteString(fmt.Sprintf("%.3f%%", stockInfo.VolumeRatio*100))
	// message.WriteString("\n\n")

	// 財務指標
	message.WriteString("💼財務指標:\n")
	message.WriteString("本益比: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.PE))
	message.WriteString("\n本淨比: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.PB))
	message.WriteString("\n市值: ")
	marketCapStr := fmt.Sprintf("%.2f", stockInfo.MarketCap/1000000000000) // 轉換為兆元
	message.WriteString(marketCapStr)
	message.WriteString(" 兆\n每股淨值: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.BookValue))
	message.WriteString("\n近四季EPS: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.EPS))
	message.WriteString("\n營季EPS: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.QuarterEPS))
	message.WriteString("\n年股利: ")
	message.WriteString(fmt.Sprintf("%.2f", stockInfo.Dividend))
	message.WriteString("\n殖利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", stockInfo.DividendRate))
	message.WriteString("\n\n")

	// 獲利能力
	message.WriteString("💡獲利能力:\n")
	message.WriteString("毛利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", stockInfo.GrossMargin))
	message.WriteString("\n營益率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", stockInfo.OperMargin))
	message.WriteString("\n淨利率: ")
	message.WriteString(fmt.Sprintf("%.2f%%", stockInfo.NetMargin))
	message.WriteString("</pre>")
	return message.String()
}
