package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/database"
	"github.com/wishwall/wishwall/internal/handler"
	"github.com/wishwall/wishwall/internal/middleware"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/router"
	"github.com/wishwall/wishwall/internal/service"
	"github.com/wishwall/wishwall/internal/util"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	logger := util.InitLogger(cfg.AppEnv)

	db, err := database.Open(cfg, logger)
	if err != nil {
		logger.Error(constants.LogDBMigrateFailed, "error", err)
		os.Exit(1)
	}
	rdb, err := database.OpenRedis(cfg, logger)
	if err != nil {
		logger.Error("redis connect failed", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	minioClient, err := database.OpenMinIO(cfg, logger)
	if err != nil {
		logger.Error("minio connect failed", "error", err)
		os.Exit(1)
	}

	// 仓储层
	userRepo := repository.NewUserRepository(db)
	wishRepo := repository.NewWishRepository(db)
	claimRepo := repository.NewWishClaimRepository(db)
	acceptanceRepo := repository.NewWishAcceptanceRepository(db)
	blessRepo := repository.NewBlessingRepository(db)
	capsuleRepo := repository.NewTimeCapsuleRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	txManager := repository.NewTxManager(db)

	// 服务层（构造器注入，依赖单向）
	auditSvc := service.NewAuditService(auditRepo, userRepo, logger)
	badgeSvc := service.NewBadgeService(badgeRepo, wishRepo, rdb, logger)
	userSvc := service.NewUserService(userRepo, cfg, auditSvc, logger)
	wishSvc := service.NewWishService(wishRepo, claimRepo, acceptanceRepo, blessRepo, userRepo, badgeSvc, auditSvc, logger)
	claimSvc := service.NewWishClaimService(txManager, wishRepo, claimRepo, acceptanceRepo, userRepo, badgeSvc, auditSvc, logger)
	acceptanceSvc := service.NewWishAcceptanceService(txManager, wishRepo, claimRepo, acceptanceRepo, badgeSvc, auditSvc, logger)
	blessSvc := service.NewBlessingService(blessRepo, wishRepo, userRepo, badgeSvc, auditSvc, logger)
	capsuleSvc := service.NewTimeCapsuleService(capsuleRepo, auditSvc, logger)
	uploadSvc := service.NewUploadService(cfg, minioClient, auditSvc, logger)

	// 处理器
	handlers := &router.Handlers{
		User:       handler.NewUserHandler(userSvc),
		Wish:       handler.NewWishHandler(wishSvc),
		Claim:      handler.NewWishClaimHandler(claimSvc),
		Acceptance: handler.NewWishAcceptanceHandler(acceptanceSvc),
		Bless:      handler.NewBlessingHandler(blessSvc),
		Capsule:    handler.NewTimeCapsuleHandler(capsuleSvc),
		Badge:      handler.NewBadgeHandler(badgeSvc),
		Audit:      handler.NewAuditLogHandler(auditSvc),
		Upload:     handler.NewUploadHandler(uploadSvc),
		Health:     handler.NewHealthHandler(),
	}

	// 默认管理员种子
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@123456"
	}
	if err := userSvc.SeedAdmin("admin", "admin@wishwall.local", adminPassword); err != nil {
		logger.Warn("seed admin failed", "error", err)
	}

	// 时光胶囊到期解锁定时任务
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go runCapsuleUnlocker(ctx, capsuleSvc, logger)

	gin.SetMode(gin.ReleaseMode)
	engine := router.NewRouter(cfg, handlers, middleware.AuditMiddleware(auditSvc))
	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		logger.Info(constants.LogSrvStarted, "port", cfg.ServerPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info(constants.LogSrvShutdown)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func runCapsuleUnlocker(ctx context.Context, capsuleSvc service.TimeCapsuleService, logger *slog.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			unlocked, err := capsuleSvc.UnlockDue()
			if err != nil {
				logger.Warn("capsule unlock tick failed", "error", err)
				continue
			}
			if len(unlocked) > 0 {
				logger.Info(constants.LogCapsuleUnlocked, "count", len(unlocked))
			}
		case <-ctx.Done():
			return
		}
	}
}
