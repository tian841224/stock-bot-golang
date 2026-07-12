package tgbot

import (
	"context"
	"net/http"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 部分指令會串連多個外部 API 呼叫（每個最長 10 秒）與圖表產生/上傳，
// 30 秒在慢速情況下偏緊，容易讓耗時較長的指令被中途取消卻無任何回覆，故拉寬至 60 秒。
const asyncProcessingTimeout = 60 * time.Second

type telegramUpdateProcessor interface {
	ProcessUpdate(ctx context.Context, update *tgbotapi.Update) error
}

type TgHandler struct {
	cfg              *config.Config
	messageProcessor telegramUpdateProcessor
	logger           logger.Logger
}

func NewTgHandler(
	cfg *config.Config,
	messageProcessor telegramUpdateProcessor,
	log logger.Logger,
) *TgHandler {
	return &TgHandler{
		cfg:              cfg,
		messageProcessor: messageProcessor,
		logger:           log,
	}
}

// Webhook 驗證 X-Telegram-Bot-Api-Secret-Token 並回應 200
func (h *TgHandler) Webhook(c *gin.Context) {
	if h.cfg.TELEGRAM_BOT_SECRET_TOKEN != "" {
		headerToken := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")
		if headerToken == "" || headerToken != h.cfg.TELEGRAM_BOT_SECRET_TOKEN {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	// 讀取並解析 update
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		// 從 context 取出已帶有 request_id 的 logger
		reqLogger := logger.FromContext(c.Request.Context(), h.logger)
		reqLogger.Error("failed to parse telegram webhook JSON", logger.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 先回應 200，背景處理，避免 Telegram 重送
	c.Status(http.StatusOK)

	// 複製 request context（包含 request_id logger）傳入 goroutine
	reqCtx := c.Request.Context()

	go func(u tgbotapi.Update) {
		processCtx, cancel := context.WithTimeout(logger.DetachContext(reqCtx, h.logger), asyncProcessingTimeout)
		defer cancel()

		// 取出 HTTP middleware 已注入的帶 request_id logger
		// 再加上 chat_id 方便過濾特定使用者的訊息
		var chatID int64
		if u.Message != nil {
			chatID = u.Message.Chat.ID
		}
		reqLogger := logger.FromContext(processCtx, h.logger).With(logger.Int64("chat_id", chatID))
		processCtx = logger.WithLogger(processCtx, reqLogger)

		defer func() {
			if r := recover(); r != nil {
				reqLogger.Error("panic recovering telegram update", logger.Any("recover", r))
			}
		}()

		if err := h.messageProcessor.ProcessUpdate(processCtx, &u); err != nil {
			reqLogger.Error("failed to process telegram update", logger.Error(err))
		}
	}(update)
}
