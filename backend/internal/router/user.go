package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterUserRoutes 注册用户/认证路由。
func RegisterUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler, auth gin.HandlerFunc) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)

	users := rg.Group("/users")
	users.GET("/:id", h.GetByID)

	me := rg.Group("/users/me")
	me.Use(auth)
	me.GET("", h.Me)
	me.PUT("", h.UpdateProfile)
}
