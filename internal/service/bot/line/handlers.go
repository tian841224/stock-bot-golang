package line

import (
	"strconv"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/db/models"
	"github.com/tian841224/stock-bot/internal/service/user"
	"github.com/tian841224/stock-bot/pkg/logger"

	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/linebot"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

// LineServiceHandler LINE 服務處理器介面
type LineServiceHandler interface {
	HandleTextMessage(event *linebot.Event, message *linebot.TextMessage) error
}

// lineServiceHandler 處理對話邏輯
type lineServiceHandler struct {
	botClient      *linebotInfra.LineBotClient
	commandHandler *LineCommandHandler
	userService    user.UserService
}

// NewBotService 創建 service
func NewBotService(
	botClient *linebotInfra.LineBotClient,
	commandHandler *LineCommandHandler,
	userService user.UserService,
) LineServiceHandler {
	return &lineServiceHandler{
		botClient:      botClient,
		commandHandler: commandHandler,
		userService:    userService,
	}
}

// parseMessageArgs 解析訊息參數
func parseMessageArgs(messageText string) (command string, arg1 string, arg2 string) {
	parts := strings.Fields(messageText)
	if len(parts) == 0 {
		return "", "", ""
	}

	command = parts[0]
	if len(parts) > 1 {
		arg1 = parts[1]
	}
	if len(parts) > 2 {
		arg2 = parts[2]
	}
	return command, arg1, arg2
}

// getDefaultDateForTodayPrice 取得今日股價的預設日期
func getDefaultDateForTodayPrice() string {
	taipeiLocation, _ := time.LoadLocation("Asia/Taipei")
	now := time.Now().In(taipeiLocation)
	marketOpenTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, taipeiLocation)

	if now.Before(marketOpenTime) {
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	}
	return now.Format("2006-01-02")
}

// parseMarketInfoCount 解析大盤資訊的顯示筆數
func parseMarketInfoCount(arg1 string) int {
	count := 1
	if arg1 == "" {
		return count
	}

	parsedCount, err := strconv.Atoi(arg1)
	if err == nil && parsedCount > 0 {
		count = parsedCount
	}
	return count
}

// HandleTextMessage 處理文字訊息
func (s *lineServiceHandler) HandleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if message.Text == "" {
		return nil
	}

	userID := event.Source.UserID
	messageText := message.Text

	logger.Log.Info("收到 LINE 訊息",
		zap.String("user_id", userID),
		zap.String("message", messageText))

	_, err := s.userService.GetOrCreate(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("建立或取得使用者失敗", zap.Error(err))
		return s.botClient.ReplyMessage(event.ReplyToken, "系統錯誤，請稍後再試")
	}

	command, arg1, arg2 := parseMessageArgs(messageText)
	if command == "" {
		return nil
	}

	return s.executeCommand(command, userID, event.ReplyToken, arg1, arg2, messageText)
}

// executeCommand 執行對應的命令
func (s *lineServiceHandler) executeCommand(command, userID, replyToken, arg1, arg2, messageText string) error {
	commandMap := map[string]func() error{
		"/start": func() error {
			return s.commandHandler.CommandStart(replyToken)
		},
		"/k": func() error {
			return s.commandHandler.CommandHistoricalCandles(replyToken, arg1)
		},
		"/p": func() error {
			return s.commandHandler.CommandPerformanceChart(replyToken, arg1)
		},
		"/d": func() error {
			dateArg := arg2
			if dateArg == "" {
				dateArg = getDefaultDateForTodayPrice()
			}
			return s.commandHandler.CommandTodayStockPrice(replyToken, arg1, dateArg)
		},
		"/t": func() error {
			return s.commandHandler.CommandTopVolumeItems(replyToken)
		},
		"/i": func() error {
			return s.commandHandler.CommandStockInfo(replyToken, arg1, arg2)
		},
		"/r": func() error {
			return s.commandHandler.CommandRevenue(replyToken, arg1)
		},
		"/n": func() error {
			return s.commandHandler.CommandNews(replyToken, arg1)
		},
		"/m": func() error {
			count := parseMarketInfoCount(arg1)
			return s.commandHandler.CommandDailyMarketInfo(replyToken, count)
		},
		"/sub": func() error {
			return s.commandHandler.CommandSubscribe(userID, replyToken, arg1)
		},
		"/unsub": func() error {
			return s.commandHandler.CommandUnsubscribe(userID, replyToken, arg1)
		},
		"/add": func() error {
			return s.commandHandler.CommandAddStock(userID, replyToken, arg1)
		},
		"/del": func() error {
			return s.commandHandler.CommandDeleteStock(userID, replyToken, arg1)
		},
		"/list": func() error {
			return s.commandHandler.CommandListSubscriptions(userID, replyToken)
		},
		"test": func() error {
			return s.botClient.ReplyMessage(replyToken, "新增成功")
		},
	}

	if handler, exists := commandMap[command]; exists {
		return handler()
	}

	reply := "你說了: " + messageText
	return s.botClient.ReplyMessage(replyToken, reply)
}
