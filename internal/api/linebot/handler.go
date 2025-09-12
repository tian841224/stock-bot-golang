package linebot

import (
	linebotInfra "stock-bot/internal/infrastructure/linebot"
	"stock-bot/internal/service/bot/line"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

// LineBotHandler 處理 webhook 請求
type LineBotHandler struct {
	service   *line.LineBotService
	botClient *linebotInfra.LineBotClient
}

// NewLineBotHandler 創建 handler
func NewLineBotHandler(service *line.LineBotService, botClient *linebotInfra.LineBotClient) *LineBotHandler {
	return &LineBotHandler{
		service:   service,
		botClient: botClient,
	}	
}

// Webhook 處理 LINE webhook 事件
func (h *LineBotHandler) Webhook(c *gin.Context) {
	events, err := h.botClient.Client.ParseRequest(c.Request)
	if err != nil {
		//logger.Log.Error("Failed to parse webhook")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := h.service.HandleTextMessage(event, message); err != nil {
					//logger.Log.Error("Failed to handle text message")
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}
	c.Status(200)
}
