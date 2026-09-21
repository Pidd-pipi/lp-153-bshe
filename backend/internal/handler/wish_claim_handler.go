package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// WishClaimHandler 心愿认领处理器。
type WishClaimHandler struct {
	claim service.WishClaimService
}

// NewWishClaimHandler 构造认领处理器。
func NewWishClaimHandler(claim service.WishClaimService) *WishClaimHandler {
	return &WishClaimHandler{claim: claim}
}

// Claim POST /api/v1/wishes/:id/claim
func (h *WishClaimHandler) Claim(c *gin.Context) {
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	claim, err := h.claim.Claim(c.Request.Context(), userID, wishID, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgWishClaimed, "data": dto.ToWishClaimResponse(claim, "", "")})
}

// UpdateProgress PUT /api/v1/claims/:id/progress
func (h *WishClaimHandler) UpdateProgress(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "认领 id 参数非法")
		return
	}
	var req dto.UpdateProgressRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	claim, err := h.claim.UpdateProgress(c.Request.Context(), userID, claimID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgUpdateSuccess, "data": dto.ToWishClaimResponse(claim, "", "")})
}

// Submit POST /api/v1/claims/:id/submit 圆梦人送交验收。
func (h *WishClaimHandler) Submit(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "认领 id 参数非法")
		return
	}
	var req dto.SubmitClaimRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	claim, err := h.claim.SubmitForAcceptance(c.Request.Context(), userID, claimID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgClaimSubmitted, "data": dto.ToWishClaimResponse(claim, "", "")})
}

// Review POST /api/v1/claims/:id/review 发布者验收（确认/驳回）。
func (h *WishClaimHandler) Review(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "认领 id 参数非法")
		return
	}
	var req dto.ReviewClaimRequest
	if !bindJSON(c, &req) {
		return
	}
	// handler 层再次包装：驳回必须写原因（错误信息携带实体名 wish_claim 与字段名 reason）。
	if req.Action == "reject" && req.Reason == "" {
		responseError(c, 422, constants.CodeRejectReasonRequired, constants.MsgRejectReason+"（实体：wish_claim，字段：reason，角色：发布者）")
		return
	}
	userID := middleware.CurrentUserID(c)
	claim, err := h.claim.Review(c.Request.Context(), userID, claimID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	message := constants.MsgClaimAccepted
	if req.Action == "reject" {
		message = constants.MsgClaimRejected
	}
	c.JSON(200, gin.H{"code": 0, "message": message, "data": dto.ToWishClaimResponse(claim, "", "")})
}

// Mine GET /api/v1/claims/mine
func (h *WishClaimHandler) Mine(c *gin.Context) {
	var q dto.PageQuery
	if !bindQuery(c, &q) {
		return
	}
	userID := middleware.CurrentUserID(c)
	result, err := h.claim.ListMine(userID, q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}

// GetByWish GET /api/v1/wishes/:id/claim
func (h *WishClaimHandler) GetByWish(c *gin.Context) {
	wishID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "心愿 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	claim, err := h.claim.GetByWishID(userID, wishID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToWishClaimResponse(claim, "", "")})
}
