package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/util"
)

// RequireRole RBAC 权限中间件：仅允许指定角色访问。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role := CurrentRole(c)
		if !allowed[role] {
			util.Logger().Warn(constants.LogForbidden, "request_id", GetRequestID(c), "role", role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    constants.CodeForbidden,
				"message": constants.MsgNeedAdmin + "：" + constants.MsgForbiddenRole,
			})
			return
		}
		c.Next()
	}
}
