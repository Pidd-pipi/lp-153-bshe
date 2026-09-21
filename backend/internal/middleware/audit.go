package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/service"
)

// AuditMiddleware 通用审计中间件：为写操作自动落一条审计日志（service 层另有业务埋点）。
func AuditMiddleware(audit service.AuditService) gin.HandlerFunc {
	writeMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true, "DELETE": true}
	return func(c *gin.Context) {
		c.Next()
		if !writeMethods[c.Request.Method] {
			return
		}
		action := strings.ToLower(c.Request.Method) + "_" + entityFromPath(c.Request.URL.Path)
		_ = audit.Record(&model.AuditLog{
			UserID:     CurrentUserID(c),
			Action:     action,
			EntityType: entityFromPath(c.Request.URL.Path),
			EntityID:   idFromPath(c.Request.URL.Path),
			Detail:     "写操作审计：" + c.Request.Method + " " + c.Request.URL.Path,
			IP:         c.ClientIP(),
			RequestID:  GetRequestID(c),
		})
	}
}

func entityFromPath(path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	// /api/v1/<entity>/...
	for i, seg := range segments {
		if seg == "v1" && i+1 < len(segments) {
			return segments[i+1]
		}
	}
	return "unknown"
}

func idFromPath(path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	for i, seg := range segments {
		if seg == "v1" && i+2 < len(segments) {
			return segments[i+2]
		}
	}
	return ""
}
