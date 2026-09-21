package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
)

// TimeCapsuleHandler 时光胶囊处理器。
type TimeCapsuleHandler struct {
	capsule service.TimeCapsuleService
}

// NewTimeCapsuleHandler 构造胶囊处理器。
func NewTimeCapsuleHandler(capsule service.TimeCapsuleService) *TimeCapsuleHandler {
	return &TimeCapsuleHandler{capsule: capsule}
}

// Create POST /api/v1/capsules
func (h *TimeCapsuleHandler) Create(c *gin.Context) {
	var req dto.CreateCapsuleRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	capsule, err := h.capsule.Create(userID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgCapsuleCreated, "data": dto.ToCapsuleResponse(capsule, true)})
}

// ListMine GET /api/v1/capsules/mine
func (h *TimeCapsuleHandler) ListMine(c *gin.Context) {
	var q dto.PageQuery
	if !bindQuery(c, &q) {
		return
	}
	userID := middleware.CurrentUserID(c)
	result, err := h.capsule.ListMine(userID, q)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": result})
}

// GetByID GET /api/v1/capsules/:id
func (h *TimeCapsuleHandler) GetByID(c *gin.Context) {
	capsuleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "胶囊 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	capsule, err := h.capsule.GetByID(userID, capsuleID)
	if err != nil {
		handleError(c, err)
		return
	}
	mask := capsule.Status == constants.CapsuleStatusLocked
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToCapsuleResponse(capsule, mask)})
}

// Delete DELETE /api/v1/capsules/:id
func (h *TimeCapsuleHandler) Delete(c *gin.Context) {
	capsuleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "胶囊 id 参数非法")
		return
	}
	userID := middleware.CurrentUserID(c)
	if err := h.capsule.Delete(userID, capsuleID, c.ClientIP(), middleware.GetRequestID(c)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgDeleteSuccess, "data": nil})
}
