package tg

import (
	"fmt"
	"strconv"
	"time"

	"stock-bot/internal/db/models"
	"stock-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleStart 處理 /start 命令
func (s *TgService) handleStart(userID int64) error {
	text := `台股機器人指令指南🤖

📊 基本K線圖
格式：/k [股票代碼] [時間範圍]

時間範圍選項（預設：d）：
- h - 時K線
- d - 日K線
- w - 週K線
- m - 月K線
- 5m - 5分K線
- 15m - 15分K線
- 30m - 30分K線
- 60m - 60分K線

股票資訊指令
- /d [股票代碼] - 查詢股票詳細資訊
- /p [股票代碼] - 查詢股票績效
- /n [股票代碼] - 查詢股票新聞
- /yn [股票代碼] - 查詢Yahoo股票新聞（預設：台股新聞）
- /i [股票代碼] - 查詢當日收盤資訊 (可指定日期 ex: /i 2330 20250101)

市場總覽指令
- /m - 查詢大盤資訊
- /t - 查詢當日交易量前20名

訂閱股票資訊
- /add [股票代碼] - 訂閱 股票
- /del [股票代碼] - 取消訂閱 股票
- /sub 1 - 訂閱 當日個股資訊
- /sub 2 - 訂閱 觀察清單新聞
- /sub 3 - 訂閱 當日市場成交行情
- /sub 4 - 訂閱 當日交易量前20名

查詢指令
- /list - 查詢已訂閱功能及股票

(取消訂閱 unsub + 代號)`

	return s.sendMessage(userID, text)
}

// handleKline 處理 /k 命令 - K線圖
func (s *TgService) handleKline(userID int64, symbol, timeRange string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 驗證股票代號
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return s.sendMessage(userID, "查無此股票代號，請重新確認")
	}

	// 轉換時間範圍
	timeRangeText := s.convertTimeRange(timeRange)

	// 取得 K 線圖（這裡需要實際的圖表服務，暫時返回文字訊息）
	imageData, _, err := s.stockService.GetStockAnalysis(symbol)
	if err != nil {
		logger.Log.Error("取得股票分析圖表失敗", zap.Error(err))
		return s.sendMessage(userID, "取得 K 線圖失敗，請稍後再試")
	}

	// 發送圖片
	photo := tgbotapi.NewPhoto(userID, tgbotapi.FileBytes{
		Name:  "kline.png",
		Bytes: imageData,
	})
	photo.Caption = fmt.Sprintf("%s(%s) K線圖　💹", stockName, symbol)

	_, err = s.botClient.Send(photo)
	if err != nil {
		logger.Log.Error("發送圖片失敗", zap.Error(err))
		return s.sendMessage(userID, fmt.Sprintf("%s(%s) %s K線圖", stockName, symbol, timeRangeText))
	}

	return nil
}

// handlePerformance 處理 /p 命令 - 股票績效
func (s *TgService) handlePerformance(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 驗證股票代號並取得基本資訊
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return s.sendMessage(userID, "查無此股票代號，請重新確認")
	}

	// 取得績效圖表
	imageData, _, err := s.stockService.GetStockAnalysis(symbol)
	if err != nil {
		logger.Log.Error("取得股票績效圖表失敗", zap.Error(err))
		return s.sendMessage(userID, "取得績效資料失敗，請稍後再試")
	}

	// 發送圖片
	photo := tgbotapi.NewPhoto(userID, tgbotapi.FileBytes{
		Name:  "performance.png",
		Bytes: imageData,
	})
	photo.Caption = fmt.Sprintf("%s(%s) 績效表現　✨", stockName, symbol)

	_, err = s.botClient.Send(photo)
	if err != nil {
		logger.Log.Error("發送圖片失敗", zap.Error(err))
		return s.sendMessage(userID, fmt.Sprintf("%s(%s) 績效表現", stockName, symbol))
	}

	return nil
}

// handleDetailPrice 處理 /d 命令 - 股票詳細價格資訊
func (s *TgService) handleDetailPrice(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 取得股票價格資訊
	stockInfo, err := s.stockService.GetStockPrice(symbol)
	if err != nil {
		logger.Log.Error("取得股票價格失敗", zap.Error(err))
		return s.sendMessage(userID, "查無此股票資料，請重新確認")
	}

	// 建立詳細資訊訊息
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
成交股數：%d
成交筆數：%d</code>`,
		stockInfo.Date,
		stockInfo.StockName, stockInfo.StockID, emoji,
		stockInfo.OpenPrice,
		stockInfo.ClosePrice,
		stockInfo.UpDownSign, stockInfo.ChangeAmount, stockInfo.PercentageChange,
		stockInfo.HighPrice,
		stockInfo.LowPrice,
		stockInfo.Volume,
		stockInfo.Transaction)

	return s.sendMessageHTML(userID, message)
}

// handleNews 處理 /n 命令 - 股票新聞
func (s *TgService) handleNews(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 驗證股票代號
	valid, stockName, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return s.sendMessage(userID, "查無此股票代號，請重新確認")
	}

	// 這裡需要實際的新聞服務，暫時返回模擬資料
	message := fmt.Sprintf("⚡️%s(%s)-即時新聞\n\n暫無新聞資料，功能開發中...", stockName, symbol)
	return s.sendMessage(userID, message)
}

// handleYahooNews 處理 /yn 命令 - Yahoo 股票新聞
func (s *TgService) handleYahooNews(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 這裡需要實際的 Yahoo 新聞服務，暫時返回模擬資料
	message := fmt.Sprintf("⚡️%s-即時新聞\n\n暫無新聞資料，功能開發中...", symbol)
	return s.sendMessage(userID, message)
}

// handleDailyMarketInfo 處理 /m 命令 - 大盤資訊
func (s *TgService) handleDailyMarketInfo(userID int64, count int) error {
	marketInfoList, err := s.stockService.GetDailyMarketInfo(count)
	if err != nil {
		logger.Log.Error("取得大盤資訊失敗", zap.Error(err))
		return s.sendMessage(userID, "查無資料，請確認後再試")
	}

	if len(marketInfoList) == 0 {
		return s.sendMessage(userID, "查無資料，請確認後再試")
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

	return s.sendMessageHTML(userID, messageText)
}

// handleTopVolumeItems 處理 /t 命令 - 交易量前20名
func (s *TgService) handleTopVolumeItems(userID int64) error {
	topItems, err := s.stockService.GetTopVolumeItems()
	if err != nil {
		logger.Log.Error("取得交易量前20名失敗", zap.Error(err))
		return s.sendMessage(userID, "查無資料，請確認後再試")
	}

	if len(topItems) == 0 {
		return s.sendMessage(userID, "查無資料，請確認後再試")
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
		messageText += fmt.Sprintf("成交股數：%d\n", item.Volume)
		messageText += fmt.Sprintf("成交筆數：%d\n", item.Transaction)
		messageText += fmt.Sprintf("開盤價：%.2f\n", item.OpenPrice)
		messageText += fmt.Sprintf("收盤價：%.2f\n", item.ClosePrice)
		messageText += fmt.Sprintf("漲跌幅：%s%.2f (%s)\n", item.UpDownSign, item.ChangeAmount, item.PercentageChange)
		messageText += fmt.Sprintf("最高價：%.2f\n", item.HighPrice)
		messageText += fmt.Sprintf("最低價：%.2f\n", item.LowPrice)
		messageText += "</code>\n"
	}

	return s.sendMessageHTML(userID, messageText)
}

// handleStockInfo 處理 /i 命令 - 股票資訊（可指定日期）
func (s *TgService) handleStockInfo(userID int64, symbol, date string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 取得股票價格資訊
	stockInfo, err := s.stockService.GetStockPrice(symbol, date)
	if err != nil {
		logger.Log.Error("取得股票資訊失敗", zap.Error(err))
		return s.sendMessage(userID, "查無資料，請確認後再試")
	}

	// 格式化日期顯示
	displayDate := stockInfo.Date
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
成交股數：%d
成交筆數：%d</code>`,
		displayDate,
		stockInfo.StockName, stockInfo.StockID, emoji,
		stockInfo.OpenPrice,
		stockInfo.ClosePrice,
		stockInfo.UpDownSign, stockInfo.ChangeAmount, stockInfo.PercentageChange,
		stockInfo.HighPrice,
		stockInfo.LowPrice,
		stockInfo.Volume,
		stockInfo.Transaction)

	return s.sendMessageHTML(userID, message)
}

// handleSubscribe 處理 /sub 命令 - 訂閱功能
func (s *TgService) handleSubscribe(userID int64, item string) error {
	return s.updateUserSubscription(userID, item, "active")
}

// handleUnsubscribe 處理 /unsub 命令 - 取消訂閱功能
func (s *TgService) handleUnsubscribe(userID int64, item string) error {
	return s.updateUserSubscription(userID, item, "inactive")
}

// updateUserSubscription 更新使用者訂閱狀態
func (s *TgService) updateUserSubscription(userID int64, item, status string) error {
	subscriptionItem, exists := s.subscriptionItemMap[item]
	if !exists {
		return s.sendMessage(userID, fmt.Sprintf("無效的訂閱項目: %s", item))
	}

	// 取得使用者資料
	user, err := s.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return s.sendMessage(userID, "無法取得使用者")
	}

	// 檢查是否已經有此訂閱項目
	existingSubscription, err := s.userSubscriptionRepo.GetUserSubscriptionByItem(user.ID, subscriptionItem)
	if err != nil {
		// 如果沒有找到訂閱項目，且是要訂閱，則新增
		if status == "active" {
			if err := s.userSubscriptionRepo.AddUserSubscriptionItem(user.ID, subscriptionItem); err != nil {
				logger.Log.Error("新增訂閱項目失敗", zap.Error(err))
				return s.sendMessage(userID, "訂閱失敗，請稍後再試")
			}
			return s.sendMessage(userID, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
		} else {
			return s.sendMessage(userID, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 如果狀態相同，不需要更新
	if existingSubscription.Status == status {
		if status == "active" {
			return s.sendMessage(userID, fmt.Sprintf("已訂閱：%s", subscriptionItem.GetName()))
		} else {
			return s.sendMessage(userID, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 更新訂閱狀態
	if err := s.userSubscriptionRepo.UpdateUserSubscriptionItem(user.ID, subscriptionItem, status); err != nil {
		logger.Log.Error("更新訂閱狀態失敗", zap.Error(err))
		return s.sendMessage(userID, "操作失敗，請稍後再試")
	}

	if status == "active" {
		return s.sendMessage(userID, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
	} else {
		return s.sendMessage(userID, fmt.Sprintf("取消訂閱成功：%s", subscriptionItem.GetName()))
	}
}

// handleAddStock 處理 /add 命令 - 新增股票訂閱
func (s *TgService) handleAddStock(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 驗證股票代號
	valid, _, err := s.stockService.ValidateStockID(symbol)
	if err != nil || !valid {
		return s.sendMessage(userID, "無此股票代號，請重新確認")
	}

	// 取得使用者資料
	user, err := s.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return s.sendMessage(userID, "無法取得使用者")
	}

	// 新增股票訂閱
	success, err := s.userSubscriptionRepo.AddUserSubscriptionStock(user.ID, symbol)
	if err != nil {
		logger.Log.Error("新增股票訂閱失敗", zap.Error(err))
		return s.sendMessage(userID, "訂閱失敗，請稍後再試")
	}

	if !success {
		return s.sendMessage(userID, "已訂閱過此股票")
	}

	return s.sendMessage(userID, "訂閱成功")
}

// handleDeleteStock 處理 /del 命令 - 刪除股票訂閱
func (s *TgService) handleDeleteStock(userID int64, symbol string) error {
	if symbol == "" {
		return s.sendMessage(userID, "請輸入股票代號")
	}

	// 取得使用者資料
	user, err := s.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return s.sendMessage(userID, "無法取得使用者")
	}

	// 刪除股票訂閱
	success, err := s.userSubscriptionRepo.DeleteUserSubscriptionStock(user.ID, symbol)
	if err != nil {
		logger.Log.Error("刪除股票訂閱失敗", zap.Error(err))
		return s.sendMessage(userID, "取消訂閱失敗，請稍後再試")
	}

	if !success {
		return s.sendMessage(userID, "取消訂閱失敗，請檢查是否已訂閱")
	}

	return s.sendMessage(userID, "取消訂閱成功")
}

// handleListSubscriptions 處理 /list 命令 - 列出訂閱項目
func (s *TgService) handleListSubscriptions(userID int64) error {
	// 取得使用者資料
	user, err := s.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return s.sendMessage(userID, "無法取得使用者")
	}

	// 取得使用者訂閱項目
	subscriptions, err := s.userSubscriptionRepo.GetUserSubscriptionList(user.ID)
	if err != nil {
		logger.Log.Error("取得使用者訂閱項目失敗", zap.Error(err))
		return s.sendMessage(userID, "取得訂閱清單失敗")
	}

	// 取得使用者訂閱股票
	subscriptionStocks, err := s.userSubscriptionRepo.GetUserSubscriptionStockList(user.ID)
	if err != nil {
		logger.Log.Error("取得使用者訂閱股票失敗", zap.Error(err))
		return s.sendMessage(userID, "取得訂閱清單失敗")
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

	return s.sendMessageHTML(userID, messageText)
}
