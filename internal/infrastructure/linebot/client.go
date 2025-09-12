package linebot

import (
	"stock-bot/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBotClient struct {
	Client *linebot.Client
}

// NewBot 初始化 LINE Bot
func NewBot(cfg config.Config) (*LineBotClient, error) {
	client, err := linebot.New(cfg.CHANNEL_SECRET, cfg.CHANNEL_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	return &LineBotClient{Client: client}, nil
}

// ReplyMessage 回覆文字訊息
func (b *LineBotClient) ReplyMessage(replyToken, text string) error {
	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do()
	return err
}
