package tgbot

import (
	"net/http"

	"github.com/tian841224/stock-bot/config"

	"github.com/tian841224/stock-bot/internal/service/bot/tg"
	"github.com/tian841224/stock-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgHandler struct {
	cfg              *config.Config
	tgServiceHandler *tg.TgHandler
}

func NewTgHandler(cfg *config.Config, tgServiceHandler *tg.TgHandler) *TgHandler {
	return &TgHandler{cfg: cfg, tgServiceHandler: tgServiceHandler}
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 先回應 200，背景處理，避免 Telegram 重送
	c.Status(http.StatusOK)

	go func(u tgbotapi.Update) {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("處理 Telegram 更新發生 panic", zap.Any("recover", r))
			}
		}()
		if err := h.tgServiceHandler.ProcessUpdate(&u); err != nil {
			logger.Log.Error("處理 Telegram 更新失敗", zap.Error(err), zap.Int("update_id", u.UpdateID))
		}
	}(update)
}
