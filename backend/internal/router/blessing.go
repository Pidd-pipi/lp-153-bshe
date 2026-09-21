package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterBlessingRoutes 注册祝福留言路由。
func RegisterBlessingRoutes(rg *gin.RouterGroup, h *handler.BlessingHandler, auth gin.HandlerFunc) {
	rg.POST("/wishes/:id/blessings", auth, h.Create)
	rg.GET("/wishes/:id/blessings", h.List)
}
