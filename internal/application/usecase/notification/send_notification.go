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
			u.logger.Error("get stock price failed",
				logger.String("op", "SendStockPriceNotification"),
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatStockPrice(stockPrice, valueobject.UserTypeTelegram)

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("parse account_id failed",
				logger.String("op", "SendStockPriceNotification"),
				logger.String("account_id", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("send message failed",
				logger.String("op", "SendStockPriceNotification"),
				logger.Int64("account_id", accountID),
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
			u.logger.Error("get stock news failed",
				logger.String("op", "SendStockNewsNotification"),
				logger.String("symbol", subscriptionSymbol.StockSymbol.Symbol),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatTelegramNewsMessage(*stockNews, subscriptionSymbol.StockSymbol.Symbol, subscriptionSymbol.StockSymbol.Name)

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("parse account_id failed",
				logger.String("op", "SendStockNewsNotification"),
				logger.String("account_id", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessageWithKeyboard(accountID, data.Text, data.InlineKeyboardMarkup)
		if err != nil {
			u.logger.Error("send message failed",
				logger.String("op", "SendStockNewsNotification"),
				logger.Int64("account_id", accountID),
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
			u.logger.Error("get daily market info failed",
				logger.String("op", "SendMarketInfoNotification"),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatDailyMarketInfo(stockPrice, valueobject.UserTypeTelegram)

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("parse account_id failed",
				logger.String("op", "SendMarketInfoNotification"),
				logger.String("account_id", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("send message failed",
				logger.String("op", "SendMarketInfoNotification"),
				logger.Int64("account_id", accountID),
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
			u.logger.Error("get top volume stocks failed",
				logger.String("op", "SendTopVolumeNotification"),
				logger.Error(err),
			)
			continue
		}

		data := u.formatterPort.FormatTopVolumeStock(topVolumeStocks, valueobject.UserTypeTelegram)

		accountID, err := strconv.ParseInt(subscriptionSymbol.User.AccountID, 10, 64)
		if err != nil {
			u.logger.Error("parse account_id failed",
				logger.String("op", "SendTopVolumeNotification"),
				logger.String("account_id", subscriptionSymbol.User.AccountID),
				logger.Error(err),
			)
			continue
		}

		err = u.client.SendMessage(accountID, data)
		if err != nil {
			u.logger.Error("send message failed",
				logger.String("op", "SendTopVolumeNotification"),
				logger.Int64("account_id", accountID),
				logger.Error(err),
			)
			continue
		}
	}
	return nil
}
