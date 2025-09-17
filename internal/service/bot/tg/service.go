package tg

import (
	"fmt"
	"time"

	"stock-bot/internal/repository"
	"stock-bot/internal/service/twstock"
	stockDto "stock-bot/internal/service/twstock/dto"
	"stock-bot/pkg/logger"

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

// GetTodayStockPrice 取得今日股價詳細資訊
func (s *TgService) GetTodayStockPrice(symbol string) (string, error) {
	if symbol == "" {
		return "", fmt.Errorf("股票代號不能為空")
	}

	// 取得今日股價資訊
	stockInfo, err := s.stockService.GetStockPrice(symbol)
	if err != nil {
		logger.Log.Error("取得今日股價失敗", zap.Error(err))
		return "", fmt.Errorf("查無此股票資料，請重新確認")
	}

	// 建立今日股價詳細訊息
	emoji := ""
	if stockInfo.UpDownSign == "+" {
		emoji = "📈"
	} else if stockInfo.UpDownSign == "-" {
		emoji = "📉"
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
		stockInfo.Date,
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

// GetStockKlineImage 取得股票 K 線圖
func (s *TgService) GetStockKlineImage(symbol, timeRange string) ([]byte, string, string, error) {
	if symbol == "" {
		return nil, "", "", fmt.Errorf("請輸入股票代號")
	}

	// 驗證股票代號
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return nil, "", "", fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 轉換時間範圍
	timeRangeText := s.convertTimeRange(timeRange)

	// 取得 K 線圖
	imageData, _, err := s.stockService.GetStockAnalysis(symbol)
	if err != nil {
		logger.Log.Error("取得股票分析圖表失敗", zap.Error(err))
		return nil, "", "", fmt.Errorf("取得 K 線圖失敗，請稍後再試")
	}

	caption := fmt.Sprintf("%s(%s) K線圖　💹", stockName, symbol)
	return imageData, caption, timeRangeText, nil
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

// GetStockNews 取得股票新聞
func (s *TgService) GetStockNews(symbol string) (string, error) {
	// 驗證股票代號
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return "", fmt.Errorf("查無此股票代號，請重新確認")
	}

	// 這裡需要實際的新聞服務，暫時返回模擬資料
	message := fmt.Sprintf("⚡️%s(%s)-即時新聞\n\n暫無新聞資料，功能開發中...", stockName, symbol)
	return message, nil
}

// GetYahooStockNews 取得 Yahoo 股票新聞
func (s *TgService) GetYahooStockNews(symbol string) (string, error) {
	// 這裡需要實際的 Yahoo 新聞服務，暫時返回模擬資料
	message := fmt.Sprintf("⚡️%s-即時新聞\n\n暫無新聞資料，功能開發中...", symbol)
	return message, nil
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
		if item.UpDownSign == "+" {
			emoji = "📈"
		} else if item.UpDownSign == "-" {
			emoji = "📉"
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
	if stockInfo.UpDownSign == "+" {
		emoji = "📈"
	} else if stockInfo.UpDownSign == "-" {
		emoji = "📉"
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

// convertTimeRange 轉換時間範圍顯示文字
func (s *TgService) convertTimeRange(timeRange string) string {
	switch timeRange {
	case "h":
		return "分時"
	case "d":
		return "日K"
	case "w":
		return "週K"
	case "m":
		return "月K"
	case "5m":
		return "5分"
	case "15m":
		return "15分"
	case "30m":
		return "30分"
	case "60m":
		return "60分"
	default:
		return "日K" // 預設值
	}
}

// GetDailyMarketInfoFormatted 取得格式化的大盤資訊
func (s *TgService) GetDailyMarketInfoFormatted(count int) (string, error) {
	marketInfoList, err := s.stockService.GetDailyMarketInfo(count)
	if err != nil {
		logger.Log.Error("取得大盤資訊失敗", zap.Error(err))
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	if len(marketInfoList) == 0 {
		return "", fmt.Errorf("查無資料，請確認後再試")
	}

	messageText := "<b>台灣股市大盤資訊</b>\n\n"
	for _, row := range marketInfoList {
		messageText += fmt.Sprintf("<b>%s</b>\n", row.Date)
		messageText += "<code>"
		messageText += fmt.Sprintf("成交股數：%s\n", row.Volume)
		messageText += fmt.Sprintf("成交金額：%s\n", row.Amount)
		messageText += fmt.Sprintf("成交筆數：%s\n", row.Transaction)
		messageText += fmt.Sprintf("發行量加權股價指數：%s\n", row.Index)
		messageText += fmt.Sprintf("漲跌點數：%s\n", row.Change)
		messageText += "</code>\n"
	}

	return messageText, nil
}

// formatPerformanceTable 格式化股票績效為HTML表格
func (s *TgService) formatPerformanceTable(stockName, symbol string, performanceData *stockDto.StockPerformanceResponseDto) string {

	// 使用 <pre> 標籤來保持格式對齊，並加上邊框效果
	result := "<pre>"
	result += fmt.Sprintf("<b>%s(%s) 績效表現 ✨</b>", stockName, symbol)
	result += "┌─────────┬─────────────┐\n"
	result += "│ Period  │ Performance │\n"
	result += "├─────────┼─────────────┤\n"

	// 加入每行資料
	for _, data := range performanceData.Data {
		// 確保中文字元對齊，使用固定寬度格式
		periodFormatted := fmt.Sprintf("%-7s", data.Period)
		performanceFormatted := fmt.Sprintf("%-11s", data.Performance)
		result += fmt.Sprintf("│ %s │ %s │\n", periodFormatted, performanceFormatted)
	}

	result += "└─────────┴─────────────┘"
	result += "</pre>"

	return result
}
