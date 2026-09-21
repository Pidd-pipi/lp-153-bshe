package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/handler"
)

// RegisterTimeCapsuleRoutes 注册时光胶囊路由。
func RegisterTimeCapsuleRoutes(rg *gin.RouterGroup, h *handler.TimeCapsuleHandler, auth gin.HandlerFunc) {
	rg.POST("/capsules", auth, h.Create)
	rg.GET("/capsules/mine", auth, h.ListMine)
	rg.GET("/capsules/:id", auth, h.GetByID)
	rg.DELETE("/capsules/:id", auth, h.Delete)
}
