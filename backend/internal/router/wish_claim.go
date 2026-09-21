package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterWishClaimRoutes 注册心愿认领路由。
func RegisterWishClaimRoutes(rg *gin.RouterGroup, h *handler.WishClaimHandler, auth gin.HandlerFunc) {
	rg.POST("/wishes/:id/claim", auth, h.Claim)
	rg.GET("/wishes/:id/claim", auth, h.GetByWish)
	rg.GET("/claims/mine", auth, h.Mine)
	rg.PUT("/claims/:id/progress", auth, h.UpdateProgress)
	rg.POST("/claims/:id/complete", auth, h.Complete)
}
