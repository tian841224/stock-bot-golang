// Package tg 提供 Telegram Bot 的處理器功能
package tgbot

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

// TgServiceHandler Telegram 服務處理器介面
type TgServiceHandler interface {
	ProcessUpdate(update *tgbotapi.Update) error
}

type tgServiceHandler struct {
	commandHandler *TgCommandHandler
	userService    user.UserService
}

func NewTgServiceHandler(commandHandler *TgCommandHandler, userService user.UserService) TgServiceHandler {
	return &tgServiceHandler{
		commandHandler: commandHandler,
		userService:    userService,
	}
}

func (s *tgServiceHandler) ProcessUpdate(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	return s.processCommand(update.Message)
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

func (s *tgServiceHandler) processCommand(message *tgbotapi.Message) error {
	if message.Text == "" {
		return nil
	}

	userID := message.Chat.ID
	messageText := message.Text

	logger.Log.Info("收到 Telegram 訊息",
		zap.Int64("user_id", userID),
		zap.String("message", messageText))

	_, err := s.userService.GetOrCreate(strconv.FormatInt(userID, 10), models.UserTypeTelegram)
	if err != nil {
		logger.Log.Error("建立或取得使用者失敗", zap.Error(err))
		return s.commandHandler.botClient.SendMessage(userID, "系統錯誤，請稍後再試")
	}

	command, arg1, arg2 := parseMessageArgs(messageText)
	if command == "" {
		return nil
	}

	return s.executeCommand(command, userID, arg1, arg2)
}

// executeCommand 執行對應的命令
func (s *tgServiceHandler) executeCommand(command string, userID int64, arg1, arg2 string) error {
	commandMap := map[string]func() error{
		"/start": func() error {
			return s.commandHandler.CommandStart(userID)
		},
		"/k": func() error {
			return s.commandHandler.CommandHistoricalCandles(userID, arg1)
		},
		"/p": func() error {
			return s.commandHandler.CommandPerformanceChart(userID, arg1)
		},
		"/d": func() error {
			dateArg := arg2
			if dateArg == "" {
				dateArg = getDefaultDateForTodayPrice()
			}
			return s.commandHandler.CommandTodayStockPrice(userID, arg1, dateArg)
		},
		"/t": func() error {
			return s.commandHandler.CommandTopVolumeItems(userID)
		},
		"/i": func() error {
			return s.commandHandler.CommandStockInfo(userID, arg1, arg2)
		},
		"/r": func() error {
			return s.commandHandler.CommandRevenue(userID, arg1)
		},
		"/n": func() error {
			return s.commandHandler.CommandNews(userID, arg1)
		},
		"/m": func() error {
			count := parseMarketInfoCount(arg1)
			return s.commandHandler.CommandDailyMarketInfo(userID, count)
		},
		"/sub": func() error {
			return s.commandHandler.CommandSubscribe(userID, arg1)
		},
		"/unsub": func() error {
			return s.commandHandler.CommandUnsubscribe(userID, arg1)
		},
		"/add": func() error {
			return s.commandHandler.CommandAddStock(userID, arg1)
		},
		"/del": func() error {
			return s.commandHandler.CommandDeleteStock(userID, arg1)
		},
		"/list": func() error {
			return s.commandHandler.CommandListSubscriptions(userID)
		},
	}

	if handler, exists := commandMap[command]; exists {
		return handler()
	}

	return nil
}
