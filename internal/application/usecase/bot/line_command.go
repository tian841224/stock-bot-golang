package bot

import (
	"context"
	"errors"
	"time"

	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/line"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/imgbb"
)

type LineCommandUsecase interface {
	GetUseGuideMessage(replyToken string) error
	GetDailyMarketInfo(ctx context.Context, replyToken string, count int) error
	GetStockPerformance(ctx context.Context, symbol string, replyToken string) error
	GetStockPerformanceChart(ctx context.Context, symbol string, replyToken string) error
	GetTopVolumeStock(ctx context.Context, replyToken string) error
	GetStockPrice(ctx context.Context, symbol string, date *time.Time, replyToken string) error
	GetStockRevenueChart(ctx context.Context, symbol string, replyToken string) error
	GetHistoricalCandlesChart(ctx context.Context, symbol string, replyToken string) error
	GetStockCompanyInfo(ctx context.Context, symbol string, replyToken string) error
	GetStockNews(ctx context.Context, symbol string, replyToken string) error
}

var _ LineCommandUsecase = (*lineCommandUsecase)(nil)

type lineCommandUsecase struct {
	botCommandUsecase BotCommandUsecase
	client            *linebotInfra.LineBotClient
	imgbbClient       *imgbb.ImgBBClient
}

var UserTypeLine = valueobject.UserTypeLine

func NewLineBotCommandUsecase(
	botCommandUsecase BotCommandUsecase,
	client *linebotInfra.LineBotClient,
	imgbbClient *imgbb.ImgBBClient,
) LineCommandUsecase {
	return &lineCommandUsecase{
		botCommandUsecase: botCommandUsecase,
		client:            client,
		imgbbClient:       imgbbClient,
	}
}

func (u *lineCommandUsecase) GetUseGuideMessage(replyToken string) error {
	message := u.botCommandUsecase.GetUseGuideMessage()
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetDailyMarketInfo(ctx context.Context, replyToken string, count int) error {
	message, err := u.botCommandUsecase.GetDailyMarketInfo(ctx, UserTypeLine, count)
	if err != nil {
		return err
	}
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetStockPerformance(ctx context.Context, symbol string, replyToken string) error {
	message, err := u.botCommandUsecase.GetStockPerformance(ctx, UserTypeLine, symbol)
	if err != nil {
		return err
	}
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetStockPerformanceChart(ctx context.Context, symbol string, replyToken string) error {
	chart, err := u.botCommandUsecase.GetStockPerformanceChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return errors.New("圖表資料為空")
	}

	return u.client.ReplyPhoto(replyToken, chart.Data, chart.FileName, u.imgbbClient)
}

func (u *lineCommandUsecase) GetTopVolumeStock(ctx context.Context, replyToken string) error {
	message, err := u.botCommandUsecase.GetTopVolumeStock(ctx, UserTypeLine)
	if err != nil {
		return err
	}
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetStockPrice(ctx context.Context, symbol string, date *time.Time, replyToken string) error {
	message, err := u.botCommandUsecase.GetStockPrice(ctx, UserTypeLine, symbol, date)
	if err != nil {
		return err
	}
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetStockRevenueChart(ctx context.Context, symbol string, replyToken string) error {
	chart, err := u.botCommandUsecase.GetStockRevenueChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return errors.New("圖表資料為空")
	}

	return u.client.ReplyPhoto(replyToken, chart.Data, chart.FileName, u.imgbbClient)
}

func (u *lineCommandUsecase) GetHistoricalCandlesChart(ctx context.Context, symbol string, replyToken string) error {
	chart, err := u.botCommandUsecase.GetHistoricalCandlesChart(ctx, symbol)
	if err != nil {
		return err
	}

	if chart == nil {
		return errors.New("圖表資料為空")
	}

	return u.client.ReplyPhoto(replyToken, chart.Data, chart.FileName, u.imgbbClient)
}

func (u *lineCommandUsecase) GetStockCompanyInfo(ctx context.Context, symbol string, replyToken string) error {
	message, err := u.botCommandUsecase.GetStockCompanyInfo(ctx, UserTypeLine, symbol)
	if err != nil {
		return err
	}
	return u.client.ReplyMessage(replyToken, message)
}

func (u *lineCommandUsecase) GetStockNews(ctx context.Context, symbol string, replyToken string) error {
	newsMessage, err := u.botCommandUsecase.GetStockNewsForLine(ctx, symbol)
	if err != nil {
		return err
	}

	// 優先使用 Flex Message
	if newsMessage.UseFlexMessage && newsMessage.FlexContainer != nil {
		return u.client.ReplyFlexMessage(replyToken, newsMessage.Text, newsMessage.FlexContainer)
	}

	// 降級方案：Carousel Template
	if len(newsMessage.CarouselColumns) > 0 {
		return u.client.ReplyCarousel(replyToken, newsMessage.CarouselColumns)
	}

	// 最終降級：純文字
	return u.client.ReplyMessage(replyToken, newsMessage.Text)
}
