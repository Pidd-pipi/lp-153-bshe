package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/service"
)

// AuditLogHandler 审计日志处理器（仅管理员）。
type AuditLogHandler struct {
	audit service.AuditService
}

// NewAuditLogHandler 构造审计处理器。
func NewAuditLogHandler(audit service.AuditService) *AuditLogHandler {
	return &AuditLogHandler{audit: audit}
}

// List GET /api/v1/audit-logs
func (h *AuditLogHandler) List(c *gin.Context) {
	var q dto.AuditLogQuery
	if !bindQuery(c, &q) {
		return
	}
	result, err := h.audit.List(q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}
