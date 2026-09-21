package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// UserService 用户服务：注册/登录/资料维护。
type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest, ip, requestID string) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest, ip, requestID string) (*dto.LoginResponse, error)
	GetByID(userID uint64) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest, ip, requestID string) (*model.User, error)
	SeedAdmin(username, email, password string) error
}

type userService struct {
	repo  repository.UserRepository
	cfg   *config.Config
	audit AuditService
	logger *slog.Logger
}

// NewUserService 构造用户服务。
func NewUserService(repo repository.UserRepository, cfg *config.Config, audit AuditService, logger *slog.Logger) UserService {
	return &userService{repo: repo, cfg: cfg, audit: audit, logger: logger}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest, ip, requestID string) (*model.User, error) {
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, constants.MsgUsernameExists, errors.New("username taken"))
	}
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, constants.MsgEmailExists, errors.New("email taken"))
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Nickname:     nickname,
		Role:         constants.RoleUser,
		Status:       "active",
	}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(constants.CodeUserExists, constants.MsgEmailExists, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserRegistered, "user_id", user.ID, "username", user.Username)
	_ = s.audit.Record(&model.AuditLog{
		UserID: user.ID, Action: "register", EntityType: "user", EntityID: u64str(user.ID),
		Detail: "用户注册", IP: ip, RequestID: requestID,
	})
	return user, nil
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest, ip, requestID string) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByAccount(req.Account)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInvalidCredential, constants.MsgInvalidCredential, err)
	}
	if !util.CheckPassword(user.PasswordHash, req.Password) {
		return nil, util.NewAppError(constants.CodeInvalidCredential, constants.MsgInvalidCredential, errors.New("password mismatch"))
	}
	if user.Status != "active" {
		return nil, util.NewAppError(constants.CodeUserBanned, constants.MsgUserBanned, errors.New("user banned"))
	}
	token, err := util.GenerateToken(s.cfg.JWTSecret, user.ID, user.Username, user.Role, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserLoggedIn, "user_id", user.ID, "username", user.Username)
	_ = s.audit.Record(&model.AuditLog{
		UserID: user.ID, Action: "login", EntityType: "user", EntityID: u64str(user.ID),
		Detail: "用户登录", IP: ip, RequestID: requestID,
	})
	return &dto.LoginResponse{Token: token, User: dto.ToUserResponse(user)}, nil
}

func (s *userService) GetByID(userID uint64) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, constants.MsgInvalidCredential+"：用户不存在", err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest, ip, requestID string) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeUserNotFound, constants.MsgInvalidCredential+"：用户不存在", err)
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if err := s.repo.Update(user); err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserUpdated, "user_id", user.ID, "nickname", user.Nickname)
	_ = s.audit.Record(&model.AuditLog{
		UserID: user.ID, Action: "update_profile", EntityType: "user", EntityID: u64str(user.ID),
		Detail: "更新个人资料", IP: ip, RequestID: requestID,
	})
	return user, nil
}

func (s *userService) SeedAdmin(username, email, password string) error {
	if _, err := s.repo.FindByUsername(username); err == nil {
		return nil
	}
	hash, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	admin := &model.User{
		Username: username, Email: email, PasswordHash: hash,
		Nickname: "管理员", Role: constants.RoleAdmin, Status: "active",
	}
	if err := s.repo.Create(admin); err != nil && !errors.Is(err, repository.ErrDuplicate) {
		return err
	}
	s.logger.Info(constants.LogSeedAdminCreated, "username", username)
	return nil
}
