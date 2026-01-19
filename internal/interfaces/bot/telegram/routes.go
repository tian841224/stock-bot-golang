package tgbot

import "github.com/gin-gonic/gin"

// RegisterRoutes 註冊 Telegram Bot 的路由
func RegisterRoutes(r *gin.Engine, handler *TgHandler, path string) {
	r.POST(path, handler.Webhook)
}
