package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/service"
	"github.com/wishwall/wishwall/internal/util"
)

// UserHandler 用户相关处理器。
type UserHandler struct {
	user service.UserService
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(user service.UserService) *UserHandler {
	return &UserHandler{user: user}
}

// Register POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.user.Register(c.Request.Context(), req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	util.Logger().Info(constants.LogUserRegistered, "user_id", user.ID)
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgRegisterSuccess, "data": dto.ToUserResponse(user)})
}

// Login POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	resp, err := h.user.Login(c.Request.Context(), req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	util.Logger().Info(constants.LogUserLoggedIn, "user_id", resp.User.ID)
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgLoginSuccess, "data": resp})
}

// Me GET /api/v1/users/me
func (h *UserHandler) Me(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	user, err := h.user.GetByID(userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToUserResponse(user)})
}

// GetByID GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responseError(c, 400, constants.CodeBadRequest, "用户 id 参数非法")
		return
	}
	user, err := h.user.GetByID(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": dto.ToUserResponse(user)})
}

// UpdateProfile PUT /api/v1/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)
	user, err := h.user.UpdateProfile(c.Request.Context(), userID, req, c.ClientIP(), middleware.GetRequestID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	util.Logger().Info(constants.LogUserUpdated, "user_id", user.ID)
	c.JSON(200, gin.H{"code": 0, "message": constants.MsgUpdateSuccess, "data": dto.ToUserResponse(user)})
}
