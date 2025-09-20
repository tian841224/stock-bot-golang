package tg

import (
	"strconv"
	"strings"

	"stock-bot/internal/db/models"
	"stock-bot/internal/service/user"
	"stock-bot/pkg/logger"

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
		return s.commandHandler.CommandKline(userID, arg1, arg2)
	case "/p":
		return s.commandHandler.CommandPerformance(userID, arg1)
	case "/pc":
		return s.commandHandler.CommandPerformanceChart(userID, arg1)
	case "/pb":
		return s.commandHandler.CommandPerformanceBarChart(userID, arg1)
	case "/d":
		return s.commandHandler.CommandTodayStockPrice(userID, arg1, arg2)
	case "/t":
		return s.commandHandler.CommandTopVolumeItems(userID)
	case "/i":
		return s.commandHandler.CommandStockInfo(userID, arg1, arg2)
	case "/r":
		return s.commandHandler.CommandRevenue(userID, arg1)
	// case "/n":
	// 	return s.commandHandler.sendMessage(userID, "新聞功能暫時停用")
	// case "/yn":
	// 	return s.commandHandler.sendMessage(userID, "Yahoo新聞功能暫時停用")
	// case "/m":
	// 	count := 1
	// 	if arg1 != "" {
	// 		if c, err := strconv.Atoi(arg1); err == nil {
	// 			count = c
	// 		}
	// 	}
	// 	return s.commandHandler.CommandDailyMarketInfo(userID, count)

	case "/sub":
		return s.commandHandler.CommandSubscribe(userID, arg1)
	case "/unsub":
		return s.commandHandler.CommandUnsubscribe(userID, arg1)
	case "/add":
		return s.commandHandler.CommandAddStock(userID, arg1)
	case "/del":
		return s.commandHandler.CommandDeleteStock(userID, arg1)
	case "/list":
		return s.commandHandler.CommandListSubscriptions(userID)
	default:
		return nil
	}
}
