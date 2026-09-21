package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/constants"
)

// Body 统一响应包裹：{ "code": 0, "message": "ok", "data": ... }。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 成功响应。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: constants.CodeOK, Message: constants.MsgOK, Data: data})
}

// Error 错误响应。
func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message})
}
