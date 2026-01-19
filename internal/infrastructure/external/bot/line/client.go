package linebot

import (
	"bytes"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/external/imgbb"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type LineBotClient struct {
	Client *linebot.Client
	logger logger.Logger
}

// NewBot 初始化 LINE Bot
func NewBot(cfg config.Config, log logger.Logger) (*LineBotClient, error) {
	client, err := linebot.New(cfg.CHANNEL_SECRET, cfg.CHANNEL_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	return &LineBotClient{Client: client, logger: log}, nil
}

// ReplyMessage 回覆文字訊息
func (b *LineBotClient) ReplyMessage(replyToken, text string) error {
	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		b.logger.Error("發送訊息失敗", logger.Error(err))
	}
	return err
}

// ReplyMessageWithButtons 回覆帶有按鈕的訊息
func (b *LineBotClient) ReplyMessageWithButtons(replyToken, text string, buttons []linebot.TemplateAction) error {
	if len(buttons) == 0 {
		return b.ReplyMessage(replyToken, text)
	}

	template := linebot.NewButtonsTemplate(
		"", "", text, buttons...,
	)

	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTemplateMessage("按鈕", template)).Do()
	if err != nil {
		b.logger.Error("發送帶有按鈕的訊息失敗", logger.Error(err))
	}
	return err
}

// ReplyImage 回覆圖片訊息
func (b *LineBotClient) ReplyImage(replyToken, imageURL string) error {
	imageMessage := linebot.NewImageMessage(imageURL, imageURL)
	_, err := b.Client.ReplyMessage(replyToken, imageMessage).Do()
	if err != nil {
		b.logger.Error("發送圖片訊息失敗", logger.Error(err))
	}
	return err
}

// ReplyPhoto 上傳圖片並回覆（需要 ImgBB 客戶端）
func (b *LineBotClient) ReplyPhoto(replyToken string, data []byte, caption string, imgbbClient *imgbb.ImgBBClient) error {
	// 如果沒有 ImgBB 客戶端，只發送文字訊息
	if imgbbClient == nil {
		b.logger.Warn("ImgBB 客戶端未設定，只發送文字訊息")
		return b.ReplyMessage(replyToken, caption)
	}

	// 上傳圖片到 ImgBB
	options := &imgbb.UploadOptions{
		Name: "stock_chart",
	}

	reader := bytes.NewReader(data)
	resp, err := imgbbClient.UploadFromFile(reader, "chart.png", options)
	if err != nil {
		b.logger.Error("上傳圖片到 ImgBB 失敗", logger.Error(err))
		// 如果上傳失敗，只發送文字訊息
		return b.ReplyMessage(replyToken, caption)
	}

	// 發送圖片
	return b.ReplyImage(replyToken, resp.Data.URL)
}

// ReplyCarousel 回覆輪播模板訊息
func (b *LineBotClient) ReplyCarousel(replyToken string, columns []*linebot.CarouselColumn) error {
	if len(columns) == 0 {
		return nil
	}

	template := linebot.NewCarouselTemplate(columns...)
	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTemplateMessage("新聞列表", template)).Do()
	if err != nil {
		b.logger.Error("發送輪播訊息失敗", logger.Error(err))
	}
	return err
}

// ReplyFlexMessage 回覆 Flex Message
func (b *LineBotClient) ReplyFlexMessage(replyToken string, altText string, contents linebot.FlexContainer) error {
	_, err := b.Client.ReplyMessage(replyToken, linebot.NewFlexMessage(altText, contents)).Do()
	if err != nil {
		b.logger.Error("發送 Flex 訊息失敗", logger.Error(err))
	}
	return err
}
