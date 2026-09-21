package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterWishAcceptanceRoutes 注册心愿验收路由。
func RegisterWishAcceptanceRoutes(rg *gin.RouterGroup, h *handler.WishAcceptanceHandler, auth gin.HandlerFunc) {
	rg.POST("/claims/:id/acceptance", auth, h.Submit)
	rg.GET("/wishes/:id/acceptance", auth, h.GetByWish)
	rg.POST("/acceptances/:id/confirm", auth, h.Confirm)
	rg.POST("/acceptances/:id/reject", auth, h.Reject)
}
