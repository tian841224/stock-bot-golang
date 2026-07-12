package linebot

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/external/bot/line"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// 部分指令會串連多個外部 API 呼叫（每個最長 10 秒）與圖表產生/上傳，
// 30 秒在慢速情況下偏緊，容易讓耗時較長的指令被中途取消卻無任何回覆，故拉寬至 60 秒。
const asyncProcessingTimeout = 60 * time.Second

type lineMessageProcessor interface {
	ProcessTextMessage(ctx context.Context, event *linebot.Event, message *linebot.TextMessage) error
}

// LineBotHandler 處理 webhook 請求
type LineBotHandler struct {
	botClient        *linebotInfra.LineBotClient
	messageProcessor lineMessageProcessor
	logger           logger.Logger
}

// NewLineBotHandler 創建 handler
func NewLineBotHandler(
	botClient *linebotInfra.LineBotClient,
	messageProcessor lineMessageProcessor,
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
		reqLogger := logger.FromContext(c.Request.Context(), h.logger)
		reqLogger.Error("failed to parse LINE webhook request", logger.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 先回應 200，背景處理，避免 LINE 平台重送
	c.Status(http.StatusOK)

	// 複製 request context（包含 request_id logger）傳入 goroutine
	reqCtx := c.Request.Context()

	// 在 goroutine 中處理事件，避免 webhook 超時
	go func(evts []*linebot.Event) {
		processCtx, cancel := context.WithTimeout(logger.DetachContext(reqCtx, h.logger), asyncProcessingTimeout)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				logger.FromContext(processCtx, h.logger).Error("panic recovering LINE update", logger.Any("recover", r))
			}
		}()

		for _, event := range evts {
			if event.Type == linebot.EventTypeMessage {
				// 從 HTTP middleware 取出帶 request_id 的 logger，再加上 reply_token
				reqLogger := logger.FromContext(processCtx, h.logger).With(
					logger.String("reply_token", event.ReplyToken),
				)
				eventCtx := logger.WithLogger(processCtx, reqLogger)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if err := h.messageProcessor.ProcessTextMessage(eventCtx, event, message); err != nil {
						reqLogger.Error("failed to process LINE text message", logger.Error(err))
					}
				}
			}
		}
	}(events)
}
