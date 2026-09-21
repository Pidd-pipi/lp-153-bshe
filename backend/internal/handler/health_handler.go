package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/util"
)

// HealthHandler 健康检查处理器。
type HealthHandler struct{}

// NewHealthHandler 构造健康检查处理器。
func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// Healthz GET /healthz 返回服务健康状态。
func (h *HealthHandler) Healthz(c *gin.Context) {
	util.Logger().Info(constants.LogRequestCompleted, "path", "/healthz", "status", 200)
	c.JSON(200, gin.H{"status": "ok", "service": "wishwall", "time": util.FormatDateTime(now())})
}
