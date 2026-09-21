package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterWishRoutes 注册心愿路由。
func RegisterWishRoutes(rg *gin.RouterGroup, h *handler.WishHandler, auth gin.HandlerFunc) {
	rg.POST("/wishes", auth, h.Create)
	rg.GET("/wishes", h.List)
	rg.GET("/wishes/mine", auth, h.Mine)
	rg.GET("/wishes/:id", h.GetByID)
	rg.PUT("/wishes/:id", auth, h.Update)
	rg.DELETE("/wishes/:id", auth, h.Delete)
	rg.POST("/wishes/:id/like", auth, h.Like)
	rg.GET("/discover", h.Discover)
	rg.GET("/discover/leaderboard", h.Leaderboard)
}
