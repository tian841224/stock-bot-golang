package tg

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"stock-bot/internal/db/models"
	"stock-bot/internal/repository"
	"stock-bot/internal/service/user"
	"stock-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgCommandHandler struct {
	botClient            *tgbotapi.BotAPI
	tgService            *TgService
	userService          user.UserService
	userSubscriptionRepo repository.UserSubscriptionRepository
	subscriptionItemMap  map[string]models.SubscriptionItem
}

func NewTgCommandHandler(
	botClient *tgbotapi.BotAPI,
	tgService *TgService,
	userService user.UserService,
	userSubscriptionRepo repository.UserSubscriptionRepository,
) *TgCommandHandler {
	return &TgCommandHandler{
		botClient:            botClient,
		tgService:            tgService,
		userService:          userService,
		userSubscriptionRepo: userSubscriptionRepo,
		subscriptionItemMap:  models.SubscriptionItemMap,
	}
}

// commandStart 處理 /start 命令
func (c *TgCommandHandler) CommandStart(userID int64) error {
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
- /d [股票代碼] - 查詢今日股價詳細資訊
- /d [股票代碼] [日期] - 查詢指定日期股價 (格式: YYYY-MM-DD)
- /p [股票代碼] - 查詢股票績效
- /pc [股票代碼] - 查詢股票績效圖表 (折線圖)
- /pb [股票代碼] - 查詢股票績效圖表 (柱狀圖)
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

	return c.sendMessage(userID, text)
}

// CommandKline 處理 /k 命令 - K線圖
func (c *TgCommandHandler) CommandKline(userID int64, symbol, timeRange string) error {
	// 呼叫業務邏輯
	imageData, caption, timeRangeText, err := c.tgService.GetStockKlineImage(symbol, timeRange)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 發送圖片
	photo := tgbotapi.NewPhoto(userID, tgbotapi.FileBytes{
		Name:  "kline.png",
		Bytes: imageData,
	})
	photo.Caption = caption

	_, err = c.botClient.Send(photo)
	if err != nil {
		logger.Log.Error("發送圖片失敗", zap.Error(err))
		return c.sendMessage(userID, fmt.Sprintf("%s %s K線圖", caption, timeRangeText))
	}

	return nil
}

// CommandPerformance 處理 /p 命令 - 股票績效
func (c *TgCommandHandler) CommandPerformance(userID int64, symbol string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得績效資料
	performanceText, err := c.tgService.GetStockPerformance(symbol)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 發送HTML格式的表格訊息
	return c.sendMessageHTML(userID, performanceText)
}

// CommandPerformanceChart 處理 /pc 命令 - 股票績效圖表 (折線圖)
func (c *TgCommandHandler) CommandPerformanceChart(userID int64, symbol string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得績效圖表資料
	chartData, caption, err := c.tgService.GetStockPerformanceWithChart(symbol, "line")
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 檢查是否有圖表資料
	if len(chartData) == 0 {
		// 如果沒有圖表資料，發送文字版本
		return c.sendMessageHTML(userID, caption)
	}

	// 發送圖表
	return c.sendPhoto(userID, chartData, caption)
}

// CommandPerformanceBarChart 處理 /pb 命令 - 股票績效圖表 (柱狀圖)
func (c *TgCommandHandler) CommandPerformanceBarChart(userID int64, symbol string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得績效圖表資料
	chartData, caption, err := c.tgService.GetStockPerformanceWithChart(symbol, "bar")
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 檢查是否有圖表資料
	if len(chartData) == 0 {
		// 如果沒有圖表資料，發送文字版本
		return c.sendMessageHTML(userID, caption)
	}

	// 發送圖表
	return c.sendPhoto(userID, chartData, caption)
}

// CommandTodayStockPrice 處理 /d 命令 - 股價詳細資訊（支援日期查詢）
func (c *TgCommandHandler) CommandTodayStockPrice(userID int64, symbol, date string) error {
	// 輸入驗證
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號\n\n使用方式：\n/d 股票代號 - 查詢今日股價\n/d 股票代號 2025-09-01 - 查詢指定日期股價")
	}

	var message string
	var err error

	// 根據是否有日期參數決定呼叫哪個方法
	if date != "" {
		// 驗證日期格式
		if !c.isValidDateFormat(date) {
			return c.sendMessage(userID, "日期格式錯誤，請使用 YYYY-MM-DD 格式\n例如：2025-09-01")
		}
		// 查詢指定日期股價
		message, err = c.tgService.GetStockPriceByDate(symbol, date)
	} else {
		message, err = c.tgService.GetStockPriceByDate(symbol, time.Now().Format("2006-01-02"))
	}

	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 發送回應
	return c.sendMessageHTML(userID, message)
}

// CommandNews 處理 /n 命令 - 股票新聞
// func (c *TgCommandHandler) CommandNews(userID int64, symbol string) error {
// 	if symbol == "" {
// 		return c.sendMessage(userID, "請輸入股票代號")
// 	}

// 	// 取得新聞資料
// 	message, err := c.tgService.GetStockNews(symbol)
// 	if err != nil {
// 		return c.sendMessage(userID, err.Error())
// 	}

// 	return c.sendMessage(userID, message)
// }

// CommandDailyMarketInfo 處理 /m 命令 - 大盤資訊
func (c *TgCommandHandler) CommandDailyMarketInfo(userID int64, count int) error {
	// 呼叫業務邏輯
	messageText, err := c.tgService.GetDailyMarketInfo(count)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	// 發送回應
	return c.sendMessageHTML(userID, messageText)
}

// CommandTopVolumeItems 處理 /t 命令 - 交易量前20名
func (c *TgCommandHandler) CommandTopVolumeItems(userID int64) error {
	// 取得交易量前20名資料
	messageText, err := c.tgService.GetTopVolumeItemsFormatted()
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	return c.sendMessageHTML(userID, messageText)
}

// CommandStockInfo 處理 /i 命令 - 股票資訊（可指定日期）
func (c *TgCommandHandler) CommandStockInfo(userID int64, symbol, date string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得股票資訊
	message, err := c.tgService.GetStockInfo(symbol)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	return c.sendMessageHTML(userID, message)
}

// CommandSubscribe 處理 /sub 命令 - 訂閱功能
func (c *TgCommandHandler) CommandSubscribe(userID int64, item string) error {
	return c.updateUserSubscription(userID, item, "active")
}

// CommandUnsubscribe 處理 /unsub 命令 - 取消訂閱功能
func (c *TgCommandHandler) CommandUnsubscribe(userID int64, item string) error {
	return c.updateUserSubscription(userID, item, "inactive")
}

// updateUserSubscription 更新使用者訂閱狀態
func (c *TgCommandHandler) updateUserSubscription(userID int64, item, status string) error {
	subscriptionItem, exists := c.subscriptionItemMap[item]
	if !exists {
		return c.sendMessage(userID, fmt.Sprintf("無效的訂閱項目: %s", item))
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(userID, "無法取得使用者")
	}

	// 檢查是否已經有此訂閱項目
	existingSubscription, err := c.userSubscriptionRepo.GetUserSubscriptionByItem(user.ID, subscriptionItem)
	if err != nil {
		// 如果沒有找到訂閱項目，且是要訂閱，則新增
		if status == "active" {
			if err := c.userSubscriptionRepo.AddUserSubscriptionItem(user.ID, subscriptionItem); err != nil {
				logger.Log.Error("新增訂閱項目失敗", zap.Error(err))
				return c.sendMessage(userID, "訂閱失敗，請稍後再試")
			}
			return c.sendMessage(userID, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
		} else {
			return c.sendMessage(userID, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 如果狀態相同，不需要更新
	if existingSubscription.Status == status {
		if status == "active" {
			return c.sendMessage(userID, fmt.Sprintf("已訂閱：%s", subscriptionItem.GetName()))
		} else {
			return c.sendMessage(userID, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 更新訂閱狀態
	if err := c.userSubscriptionRepo.UpdateUserSubscriptionItem(user.ID, subscriptionItem, status); err != nil {
		logger.Log.Error("更新訂閱狀態失敗", zap.Error(err))
		return c.sendMessage(userID, "操作失敗，請稍後再試")
	}

	if status == "active" {
		return c.sendMessage(userID, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
	} else {
		return c.sendMessage(userID, fmt.Sprintf("取消訂閱成功：%s", subscriptionItem.GetName()))
	}
}

// CommandAddStock 處理 /add 命令 - 新增股票訂閱
func (c *TgCommandHandler) CommandAddStock(userID int64, symbol string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(userID, "無法取得使用者")
	}

	// 新增股票訂閱
	message, err := c.tgService.AddUserStockSubscription(user.ID, symbol)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	return c.sendMessage(userID, message)
}

// CommandDeleteStock 處理 /del 命令 - 刪除股票訂閱
func (c *TgCommandHandler) CommandDeleteStock(userID int64, symbol string) error {
	if symbol == "" {
		return c.sendMessage(userID, "請輸入股票代號")
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(userID, "無法取得使用者")
	}

	// 刪除股票訂閱
	message, err := c.tgService.DeleteUserStockSubscription(user.ID, symbol)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	return c.sendMessage(userID, message)
}

// CommandListSubscriptions 處理 /list 命令 - 列出訂閱項目
func (c *TgCommandHandler) CommandListSubscriptions(userID int64) error {
	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(userID, "無法取得使用者")
	}

	// 取得訂閱清單
	messageText, err := c.tgService.GetUserSubscriptionList(user.ID)
	if err != nil {
		return c.sendMessage(userID, err.Error())
	}

	return c.sendMessageHTML(userID, messageText)
}

// 輔助方法

func (c *TgCommandHandler) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.botClient.Send(msg)
	if err != nil {
		logger.Log.Error("發送訊息失敗", zap.Error(err))
	}
	return err
}

func (c *TgCommandHandler) sendMessageHTML(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.botClient.Send(msg)
	if err != nil {
		logger.Log.Error("發送 HTML 訊息失敗", zap.Error(err))
	}
	return err
}

func (c *TgCommandHandler) sendDocument(chatID int64, data []byte, filename, caption string) error {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  filename,
		Bytes: data,
	})
	doc.Caption = caption
	doc.ParseMode = tgbotapi.ModeHTML
	_, err := c.botClient.Send(doc)
	if err != nil {
		logger.Log.Error("發送文件失敗", zap.Error(err))
	}
	return err
}

func (c *TgCommandHandler) sendPhoto(chatID int64, data []byte, caption string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  "chart.png",
		Bytes: data,
	})
	photo.Caption = caption
	photo.ParseMode = tgbotapi.ModeHTML
	_, err := c.botClient.Send(photo)
	if err != nil {
		logger.Log.Error("發送圖片失敗", zap.Error(err))
	}
	return err
}

// isValidDateFormat 驗證日期格式是否為 YYYY-MM-DD
func (c *TgCommandHandler) isValidDateFormat(date string) bool {
	// 檢查長度
	if len(date) != 10 {
		return false
	}

	// 使用正則表達式驗證格式
	matched, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	if err != nil || !matched {
		return false
	}

	// 嘗試解析日期以確保是有效日期
	_, err = time.Parse("2006-01-02", date)
	return err == nil
}
