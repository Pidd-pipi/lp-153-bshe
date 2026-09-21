package service

import (
	"log/slog"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// AuditService 审计日志服务：service 埋点与中间件共用。
type AuditService interface {
	Record(log *model.AuditLog) error
	List(q dto.AuditLogQuery) (*dto.PageResult, error)
}

type auditService struct {
	repo   repository.AuditLogRepository
	user   repository.UserRepository
	logger *slog.Logger
}

// NewAuditService 构造审计服务。
func NewAuditService(repo repository.AuditLogRepository, user repository.UserRepository, logger *slog.Logger) AuditService {
	return &auditService{repo: repo, user: user, logger: logger}
}

func (s *auditService) Record(log *model.AuditLog) error {
	if err := s.repo.Create(log); err != nil {
		s.logger.Error(constants.LogAuditRecorded, "user_id", log.UserID, "action", log.Action, "error", err)
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogAuditRecorded, "user_id", log.UserID, "action", log.Action, "entity_type", log.EntityType, "entity_id", log.EntityID)
	return nil
}

func (s *auditService) List(q dto.AuditLogQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	filters := map[string]any{
		"user_id":     q.UserID,
		"action":      q.Action,
		"entity_type": q.EntityType,
	}
	total, err := s.repo.Count(filters)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	logs, err := s.repo.List(filters, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items := make([]dto.AuditLogResponse, 0, len(logs))
	for i := range logs {
		username := ""
		if logs[i].UserID > 0 {
			if u, uerr := s.user.FindByID(logs[i].UserID); uerr == nil {
				username = u.Username
			}
		}
		items = append(items, dto.ToAuditLogResponse(&logs[i], username))
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}
