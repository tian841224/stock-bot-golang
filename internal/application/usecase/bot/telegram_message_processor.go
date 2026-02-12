package bot

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	tgbotapi "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramMessageProcessor 處理 Telegram 訊息的路由和編排
type TelegramMessageProcessor struct {
	tgCommandUsecase TelegramCommandUsecase
	userAccountPort  port.UserAccountPort
	tgClient         *tgbotapi.TgBotClient
	logger           logger.Logger
}

func NewTelegramMessageProcessor(
	tgCommandUsecase TelegramCommandUsecase,
	userAccountPort port.UserAccountPort,
	tgClient *tgbotapi.TgBotClient,
	log logger.Logger,
) *TelegramMessageProcessor {
	return &TelegramMessageProcessor{
		tgCommandUsecase: tgCommandUsecase,
		userAccountPort:  userAccountPort,
		tgClient:         tgClient,
		logger:           log,
	}
}

// ProcessUpdate 處理 Telegram update，包含命令路由
func (p *TelegramMessageProcessor) ProcessUpdate(ctx context.Context, update *tgbot.Update) error {
	if update.Message == nil || update.Message.Text == "" {
		return nil
	}

	chatID := update.Message.Chat.ID
	messageText := update.Message.Text

	p.logger.Info("收到 Telegram 訊息",
		logger.Int64("chat_id", chatID),
		logger.String("message", messageText))

	// 解析命令和參數
	command, arg1, arg2 := p.parseMessageArgs(messageText)
	if command == "" {
		return nil
	}

	// 路由到對應的命令處理器
	if err := p.routeCommand(ctx, command, arg1, arg2, chatID); err != nil {
		p.logger.Error("處理命令失敗",
			logger.String("command", command),
			logger.String("arg1", arg1),
			logger.String("arg2", arg2),
			logger.Int64("chat_id", chatID),
			logger.Error(err))

		// 發送錯誤訊息給使用者
		errorMsg := err.Error()
		if errorMsg == "" {
			errorMsg = "處理請求時發生錯誤，請稍後再試"
		}
		return p.sendError(chatID, errorMsg)
	}
	return nil
}

// routeCommand 路由命令到對應的處理器
func (p *TelegramMessageProcessor) routeCommand(ctx context.Context, command, arg1, arg2 string, chatID int64) error {
	switch command {
	case "/start":
		return p.tgCommandUsecase.GetUseGuideMessage(chatID)
	case "/k":
		return p.handleHistoricalCandles(ctx, chatID, arg1)
	case "/p":
		return p.handlePerformanceChart(ctx, chatID, arg1)
	case "/d":
		return p.handleStockPrice(ctx, chatID, arg1, arg2)
	case "/t":
		return p.tgCommandUsecase.GetTopVolumeStock(ctx, chatID)
	case "/i":
		return p.tgCommandUsecase.GetStockCompanyInfo(ctx, arg1, chatID)
	case "/r":
		return p.handleRevenueChart(ctx, chatID, arg1)
	case "/m":
		return p.handleDailyMarket(ctx, chatID, arg1)
	case "/n":
		return p.tgCommandUsecase.GetStockNews(ctx, arg1, chatID)
	case "/sub":
		return p.handleSubscribedItems(ctx, chatID, arg1)
	case "/unsub":
		return p.handleUnsubscribedItems(ctx, chatID, arg1)
	case "/add":
		return p.tgCommandUsecase.SubscribeStock(ctx, chatID, arg1)
	case "/del":
		return p.tgCommandUsecase.UnsubscribeStock(ctx, chatID, arg1)
	case "/list":
		return p.tgCommandUsecase.GetSubscribed(ctx, chatID)
	default:
		// return p.handleUnknownCommand(chatID)
	}
	return nil
}

// 各個命令的具體處理邏輯

func (p *TelegramMessageProcessor) handleHistoricalCandles(ctx context.Context, chatID int64, symbol string) error {
	if symbol == "" {
		return p.sendError(chatID, "請輸入股票代號\n\n使用方式：\n/k 股票代號 - 查詢K線圖")
	}
	return p.tgCommandUsecase.GetHistoricalCandlesChart(ctx, symbol, chatID)
}

func (p *TelegramMessageProcessor) handlePerformanceChart(ctx context.Context, chatID int64, symbol string) error {
	if symbol == "" {
		return p.sendError(chatID, "請輸入股票代號\n\n使用方式：\n/p 股票代號 - 查詢績效圖表")
	}
	return p.tgCommandUsecase.GetStockPerformanceChart(ctx, symbol, chatID)
}

func (p *TelegramMessageProcessor) handleStockPrice(ctx context.Context, chatID int64, symbol, rawDate string) error {
	if symbol == "" {
		return p.sendError(chatID, "請輸入股票代號\n\n使用方式：\n/d 股票代號 - 查詢今日股價\n/d 股票代號 2025-12-09 - 查詢指定日期股價")
	}

	var datePtr *time.Time
	if rawDate != "" {
		parsed, err := p.parseDate(rawDate)
		if err != nil {
			return p.sendError(chatID, "日期格式錯誤，請使用 YYYY-MM-DD 格式\n例如：2025-12-09")
		}
		datePtr = &parsed
		return p.tgCommandUsecase.GetStockPrice(ctx, symbol, datePtr, chatID)
	}

	// 如果為空則取今天
	now := time.Now()
	datePtr = &now
	if now.Hour() < 14 {
		now = now.AddDate(0, 0, -1)
		datePtr = &now
	}

	return p.tgCommandUsecase.GetStockPrice(ctx, symbol, datePtr, chatID)
}

func (p *TelegramMessageProcessor) handleRevenueChart(ctx context.Context, chatID int64, symbol string) error {
	if symbol == "" {
		return p.sendError(chatID, "請輸入股票代號\n\n使用方式：\n/r 股票代號 - 查詢月營收圖表")
	}
	return p.tgCommandUsecase.GetStockRevenueChart(ctx, symbol, chatID)
}

func (p *TelegramMessageProcessor) handleDailyMarket(ctx context.Context, chatID int64, countStr string) error {
	count := 1
	if countStr != "" {
		countInt, err := strconv.Atoi(countStr)
		if countInt <= 0 || err != nil {
			return p.sendError(chatID, "請輸入有效的數字，且大於0\n\n使用方式：\n/m [數量] - 查詢指定筆數的大盤資訊")
		}
		count = countInt
	}
	return p.tgCommandUsecase.GetDailyMarketInfo(ctx, chatID, count)
}

func (p *TelegramMessageProcessor) handleSubscribedItems(ctx context.Context, chatID int64, args string) error {
	intValue, err := strconv.Atoi(args)
	if err != nil {
		return p.sendError(chatID, "請輸入有效的訂閱類型")
	}

	item, err := valueobject.NewSubscriptionType(intValue)
	if err != nil {
		return p.sendError(chatID, "請輸入有效的訂閱類型")
	}

	return p.tgCommandUsecase.SubscribedItems(ctx, chatID, item)
}

func (p *TelegramMessageProcessor) handleUnsubscribedItems(ctx context.Context, chatID int64, args string) error {
	intValue, err := strconv.Atoi(args)
	if err != nil {
		return p.sendError(chatID, "請輸入有效的訂閱類型")
	}

	item, err := valueobject.NewSubscriptionType(intValue)
	if err != nil {
		return p.sendError(chatID, "請輸入有效的訂閱類型")
	}

	return p.tgCommandUsecase.UnsubscribedItems(ctx, chatID, item)
}

// func (p *TelegramMessageProcessor) handleUnknownCommand(chatID int64) error {
// 	return p.sendError(chatID, "指令不存在，輸入 /start 查看說明")
// }

// 輔助方法

func (p *TelegramMessageProcessor) sendError(chatID int64, message string) error {
	p.logger.Warn("發送錯誤訊息",
		logger.Int64("chat_id", chatID),
		logger.String("message", message))
	return p.tgClient.SendMessage(chatID, message)
}

func (p *TelegramMessageProcessor) parseMessageArgs(messageText string) (command, arg1, arg2 string) {
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

func (p *TelegramMessageProcessor) parseDate(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
