package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/util"
)

// handleError 统一把 handler/service 透传的错误转成标准 JSON 响应。
// handler 层再次包装 message，保证 message 中包含实体名。
func handleError(c *gin.Context, err error) {
	appErr, ok := util.AsAppError(err)
	if !ok {
		util.Logger().Error(constants.LogInternalError,
			"request_id", middleware.GetRequestID(c), "error", err)
		responseError(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		return
	}
	status := httpStatusOf(appErr.Code)
	util.Logger().Warn("handler wrapped service error",
		"request_id", middleware.GetRequestID(c), "code", appErr.Code, "message", appErr.Message)
	responseError(c, status, appErr.Code, appErr.Message)
}

func responseError(c *gin.Context, status, code int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"code": code, "message": message})
}

// httpStatusOf 错误码到 HTTP 状态映射。
func httpStatusOf(code int) int {
	switch {
	case code >= constants.CodeRateLimited:
		return http.StatusTooManyRequests
	case code >= constants.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	case code >= constants.CodeConflict:
		return http.StatusConflict
	case code >= constants.CodeNotFound:
		return http.StatusNotFound
	case code >= constants.CodeForbidden:
		return http.StatusForbidden
	case code >= constants.CodeUnauthorized:
		return http.StatusUnauthorized
	case code >= constants.CodeBadRequest:
		return http.StatusBadRequest
	}
	// 业务错误码（1xxxx-8xxxx）：按实体错误码映射 HTTP 状态。
	switch {
	case code == constants.CodeUserExists || code == constants.CodeWishAlreadyClaimed:
		return http.StatusConflict
	case code == constants.CodeUserBanned || code == constants.CodeWishNotOwner ||
		code == constants.CodeClaimNotOwner || code == constants.CodeAuditDenied:
		return http.StatusForbidden
	case code == constants.CodeInvalidCredential:
		return http.StatusUnauthorized
	case code == constants.CodeUserNotFound || code == constants.CodeWishNotFound ||
		code == constants.CodeClaimNotFound || code == constants.CodeCapsuleNotFound:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}

// bindJSON 绑定并校验 JSON 入参，失败统一返回 422。
func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		util.Logger().Warn(constants.LogInvalidParam, "request_id", middleware.GetRequestID(c), "error", err)
		responseError(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed,
			constants.MsgParamInvalid+"："+err.Error())
		return false
	}
	return true
}

// bindQuery 绑定查询参数并校验。
func bindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		util.Logger().Warn(constants.LogInvalidParam, "request_id", middleware.GetRequestID(c), "error", err)
		responseError(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed,
			constants.MsgParamInvalid+"："+err.Error())
		return false
	}
	return true
}
