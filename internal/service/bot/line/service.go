package line

import (
	"strings"

	linebotInfra "stock-bot/internal/infrastructure/linebot"

	"github.com/line/line-bot-sdk-go/linebot"
)

// BotService 處理對話邏輯
type LineBotService struct {
	botClient *linebotInfra.LineBotClient
}

// NewBotService 創建 service
func NewBotService(botClient *linebotInfra.LineBotClient) *LineBotService {
	return &LineBotService{
		botClient: botClient,
	}
}

// HandleTextMessage 處理文字訊息
func (s *LineBotService) HandleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if message.Text == "" {
		return nil
	}

	parts := strings.Split(strings.TrimSpace(message.Text), " ")
	command := parts[0]

	switch command {
	case "test":
		return s.botClient.ReplyMessage(event.ReplyToken, "新增成功")
	default:
		reply := "你說了: " + message.Text
		return s.botClient.ReplyMessage(event.ReplyToken, reply)
	}
}
