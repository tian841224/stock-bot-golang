package notification

import (
	"context"
	"strconv"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/application/usecase/stock"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	tgbotapi "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/telegram"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type SendNotificationUsecase interface {
	SendStockPriceNotification(ctx context.Context) error
	SendStockNewsNotification(ctx context.Context) error
	SendMarketInfoNotification(ctx context.Context) error
	SendTopVolumeNotification(ctx context.Context) error
}

type sendNotificationUsecase struct {
	marketDataUsecase      stock.MarketDataUsecase
	subscriptionSymbolRepo port.SubscriptionSymbolRepository
	formatterPort          port.FormatterPort
	client                 *tgbotapi.TgBotClient
	logger                 logger.Logger
}

func NewSendNotificationUsecase(
	subscriptionSymbolRepo port.SubscriptionSymbolRepository,
	marketDataUsecase stock.MarketDataUsecase,
	formatterPort port.FormatterPort,
	client *tgbotapi.TgBotClient,
	log logger.Logger,
) SendNotificationUsecase {
	return &sendNotificationUsecase{
		subscriptionSymbolRepo: subscriptionSymbolRepo,
		formatterPort:          formatterPort,
		marketDataUsecase:      marketDataUsecase,
		client:                 client,
		logger:                 log,
	}
}

// SendStockPriceNotification 推送股票股價
func (u *sendNotificationUsecase) SendStockPriceNotification(ctx context.Context) error {
	subscriptionSymbols, err := u.subscriptionSymbolRepo.GetByFeature(ctx, valueobject.SubscriptionTypeStockInfo)
	if err != nil {
		return err
	}

	for _, subscriptionSymbol := range subscriptionSymbols {
		stockPrice, err := u.marketDataUsecase.GetStockPrice(ctx, subscriptionSymbol.StockSymbol.Symbol, nil)
		if err != nil {
			u.logger.Error("SendStockPriceNotification GetStockPrice Error",
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatStockPrice(stockPrice, valueobject.UserTypeTelegram)
		if err != nil {
			u.logger.Error("SendStockPriceNotification FormatStockPrice Error",
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("SendStockPriceNotification ParseInt Error",
				logger.String("accountID", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("SendStockPriceNotification SendMessage Error",
				logger.Int64("accountID", accountID),
				logger.String("data", data),
				logger.Error(err),
			)
			continue
		}
	}
	return nil
}

// SendStockNewsNotification 推送股票新聞
func (u *sendNotificationUsecase) SendStockNewsNotification(ctx context.Context) error {
	subscriptionSymbols, err := u.subscriptionSymbolRepo.GetByFeature(ctx, valueobject.SubscriptionTypeStockNews)
	if err != nil {
		return err
	}

	for _, subscriptionSymbol := range subscriptionSymbols {
		stockNews, err := u.marketDataUsecase.GetStockNews(ctx, subscriptionSymbol.StockSymbol.Symbol, 5)
		if err != nil {
			u.logger.Error("SendStockNewsNotification GetStockNews Error",
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatTelegramNewsMessage(*stockNews, subscriptionSymbol.StockSymbol.Symbol, subscriptionSymbol.StockSymbol.Name)
		if err != nil {
			u.logger.Error("SendStockNewsNotification FormatTelegramNewsMessage Error",
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("SendStockNewsNotification ParseInt Error",
				logger.String("accountID", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessageWithKeyboard(accountID, data.Text, data.InlineKeyboardMarkup)
		if err != nil {
			u.logger.Error("SendStockNewsNotification SendMessageWithKeyboard Error",
				logger.Int64("accountID", accountID),
				logger.String("data", data.Text),
				logger.Error(err),
			)
			continue
		}
	}
	return nil
}

// SendMarketInfoNotification 推送大盤資訊
func (u *sendNotificationUsecase) SendMarketInfoNotification(ctx context.Context) error {
	subscriptionSymbols, err := u.subscriptionSymbolRepo.GetByFeature(ctx, valueobject.SubscriptionTypeDailyMarketInfo)
	if err != nil {
		return err
	}

	for _, subscriptionSymbol := range subscriptionSymbols {
		stockPrice, err := u.marketDataUsecase.GetDailyMarketInfo(ctx, 1)
		if err != nil {
			u.logger.Error("SendMarketInfoNotification GetDailyMarketInfo Error",
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatDailyMarketInfo(stockPrice, valueobject.UserTypeTelegram)
		if err != nil {
			u.logger.Error("SendMarketInfoNotification FormatDailyMarketInfo Error",
				logger.Error(err),
			)
			continue
		}

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("SendMarketInfoNotification ParseInt Error",
				logger.String("accountID", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("SendMarketInfoNotification SendMessage Error",
				logger.Int64("accountID", accountID),
				logger.String("data", data),
				logger.Error(err),
			)
			continue
		}
	}
	return nil
}

// SendTopVolumeNotification 推送交易量排行
func (u *sendNotificationUsecase) SendTopVolumeNotification(ctx context.Context) error {
	subscriptionSymbols, err := u.subscriptionSymbolRepo.GetByFeature(ctx, valueobject.SubscriptionTypeTopVolumeItems)
	if err != nil {
		return err
	}

	for _, subscriptionSymbol := range subscriptionSymbols {
		topVolumeStocks, err := u.marketDataUsecase.GetTopVolumeStock(ctx)
		if err != nil {
			u.logger.Error("SendTopVolumeNotification GetTopVolumeStock Error",
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatTopVolumeStock(topVolumeStocks, valueobject.UserTypeTelegram)
		if err != nil {
			u.logger.Error("SendTopVolumeNotification FormatTopVolumeStock Error",
				logger.Error(err),
			)
			continue
		}

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("SendTopVolumeNotification ParseInt Error",
				logger.String("accountID", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("SendTopVolumeNotification SendMessage Error",
				logger.Int64("accountID", accountID),
				logger.String("data", data),
				logger.Error(err),
			)
			continue
		}
	}
	return nil
}
