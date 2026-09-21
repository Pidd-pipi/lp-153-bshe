package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/util"
)

// Recovery 捕获 panic，记录 request_id 与堆栈，返回统一 500 JSON。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				util.Logger().Error(constants.LogInternalError,
					"request_id", GetRequestID(c),
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    constants.CodeInternalError,
					"message": constants.MsgInternalError,
				})
			}
		}()
		c.Next()
	}
}
