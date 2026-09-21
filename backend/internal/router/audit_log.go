package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterAuditLogRoutes 注册审计日志路由（仅管理员）。
func RegisterAuditLogRoutes(rg *gin.RouterGroup, h *handler.AuditLogHandler, auth, admin gin.HandlerFunc) {
	rg.GET("/audit-logs", auth, admin, h.List)
}
