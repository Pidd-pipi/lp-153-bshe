package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// WishHandler 心愿相关处理器。
type WishHandler struct {
	wish service.WishService
}

// NewWishHandler 构造心愿处理器。
func NewWishHandler(wish service.WishService) *WishHandler {
	return &WishHandler{wish: wish}
}

// Create POST /api/v1/wishes
func (h *WishHandler) Create(c *gin.Context) {
	var req dto.CreateWishRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	wish, err := h.wish.Create(userID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgWishCreated, "data": dto.ToWishResponse(wish)})
}

// List GET /api/v1/wishes
func (h *WishHandler) List(c *gin.Context) {
	var q dto.WishQuery
	if !bindQuery(c, &q) {
		return
	}
	result, err := h.wish.List(q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}

// GetByID GET /api/v1/wishes/:id
func (h *WishHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	detail, err := h.wish.GetByID(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": detail})
}

// Update PUT /api/v1/wishes/:id
func (h *WishHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	var req dto.UpdateWishRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	wish, err := h.wish.Update(userID, id, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgUpdateSuccess, "data": dto.ToWishResponse(wish)})
}

// Delete DELETE /api/v1/wishes/:id
func (h *WishHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	if err := h.wish.Delete(userID, id, c.ClientIP(), middleware.GetRequestID(c)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgDeleteSuccess, "data": nil})
}

// Mine GET /api/v1/wishes/mine
func (h *WishHandler) Mine(c *gin.Context) {
	var q dto.PageQuery
	if !bindQuery(c, &q) {
		return
	}
	userID := middleware.CurrentUserID(c)
	result, err := h.wish.ListMine(userID, q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}

// Like POST /api/v1/wishes/:id/like
func (h *WishHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	if err := h.wish.Like(userID, id, c.ClientIP(), middleware.GetRequestID(c)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": gin.H{"wish_id": id}})
}

// Discover GET /api/v1/discover
func (h *WishHandler) Discover(c *gin.Context) {
	var q dto.PageQuery
	if !bindQuery(c, &q) {
		return
	}
	stories, err := h.wish.Discover(q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": stories})
}

// Leaderboard GET /api/v1/discover/leaderboard
func (h *WishHandler) Leaderboard(c *gin.Context) {
	rows, err := h.wish.Leaderboard(10)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": rows})
}
