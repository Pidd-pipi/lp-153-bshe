package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterBadgeRoutes 注册成就徽章路由。
func RegisterBadgeRoutes(rg *gin.RouterGroup, h *handler.BadgeHandler, auth gin.HandlerFunc) {
	rg.GET("/badges/mine", auth, h.Mine)
	rg.GET("/badges/leaderboard", h.Leaderboard)
}
