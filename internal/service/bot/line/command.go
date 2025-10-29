package line

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/tian841224/stock-bot/internal/db/models"
	"github.com/tian841224/stock-bot/internal/infrastructure/imgbb"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/linebot"
	"github.com/tian841224/stock-bot/internal/repository"
	"github.com/tian841224/stock-bot/internal/service/user"
	"github.com/tian841224/stock-bot/pkg/logger"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type LineCommandHandler struct {
	botClient            *linebotInfra.LineBotClient
	lineService          *LineService
	userService          user.UserService
	userSubscriptionRepo repository.UserSubscriptionRepository
	subscriptionItemMap  map[string]models.SubscriptionItem
	imgbbClient          *imgbb.ImgBBClient
}

func NewLineCommandHandler(
	botClient *linebotInfra.LineBotClient,
	lineService *LineService,
	userService user.UserService,
	userSubscriptionRepo repository.UserSubscriptionRepository,
	imgbbClient *imgbb.ImgBBClient,
) *LineCommandHandler {
	return &LineCommandHandler{
		botClient:            botClient,
		lineService:          lineService,
		userService:          userService,
		userSubscriptionRepo: userSubscriptionRepo,
		subscriptionItemMap:  models.SubscriptionItemMap,
		imgbbClient:          imgbbClient,
	}
}

// CommandStart 處理 /start 命令
func (c *LineCommandHandler) CommandStart(userID, replyToken string) error {
	text := `台股機器人指令指南🤖

📊 圖表指令
- /k [股票代碼] - K線圖 (含月均價、最高最低價標示、成交量)
- /p [股票代碼] - 股票績效圖表 (折線圖)
- /r [股票代碼] - 月營收圖表 (柱狀圖+年增率折線)

📈 股票資訊指令
- /d [股票代碼] - 查詢當日收盤資訊 (可指定日期)
- /d [股票代碼] [日期] - 查詢指定日期股價 (格式: YYYY-MM-DD)
- /i [股票代碼] - 查詢公司資訊
- /n [股票代碼] - 查詢股票新聞

📊 市場總覽指令
- /m - 查詢最新大盤資訊 (預設1筆)
- /m [數量] - 查詢指定筆數的大盤資訊
- /t - 查詢當日交易量前20名

🔔 訂閱管理
- /add [股票代碼] - 新增訂閱股票
- /del [股票代碼] - 刪除訂閱股票
- /sub [項目] - 訂閱功能
- /unsub [項目] - 取消訂閱功能
- /list - 查詢已訂閱功能及股票

💡 使用範例：
/k 2330 - 台積電K線圖
/p 0050 - 元大台灣50績效圖表
/r 2330 - 台積電月營收圖表
/d 2330 2025-01-15 - 查詢台積電指定日期股價
/m 3 - 查詢最新3筆大盤資訊`

	return c.sendMessage(replyToken, text)
}

// 處理 /p 命令 - 股票績效圖表 (折線圖)
func (c *LineCommandHandler) CommandPerformanceChart(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	// 取得績效圖表資料
	chartData, caption, err := c.lineService.GetStockPerformanceWithChart(symbol, "line")
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	// 檢查是否有圖表資料
	if len(chartData) == 0 {
		// 如果沒有圖表資料，發送文字版本
		return c.sendMessage(replyToken, caption)
	}

	// 發送圖表
	return c.sendPhoto(userID, replyToken, chartData, caption)
}

// 處理 /d 命令 - 股價詳細資訊（支援日期查詢）
func (c *LineCommandHandler) CommandTodayStockPrice(userID, replyToken, symbol, date string) error {
	// 輸入驗證
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號\n\n使用方式：\n/d 股票代號 - 查詢今日股價\n/d 股票代號 2025-09-01 - 查詢指定日期股價")
	}

	var message string
	var err error

	// 根據是否有日期參數決定呼叫哪個方法
	if date != "" {
		// 驗證日期格式
		if !c.isValidDateFormat(date) {
			return c.sendMessage(replyToken, "日期格式錯誤，請使用 YYYY-MM-DD 格式\n例如：2025-09-01")
		}
		// 查詢指定日期股價
		message, err = c.lineService.GetStockPriceByDate(symbol, date)
	} else {
		message, err = c.lineService.GetStockPriceByDate(symbol, time.Now().Format("2006-01-02"))
	}

	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	// 發送回應
	return c.sendMessage(replyToken, message)
}

// 處理 /k 命令 - 歷史K線圖
func (c *LineCommandHandler) CommandHistoricalCandles(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	chartData, caption, err := c.lineService.GetStockHistoricalCandlesChart(symbol)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendPhoto(userID, replyToken, chartData, caption)
}

// 處理 /n 命令 - 股票新聞
func (c *LineCommandHandler) CommandNews(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	// 取得新聞資料
	newsMessage, err := c.lineService.GetTaiwanStockNews(symbol)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessageWithButtons(userID, replyToken, newsMessage.Text, newsMessage.Buttons)
}

// 處理 /m 命令 - 大盤資訊
func (c *LineCommandHandler) CommandDailyMarketInfo(userID, replyToken string, count int) error {
	// 呼叫業務邏輯
	messageText, err := c.lineService.GetDailyMarketInfo(count)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	// 發送回應
	return c.sendMessage(replyToken, messageText)
}

// 處理 /t 命令 - 交易量前20名
func (c *LineCommandHandler) CommandTopVolumeItems(userID, replyToken string) error {
	// 取得交易量前20名資料
	messageText, err := c.lineService.GetTopVolumeItemsFormatted()
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessage(replyToken, messageText)
}

// 處理 /i 命令 - 股票資訊（可指定日期）
func (c *LineCommandHandler) CommandStockInfo(userID, replyToken, symbol, date string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	// 取得股票資訊
	message, err := c.lineService.GetStockInfo(symbol)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessage(replyToken, message)
}

// 處理 /r 命令 - 股票財報
func (c *LineCommandHandler) CommandRevenue(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代碼")
	}

	chartData, caption, err := c.lineService.GetStockRevenueWithChart(symbol)

	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	// 檢查是否有圖表資料
	if len(chartData) == 0 {
		// 如果沒有圖表資料，發送文字版本
		return c.sendMessage(replyToken, caption)
	}

	// 發送圖表
	return c.sendPhoto(userID, replyToken, chartData, caption)
}

// 處理 /sub 命令 - 訂閱功能
func (c *LineCommandHandler) CommandSubscribe(userID, replyToken, item string) error {
	return c.updateUserSubscription(userID, replyToken, item, "active")
}

// 處理 /unsub 命令 - 取消訂閱功能
func (c *LineCommandHandler) CommandUnsubscribe(userID, replyToken, item string) error {
	return c.updateUserSubscription(userID, replyToken, item, "inactive")
}

// updateUserSubscription 更新使用者訂閱狀態
func (c *LineCommandHandler) updateUserSubscription(userID, replyToken, item, status string) error {
	subscriptionItem, exists := c.subscriptionItemMap[item]
	if !exists {
		return c.sendMessage(replyToken, fmt.Sprintf("無效的訂閱項目: %s", item))
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(replyToken, "無法取得使用者")
	}

	// 檢查是否已經有此訂閱項目
	existingSubscription, err := c.userSubscriptionRepo.GetUserSubscriptionByItem(user.ID, subscriptionItem)
	if err != nil {
		// 如果沒有找到訂閱項目，且是要訂閱，則新增
		if status == "active" {
			if err := c.userSubscriptionRepo.AddUserSubscriptionItem(user.ID, subscriptionItem); err != nil {
				logger.Log.Error("新增訂閱項目失敗", zap.Error(err))
				return c.sendMessage(replyToken, "訂閱失敗，請稍後再試")
			}
			return c.sendMessage(replyToken, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
		} else {
			return c.sendMessage(replyToken, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 如果狀態相同，不需要更新
	if existingSubscription.Status == status {
		if status == "active" {
			return c.sendMessage(replyToken, fmt.Sprintf("已訂閱：%s", subscriptionItem.GetName()))
		} else {
			return c.sendMessage(replyToken, fmt.Sprintf("未訂閱此項目：%s", subscriptionItem.GetName()))
		}
	}

	// 更新訂閱狀態
	if err := c.userSubscriptionRepo.UpdateUserSubscriptionItem(user.ID, subscriptionItem, status); err != nil {
		logger.Log.Error("更新訂閱狀態失敗", zap.Error(err))
		return c.sendMessage(replyToken, "操作失敗，請稍後再試")
	}

	if status == "active" {
		return c.sendMessage(replyToken, fmt.Sprintf("訂閱成功：%s", subscriptionItem.GetName()))
	} else {
		return c.sendMessage(replyToken, fmt.Sprintf("取消訂閱成功：%s", subscriptionItem.GetName()))
	}
}

// 處理 /add 命令 - 新增股票訂閱
func (c *LineCommandHandler) CommandAddStock(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(replyToken, "無法取得使用者")
	}

	// 新增股票訂閱
	message, err := c.lineService.AddUserStockSubscription(user.ID, symbol)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessage(replyToken, message)
}

// 處理 /del 命令 - 刪除股票訂閱
func (c *LineCommandHandler) CommandDeleteStock(userID, replyToken, symbol string) error {
	if symbol == "" {
		return c.sendMessage(replyToken, "請輸入股票代號")
	}

	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(replyToken, "無法取得使用者")
	}

	// 刪除股票訂閱
	message, err := c.lineService.DeleteUserStockSubscription(user.ID, symbol)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessage(replyToken, message)
}

// 處理 /list 命令 - 列出訂閱項目
func (c *LineCommandHandler) CommandListSubscriptions(userID, replyToken string) error {
	// 取得使用者資料
	user, err := c.userService.GetUserByAccountID(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("取得使用者失敗", zap.Error(err))
		return c.sendMessage(replyToken, "無法取得使用者")
	}

	// 取得訂閱清單
	messageText, err := c.lineService.GetUserSubscriptionList(user.ID)
	if err != nil {
		return c.sendMessage(replyToken, err.Error())
	}

	return c.sendMessage(replyToken, messageText)
}

// 輔助方法

// 發送訊息
func (c *LineCommandHandler) sendMessage(replyToken, text string) error {
	err := c.botClient.ReplyMessage(replyToken, text)
	if err != nil {
		logger.Log.Error("發送訊息失敗", zap.Error(err))
	}
	return err
}

// 發送帶有按鈕的訊息
func (c *LineCommandHandler) sendMessageWithButtons(userID, replyToken, text string, buttons []linebot.TemplateAction) error {
	if len(buttons) == 0 {
		return c.sendMessage(replyToken, text)
	}

	// 建立按鈕模板
	template := linebot.NewButtonsTemplate(
		"", "", text, buttons...,
	)

	_, err := c.botClient.Client.ReplyMessage(replyToken, linebot.NewTemplateMessage("按鈕", template)).Do()
	if err != nil {
		logger.Log.Error("發送帶有按鈕的訊息失敗", zap.Error(err))
	}
	return err
}

// 發送圖片
func (c *LineCommandHandler) sendPhoto(userID, replyToken string, data []byte, caption string) error {
	// 如果沒有 ImgBB 客戶端，只發送文字訊息
	if c.imgbbClient == nil {
		logger.Log.Warn("ImgBB 客戶端未設定，只發送文字訊息")
		return c.sendMessage(replyToken, caption)
	}

	// 上傳圖片到 ImgBB
	options := &imgbb.UploadOptions{
		Name: "stock_chart", // 預設檔案名稱
	}

	// 將 byte slice 轉換為 Reader
	reader := bytes.NewReader(data)
	resp, err := c.imgbbClient.UploadFromFile(reader, "chart.png", options)
	if err != nil {
		logger.Log.Error("上傳圖片到 ImgBB 失敗", zap.Error(err))
		// 如果上傳失敗，只發送文字訊息
		return c.sendMessage(replyToken, caption)
	}

	// 建立圖片訊息
	imageMessage := linebot.NewImageMessage(resp.Data.URL, resp.Data.URL)

	// 發送圖片
	_, err = c.botClient.Client.ReplyMessage(replyToken, imageMessage).Do()
	if err != nil {
		logger.Log.Error("發送圖片訊息失敗", zap.Error(err))
	}

	return err
}

// 驗證日期格式是否為 YYYY-MM-DD
func (c *LineCommandHandler) isValidDateFormat(date string) bool {
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
