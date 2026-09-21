package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// BlessingHandler 祝福留言处理器。
type BlessingHandler struct {
	blessing service.BlessingService
}

// NewBlessingHandler 构造祝福处理器。
func NewBlessingHandler(blessing service.BlessingService) *BlessingHandler {
	return &BlessingHandler{blessing: blessing}
}

// Create POST /api/v1/wishes/:id/blessings
func (h *BlessingHandler) Create(c *gin.Context) {
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	var req dto.CreateBlessingRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	blessing, err := h.blessing.Create(userID, wishID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToBlessingResponse(blessing, "")})
}

// List GET /api/v1/wishes/:id/blessings
func (h *BlessingHandler) List(c *gin.Context) {
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	var q dto.PageQuery
	if !bindQuery(c, &q) {
		return
	}
	result, err := h.blessing.ListByWish(wishID, q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}
