package bot

import (
	"context"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	tgbotapi "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
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
}

var _ TelegramCommandUsecase = (*telegramCommandUsecase)(nil)

type telegramCommandUsecase struct {
	formatterPort     port.FormatterPort
	marketDataUsecase stock.MarketDataUsecase
	botCommandUsecase BotCommandUsecase
	client            *tgbotapi.TgBotClient
}

var UserTypeTelegram = valueobject.UserTypeTelegram

func NewTgBotCommandUsecase(
	formatterPort port.FormatterPort,
	botCommandUsecase BotCommandUsecase,
	marketDataUsecase stock.MarketDataUsecase,
	client *tgbotapi.TgBotClient,
) TelegramCommandUsecase {
	return &telegramCommandUsecase{formatterPort: formatterPort, botCommandUsecase: botCommandUsecase, marketDataUsecase: marketDataUsecase, client: client}
}

func (u *telegramCommandUsecase) GetUseGuideMessage(chatID int64) error {
	message := u.botCommandUsecase.GetUseGuideMessage()
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetDailyMarketInfo(ctx context.Context, chatID int64, count int) error {
	message, err := u.botCommandUsecase.GetDailyMarketInfo(ctx, UserTypeTelegram, count)
	if err != nil {
		return err
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPerformance(ctx context.Context, symbol string, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockPerformance(ctx, UserTypeTelegram, symbol)
	if err != nil {
		return err
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPerformanceChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetStockPerformanceChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return nil
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetTopVolumeStock(ctx context.Context, chatID int64) error {
	message, err := u.botCommandUsecase.GetTopVolumeStock(ctx, UserTypeTelegram)
	if err != nil {
		return err
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockPrice(ctx context.Context, symbol string, date *time.Time, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockPrice(ctx, UserTypeTelegram, symbol, date)
	if err != nil {
		return err
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockRevenueChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetStockRevenueChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return nil
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetHistoricalCandlesChart(ctx context.Context, symbol string, chatID int64) error {
	chart, err := u.botCommandUsecase.GetHistoricalCandlesChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return nil
	}

	return u.client.SendPhoto(chatID, chart.Data, chart.FileName)
}

func (u *telegramCommandUsecase) GetStockCompanyInfo(ctx context.Context, symbol string, chatID int64) error {
	message, err := u.botCommandUsecase.GetStockCompanyInfo(ctx, UserTypeTelegram, symbol)
	if err != nil {
		return err
	}
	return u.client.SendMessage(chatID, message)
}

func (u *telegramCommandUsecase) GetStockNews(ctx context.Context, symbol string, chatID int64) error {
	newsMessage, err := u.botCommandUsecase.GetStockNewsForTelegram(ctx, symbol)
	if err != nil {
		return err
	}

	if newsMessage == nil {
		return u.client.SendMessage(chatID, "無法取得新聞資料")
	}

	if newsMessage.InlineKeyboardMarkup == nil {
		return u.client.SendMessage(chatID, newsMessage.Text)
	}

	return u.client.SendMessageWithKeyboard(chatID, newsMessage.Text, newsMessage.InlineKeyboardMarkup)
}
