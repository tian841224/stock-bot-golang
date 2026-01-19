package bot

import (
	"context"
	"strconv"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	tgbotapi "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type TelegramCommandUsecase interface {
	GetUseGuideMessage(chatID int64) error
	GetDailyMarketInfo(ctx context.Context, chatID int64, count int) error
	GetStockPerformance(ctx context.Context, symbol string, chatID int64) error
	GetStockPerformanceChart(ctx context.Context, symbol string, chatID int64) error
	GetTopVolumeStock(ctx context.Context, chatID int64) error
	GetStockPrice(ctx context.Context, symbol string, date *time.Time, chatID int64) error
	GetStockRevenueChart(ctx context.Context, symbol string, chatID int64) error
	GetHistoricalCandlesChart(ctx context.Context, symbol string, chatID int64) error
	GetStockCompanyInfo(ctx context.Context, symbol string, chatID int64) error
	GetStockNews(ctx context.Context, symbol string, chatID int64) error
	SubscribeStock(ctx context.Context, chatID int64, symbol string) error
	UnsubscribeStock(ctx context.Context, chatID int64, symbol string) error
	SubscribedItems(ctx context.Context, chatID int64, item valueobject.SubscriptionType) error
	UnsubscribedItems(ctx context.Context, chatID int64, item valueobject.SubscriptionType) error
	GetSubscribed(ctx context.Context, chatID int64) error
}

var _ TelegramCommandUsecase = (*telegramCommandUsecase)(nil)

type telegramCommandUsecase struct {
	formatterPort     port.FormatterPort
	marketDataUsecase stock.MarketDataUsecase
	botCommandUsecase BotCommandUsecase
	userAccountPort   port.UserAccountPort
	client            *tgbotapi.TgBotClient
	logger            logger.Logger
}

var UserTypeTelegram = valueobject.UserTypeTelegram

func NewTgBotCommandUsecase(
	formatterPort port.FormatterPort,
	botCommandUsecase BotCommandUsecase,
	marketDataUsecase stock.MarketDataUsecase,
	userAccountPort port.UserAccountPort,
	client *tgbotapi.TgBotClient,
	log logger.Logger,
) TelegramCommandUsecase {
	return &telegramCommandUsecase{formatterPort: formatterPort, botCommandUsecase: botCommandUsecase, marketDataUsecase: marketDataUsecase, userAccountPort: userAccountPort, client: client, logger: log}
}

func (u *telegramCommandUsecase) GetUseGuideMessage(chatID int64) error {
	message := u.botCommandUsecase.GetUseGuideMessage()
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetDailyMarketInfo(ctx context.Context, chatID int64, count int) error {
	message, err := u.botCommandUsecase.GetDailyMarketInfo(ctx, UserTypeTelegram, count)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPerformance(ctx context.Context, symbol string, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockPerformance(ctx, UserTypeTelegram, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPerformanceChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetStockPerformanceChart(ctx, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	if chart == nil {
		return u.sendError(chatID, "圖表資料為空")
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetTopVolumeStock(ctx context.Context, chatID int64) error {
	message, err := u.botCommandUsecase.GetTopVolumeStock(ctx, UserTypeTelegram)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPrice(ctx context.Context, symbol string, date *time.Time, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockPrice(ctx, UserTypeTelegram, symbol, date)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockRevenueChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetStockRevenueChart(ctx, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	if chart == nil {
		return u.sendError(chatID, "圖表資料為空")
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetHistoricalCandlesChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetHistoricalCandlesChart(ctx, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	if chart == nil {
		return u.sendError(chatID, "圖表資料為空")
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetStockCompanyInfo(ctx context.Context, symbol string, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockCompanyInfo(ctx, UserTypeTelegram, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockNews(ctx context.Context, symbol string, chatID int64) error {
	newsMessage, err := u.botCommandUsecase.GetStockNewsForTelegram(ctx, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	if newsMessage.InlineKeyboardMarkup == nil {
		return u.client.SendMessage(chatID, newsMessage.Text)
	}

	return u.client.SendMessageWithKeyboard(chatID, newsMessage.Text, newsMessage.InlineKeyboardMarkup)
}

func (u *telegramCommandUsecase) SubscribeStock(ctx context.Context, chatID int64, symbol string) error {
	userID, err := u.getUser(ctx, strconv.FormatInt(chatID, 10))
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	result, err := u.botCommandUsecase.SubscribeStock(ctx, userID, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, result)
}

func (u *telegramCommandUsecase) UnsubscribeStock(ctx context.Context, chatID int64, symbol string) error {
	userID, err := u.getUser(ctx, strconv.FormatInt(chatID, 10))
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	result, err := u.botCommandUsecase.UnsubscribeStock(ctx, userID, symbol)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, result)
}

func (u *telegramCommandUsecase) SubscribedItems(ctx context.Context, chatID int64, item valueobject.SubscriptionType) error {
	userID, err := u.getUser(ctx, strconv.FormatInt(chatID, 10))
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	result, err := u.botCommandUsecase.SubscribedItems(ctx, userID, item)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, result)
}

func (u *telegramCommandUsecase) UnsubscribedItems(ctx context.Context, chatID int64, item valueobject.SubscriptionType) error {
	userID, err := u.getUser(ctx, strconv.FormatInt(chatID, 10))
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	result, err := u.botCommandUsecase.UnsubscribedItems(ctx, userID, item)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}
	return u.client.SendMessage(chatID, result)
}
func (u *telegramCommandUsecase) GetSubscribed(ctx context.Context, chatID int64) error {
	userID, err := u.getUser(ctx, strconv.FormatInt(chatID, 10))
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	result, err := u.botCommandUsecase.GetSubscribed(ctx, userID)
	if err != nil {
		return u.sendError(chatID, err.Error())
	}

	return u.client.SendMessage(chatID, result)
}

func (u *telegramCommandUsecase) sendError(chatID int64, message string) error {
	u.logger.Warn("發送訊息失敗", logger.Int64("chat_id", chatID), logger.String("message", message))
	return u.client.SendMessage(chatID, message)
}

// ensureUser 確保使用者存在，不存在則建立
func (p *telegramCommandUsecase) getUser(ctx context.Context, accountID string) (uint, error) {
	user, err := p.userAccountPort.GetOrCreate(ctx, accountID, valueobject.UserTypeTelegram)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
