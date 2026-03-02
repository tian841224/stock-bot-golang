package tgbot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tian841224/stock-bot/internal/application/usecase/bot"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	cfg              *config.Config
	messageProcessor *bot.TelegramMessageProcessor
	logger           logger.Logger
}

func NewTgHandler(
	cfg *config.Config,
	messageProcessor *bot.TelegramMessageProcessor,
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
		h.logger.Error("failed to parse telegram webhook JSON", logger.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 先回應 200，背景處理，避免 Telegram 重送
	c.Status(http.StatusOK)

	go func(u tgbotapi.Update) {
		// 建立帶有 request_id 的子 logger，讓整次請求的 log 可串聯
		var chatID int64
		if u.Message != nil {
			chatID = u.Message.Chat.ID
		}
		requestID := fmt.Sprintf("tg-%d-%d", chatID, c.Request.Context().Value("request_nano"))
		reqLogger := h.logger.With(logger.String("request_id", requestID), logger.Int64("chat_id", chatID))

		defer func() {
			if r := recover(); r != nil {
				reqLogger.Error("panic recovering telegram update", logger.Any("recover", r))
			}
		}()

		ctx := context.Background()
		if err := h.messageProcessor.ProcessUpdate(ctx, &u); err != nil {
			reqLogger.Error("failed to process telegram update", logger.Error(err))
		}
	}(update)
}
