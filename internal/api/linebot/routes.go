package linebot

import "github.com/gin-gonic/gin"

// RegisterRoutes 註冊LINE Bot的路由
func RegisterRoutes(r *gin.Engine, handler *LineBotHandler) {
	r.POST("/webhook", handler.Webhook)
}
