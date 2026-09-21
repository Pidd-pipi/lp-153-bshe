package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// WishAcceptanceHandler 心愿验收处理器。
type WishAcceptanceHandler struct {
	acceptance service.WishAcceptanceService
}

// NewWishAcceptanceHandler 构造验收处理器。
func NewWishAcceptanceHandler(acceptance service.WishAcceptanceService) *WishAcceptanceHandler {
	return &WishAcceptanceHandler{acceptance: acceptance}
}

// Submit POST /api/v1/claims/:id/acceptance
func (h *WishAcceptanceHandler) Submit(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "认领 id 参数非法")
		return
	}
	var req dto.SubmitAcceptanceRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	acc, err := h.acceptance.Submit(c.Request.Context(), userID, claimID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgAcceptanceSubmitted, "data": dto.ToWishAcceptanceResponse(acc)})
}

// Confirm POST /api/v1/acceptances/:id/confirm
func (h *WishAcceptanceHandler) Confirm(c *gin.Context) {
	acceptanceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "验收 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	acc, err := h.acceptance.Confirm(c.Request.Context(), userID, acceptanceID, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgAcceptanceConfirmed, "data": dto.ToWishAcceptanceResponse(acc)})
}

// Reject POST /api/v1/acceptances/:id/reject
func (h *WishAcceptanceHandler) Reject(c *gin.Context) {
	acceptanceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "验收 id 参数非法")
		return
	}
	var req dto.RejectAcceptanceRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	acc, err := h.acceptance.Reject(c.Request.Context(), userID, acceptanceID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgAcceptanceRejected, "data": dto.ToWishAcceptanceResponse(acc)})
}

// GetByWish GET /api/v1/wishes/:id/acceptance
func (h *WishAcceptanceHandler) GetByWish(c *gin.Context) {
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	acc, err := h.acceptance.GetByWishID(userID, wishID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToWishAcceptanceResponse(acc)})
}
