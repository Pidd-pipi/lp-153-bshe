package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// UploadHandler 文件上传处理器。
type UploadHandler struct {
	upload service.UploadService
}

// NewUploadHandler 构造上传处理器。
func NewUploadHandler(upload service.UploadService) *UploadHandler {
	return &UploadHandler{upload: upload}
}

// Upload POST /api/v1/uploads?kind=image|audio
func (h *UploadHandler) Upload(c *gin.Context) {
	kind := c.DefaultQuery("kind", "image")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "缺少 file 字段")
		return
	}
	defer file.Close()
	userID := middleware.CurrentUserID(c)
	url, err := h.upload.Upload(c.Request.Context(), userID, file, header, kind, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgUploadSuccess, "data": gin.H{"url": url}})
}
