package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterUploadRoutes 注册文件上传路由。
func RegisterUploadRoutes(rg *gin.RouterGroup, h *handler.UploadHandler, auth gin.HandlerFunc) {
	rg.POST("/uploads", auth, h.Upload)
}
