package tg

import (
	"strconv"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/db/models"
	"github.com/tian841224/stock-bot/internal/service/user"
	"github.com/tian841224/stock-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgHandler struct {
	commandHandler *TgCommandHandler
	userService    user.UserService
}

func NewTgHandler(commandHandler *TgCommandHandler, userService user.UserService) *TgHandler {
	return &TgHandler{
		commandHandler: commandHandler,
		userService:    userService,
	}
}

func (s *TgHandler) ProcessUpdate(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	return s.processCommand(update.Message)
}

func (s *TgHandler) processCommand(message *tgbotapi.Message) error {
	if message.Text == "" {
		return nil
	}

	userID := message.Chat.ID
	messageText := message.Text

	logger.Log.Info("收到 Telegram 訊息",
		zap.Int64("user_id", userID),
		zap.String("message", messageText))

	// 確保使用者存在
	_, err := s.userService.GetOrCreate(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("建立或取得使用者失敗", zap.Error(err))
		return s.commandHandler.sendMessage(userID, "系統錯誤，請稍後再試")
	}

	parts := strings.Fields(messageText)
	if len(parts) == 0 {
		return nil
	}

	command := parts[0]
	var arg1, arg2 string
	if len(parts) > 1 {
		arg1 = parts[1]
	}
	if len(parts) > 2 {
		arg2 = parts[2]
	}

	switch command {
	case "/start":
		return s.commandHandler.CommandStart(userID)
	case "/k":
		// 歷史K線圖
		return s.commandHandler.CommandHistoricalCandles(userID, arg1)
	case "/p":
		// 績效圖表
		return s.commandHandler.CommandPerformanceChart(userID, arg1)
	case "/d":
		// 今日股價
		if arg2 == "" {
			// 取得台灣時區當前時間
			taipeiLocation, _ := time.LoadLocation("Asia/Taipei")
			now := time.Now().In(taipeiLocation)

			// 判斷是否在早上9點半之前
			marketOpenTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, taipeiLocation)

			if now.Before(marketOpenTime) {
				// 9點半前，使用前一天日期
				arg2 = now.AddDate(0, 0, -1).Format("2006-01-02")
			} else {
				// 9點半後，使用當天日期
				arg2 = now.Format("2006-01-02")
			}
		}
		return s.commandHandler.CommandTodayStockPrice(userID, arg1, arg2)
	case "/t":
		// 交易量前20名
		return s.commandHandler.CommandTopVolumeItems(userID)
	case "/i":
		// 股票資訊
		return s.commandHandler.CommandStockInfo(userID, arg1, arg2)
	case "/r":
		// 財報
		return s.commandHandler.CommandRevenue(userID, arg1)
	case "/n":
		// 新聞
		return s.commandHandler.CommandNews(userID, arg1)
	case "/m":
		// 大盤資訊
		count := 1 // 預設顯示5筆
		if arg1 != "" {
			if parsedCount, err := strconv.Atoi(arg1); err == nil && parsedCount > 0 {
				count = parsedCount
			}
		}
		return s.commandHandler.CommandDailyMarketInfo(userID, count)
	case "/sub":
		// 訂閱
		return s.commandHandler.CommandSubscribe(userID, arg1)
	case "/unsub":
		// 取消訂閱
		return s.commandHandler.CommandUnsubscribe(userID, arg1)
	case "/add":
		// 新增股票
		return s.commandHandler.CommandAddStock(userID, arg1)
	case "/del":
		// 刪除股票
		return s.commandHandler.CommandDeleteStock(userID, arg1)
	case "/list":
		// 訂閱清單
		return s.commandHandler.CommandListSubscriptions(userID)

	default:
		return nil
	}
}
