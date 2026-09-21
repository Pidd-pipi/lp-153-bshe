package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// BadgeHandler 成就徽章处理器。
type BadgeHandler struct {
	badge service.BadgeService
}

// NewBadgeHandler 构造徽章处理器。
func NewBadgeHandler(badge service.BadgeService) *BadgeHandler {
	return &BadgeHandler{badge: badge}
}

// Mine GET /api/v1/badges/mine
func (h *BadgeHandler) Mine(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	items, err := h.badge.ListMine(userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": items})
}

// Leaderboard GET /api/v1/badges/leaderboard（复用 BadgeService.Leaderboard）
func (h *BadgeHandler) Leaderboard(c *gin.Context) {
	rows, err := h.badge.Leaderboard(10)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": rows})
}
