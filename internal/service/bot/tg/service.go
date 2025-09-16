package tg

import (
	"strconv"
	"strings"

	"stock-bot/internal/db/models"
	"stock-bot/internal/repository"
	"stock-bot/internal/service/stock"
	"stock-bot/internal/service/user"
	"stock-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgService struct {
	botClient            *tgbotapi.BotAPI
	stockService         *stock.StockService
	userService          user.UserService
	userSubscriptionRepo repository.UserSubscriptionRepository
	subscriptionItemMap  map[string]models.SubscriptionItem
}

func NewTgService(
	botClient *tgbotapi.BotAPI,
	stockService *stock.StockService,
	userService user.UserService,
	userSubscriptionRepo repository.UserSubscriptionRepository,
) *TgService {
	return &TgService{
		botClient:            botClient,
		stockService:         stockService,
		userService:          userService,
		userSubscriptionRepo: userSubscriptionRepo,
		subscriptionItemMap:  models.SubscriptionItemMap,
	}
}

func (s *TgService) HandleUpdate(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	return s.handleCommand(update.Message)
}

func (s *TgService) handleCommand(message *tgbotapi.Message) error {
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
		return s.sendMessage(userID, "系統錯誤，請稍後再試")
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
		return s.handleStart(userID)
	case "/k":
		return s.handleKline(userID, arg1, arg2)
	case "/p":
		return s.handlePerformance(userID, arg1)
	case "/d":
		return s.handleDetailPrice(userID, arg1)
	case "/n":
		return s.handleNews(userID, arg1)
	case "/yn":
		return s.handleYahooNews(userID, arg1)
	case "/m":
		count := 1
		if arg1 != "" {
			if c, err := strconv.Atoi(arg1); err == nil {
				count = c
			}
		}
		return s.handleDailyMarketInfo(userID, count)
	case "/t":
		return s.handleTopVolumeItems(userID)
	case "/i":
		return s.handleStockInfo(userID, arg1, arg2)
	case "/sub":
		return s.handleSubscribe(userID, arg1)
	case "/unsub":
		return s.handleUnsubscribe(userID, arg1)
	case "/add":
		return s.handleAddStock(userID, arg1)
	case "/del":
		return s.handleDeleteStock(userID, arg1)
	case "/list":
		return s.handleListSubscriptions(userID)
	default:
		return nil
	}
}

// 輔助方法

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

func (s *TgService) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.botClient.Send(msg)
	if err != nil {
		logger.Log.Error("發送訊息失敗", zap.Error(err))
	}
	return err
}

func (s *TgService) sendMessageHTML(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := s.botClient.Send(msg)
	if err != nil {
		logger.Log.Error("發送 HTML 訊息失敗", zap.Error(err))
	}
	return err
}
