package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/util"
)

// Auth JWT 认证中间件：解析 Bearer Token 并把 user_id/username/role 写入上下文。
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		}
		if token == "" {
			abortAuth(c, constants.CodeUnauthorized, constants.MsgNeedLogin)
			return
		}
		claims, err := util.ParseToken(jwtSecret, token)
		if err != nil {
			util.Logger().Warn(constants.LogUnauthorized, "request_id", GetRequestID(c), "error", err)
			abortAuth(c, constants.CodeUnauthorized, constants.MsgNeedLogin+"：登录已过期")
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func abortAuth(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": code, "message": message})
}

// CurrentUserID 从上下文取当前用户 ID。
func CurrentUserID(c *gin.Context) uint64 {
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(uint64); ok {
			return id
		}
	}
	return 0
}

// CurrentRole 从上下文取当前用户角色。
func CurrentRole(c *gin.Context) string {
	if v, ok := c.Get("role"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
