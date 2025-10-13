package linebot

import (
	linebotInfra "stock-bot/internal/infrastructure/linebot"
	"stock-bot/internal/service/bot/line"
	"stock-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

// LineBotHandler 處理 webhook 請求
type LineBotHandler struct {
	service   *line.LineBotHandler
	botClient *linebotInfra.LineBotClient
}

// NewLineBotHandler 創建 handler
func NewLineBotHandler(service *line.LineBotHandler, botClient *linebotInfra.LineBotClient) *LineBotHandler {
	return &LineBotHandler{
		service:   service,
		botClient: botClient,
	}
}

// Webhook 處理 LINE webhook 事件
func (h *LineBotHandler) Webhook(c *gin.Context) {
	events, err := h.botClient.Client.ParseRequest(c.Request)
	if err != nil {
		logger.Log.Error("Failed to parse webhook", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := h.service.HandleTextMessage(event, message); err != nil {
					logger.Log.Error("Failed to handle text message", zap.Error(err))
					// 即使處理失敗，也要回覆 LINE 平台，避免重複發送
					// 使用 ReplyMessage 回覆錯誤訊息給使用者
					if replyErr := h.botClient.ReplyMessage(event.ReplyToken, "處理訊息時發生錯誤，請稍後再試"); replyErr != nil {
						logger.Log.Error("Failed to send error reply", zap.Error(replyErr))
					}
				}
			}
		}
	}

	// 確保回傳 200 狀態碼，告訴 LINE 平台事件已成功處理
	c.JSON(200, gin.H{"status": "success"})
}
