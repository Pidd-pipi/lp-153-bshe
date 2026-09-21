package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/handler"
	"github.com/wishwall/wishwall/internal/middleware"
)

// Handlers 全部处理器聚合，供路由注册使用。
type Handlers struct {
	User       *handler.UserHandler
	Wish       *handler.WishHandler
	Claim      *handler.WishClaimHandler
	Acceptance *handler.WishAcceptanceHandler
	Bless      *handler.BlessingHandler
	Capsule    *handler.TimeCapsuleHandler
	Badge      *handler.BadgeHandler
	Audit      *handler.AuditLogHandler
	Upload     *handler.UploadHandler
	Health     *handler.HealthHandler
}

// NewRouter 装配 Gin 引擎：全局中间件 + /healthz + /api/v1 分组路由。
func NewRouter(cfg *config.Config, hs *Handlers, auditMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimit(cfg.RateLimitPerMin))

	auth := middleware.Auth(cfg.JWTSecret)
	admin := middleware.RequireRole(constants.RoleAdmin)

	r.GET("/healthz", hs.Health.Healthz)

	api := r.Group("/api/v1")
	if auditMiddleware != nil {
		api.Use(auditMiddleware)
	}
	RegisterUserRoutes(api, hs.User, auth)
	RegisterWishRoutes(api, hs.Wish, auth)
	RegisterWishClaimRoutes(api, hs.Claim, auth)
	RegisterWishAcceptanceRoutes(api, hs.Acceptance, auth)
	RegisterBlessingRoutes(api, hs.Bless, auth)
	RegisterTimeCapsuleRoutes(api, hs.Capsule, auth)
	RegisterBadgeRoutes(api, hs.Badge, auth)
	RegisterAuditLogRoutes(api, hs.Audit, auth, admin)
	RegisterUploadRoutes(api, hs.Upload, auth)

	return r
}
