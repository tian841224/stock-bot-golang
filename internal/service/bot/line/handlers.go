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

// LineServiceHandler 處理對話邏輯
type LineServiceHandler struct {
	botClient      *linebotInfra.LineBotClient
	commandHandler *LineCommandHandler
	userService    user.UserService
}

// NewBotService 創建 service
func NewBotService(
	botClient *linebotInfra.LineBotClient,
	commandHandler *LineCommandHandler,
	userService user.UserService,
) *LineServiceHandler {
	return &LineServiceHandler{
		botClient:      botClient,
		commandHandler: commandHandler,
		userService:    userService,
	}
}

// HandleTextMessage 處理文字訊息
func (s *LineServiceHandler) HandleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if message.Text == "" {
		return nil
	}

	userID := event.Source.UserID
	messageText := message.Text

	logger.Log.Info("收到 LINE 訊息",
		zap.String("user_id", userID),
		zap.String("message", messageText))

	// 確保使用者存在
	_, err := s.userService.GetOrCreate(userID, models.UserTypeLine)
	if err != nil {
		logger.Log.Error("建立或取得使用者失敗", zap.Error(err))
		return s.botClient.ReplyMessage(event.ReplyToken, "系統錯誤，請稍後再試")
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
		return s.commandHandler.CommandStart(event.ReplyToken)
	case "/k":
		// 歷史K線圖
		return s.commandHandler.CommandHistoricalCandles(event.ReplyToken, arg1)
	case "/p":
		// 績效圖表
		return s.commandHandler.CommandPerformanceChart(event.ReplyToken, arg1)
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
		return s.commandHandler.CommandTodayStockPrice(event.ReplyToken, arg1, arg2)
	case "/t":
		// 交易量前20名
		return s.commandHandler.CommandTopVolumeItems(event.ReplyToken)
	case "/i":
		// 股票資訊
		return s.commandHandler.CommandStockInfo(event.ReplyToken, arg1, arg2)
	case "/r":
		// 財報
		return s.commandHandler.CommandRevenue(event.ReplyToken, arg1)
	case "/n":
		// 新聞
		return s.commandHandler.CommandNews(event.ReplyToken, arg1)
	case "/m":
		// 大盤資訊
		count := 1 // 預設顯示1筆
		if arg1 != "" {
			if parsedCount, err := strconv.Atoi(arg1); err == nil && parsedCount > 0 {
				count = parsedCount
			}
		}
		return s.commandHandler.CommandDailyMarketInfo(event.ReplyToken, count)
	case "/sub":
		// 訂閱
		return s.commandHandler.CommandSubscribe(userID, event.ReplyToken, arg1)
	case "/unsub":
		// 取消訂閱
		return s.commandHandler.CommandUnsubscribe(userID, event.ReplyToken, arg1)
	case "/add":
		// 新增股票
		return s.commandHandler.CommandAddStock(userID, event.ReplyToken, arg1)
	case "/del":
		// 刪除股票
		return s.commandHandler.CommandDeleteStock(userID, event.ReplyToken, arg1)
	case "/list":
		// 訂閱清單
		return s.commandHandler.CommandListSubscriptions(userID, event.ReplyToken)
	case "test":
		return s.botClient.ReplyMessage(event.ReplyToken, "新增成功")
	default:
		reply := "你說了: " + message.Text
		return s.botClient.ReplyMessage(event.ReplyToken, reply)
	}
}
