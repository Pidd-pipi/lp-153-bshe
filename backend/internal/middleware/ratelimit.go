package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/util"
)

type windowEntry struct {
	count   int
	resetAt time.Time
}

// RateLimit 简单内存滑动窗口限流：按 IP 每分钟限流。
func RateLimit(perMin int) gin.HandlerFunc {
	var mu sync.Mutex
	windows := map[string]*windowEntry{}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		entry, ok := windows[ip]
		if !ok || now.After(entry.resetAt) {
			entry = &windowEntry{count: 0, resetAt: now.Add(time.Minute)}
			windows[ip] = entry
		}
		entry.count++
		if entry.count > perMin {
			mu.Unlock()
			util.Logger().Warn(constants.LogRateLimited, "request_id", GetRequestID(c), "ip", ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    constants.CodeRateLimited,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}
		mu.Unlock()
		c.Next()
	}
}
