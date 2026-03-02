package linebot

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/tian841224/stock-bot/internal/application/usecase/bot"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/line"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// LineBotHandler 處理 webhook 請求
type LineBotHandler struct {
	botClient        *linebotInfra.LineBotClient
	messageProcessor *bot.LineMessageProcessor
	logger           logger.Logger
}

// NewLineBotHandler 創建 handler
func NewLineBotHandler(
	botClient *linebotInfra.LineBotClient,
	messageProcessor *bot.LineMessageProcessor,
	log logger.Logger,
) *LineBotHandler {
	return &LineBotHandler{
		botClient:        botClient,
		messageProcessor: messageProcessor,
		logger:           log,
	}
}

// Webhook 處理 LINE webhook 事件
func (h *LineBotHandler) Webhook(c *gin.Context) {
	events, err := h.botClient.Client.ParseRequest(c.Request)
	if err != nil {
		h.logger.Error("failed to parse LINE webhook request", logger.Error(err))
		c.AbortWithStatus(400)
		return
	}

	// 先回應 200，背景處理，避免 LINE 平台重送
	c.Status(200)

	// 在 goroutine 中處理事件，避免 webhook 超時
	go func(evts []*linebot.Event) {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("panic recovering LINE update", logger.Any("recover", r))
			}
		}()

		ctx := context.Background()
		for _, event := range evts {
			if event.Type == linebot.EventTypeMessage {
				// 建立帶有 request_id 的子 logger，讓整次請求的 log 可串聯
				reqLogger := h.logger.With(
					logger.String("request_id", "line-"+event.ReplyToken),
					logger.String("reply_token", event.ReplyToken),
				)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if err := h.messageProcessor.ProcessTextMessage(ctx, event, message); err != nil {
						reqLogger.Error("failed to process LINE text message", logger.Error(err))
					}
				}
			}
		}
	}(events)
}
