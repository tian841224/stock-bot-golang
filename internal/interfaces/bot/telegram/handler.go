package tgbot

import (
	"context"
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

	// 讀取並解析 update（可視需求擴充）
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		h.logger.Error("解析 Telegram webhook JSON 失敗", logger.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 先回應 200，背景處理，避免 Telegram 重送
	c.Status(http.StatusOK)

	go func(u tgbotapi.Update) {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("處理 Telegram 更新發生 panic", logger.Any("recover", r))
			}
		}()

		ctx := context.Background()
		if err := h.messageProcessor.ProcessUpdate(ctx, &u); err != nil {
			h.logger.Error("處理 Telegram 更新失敗", logger.Error(err))
		}
	}(update)
}
