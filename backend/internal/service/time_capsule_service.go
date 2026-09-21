package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// TimeCapsuleService 时光胶囊服务：封存/查看/解锁/删除。
type TimeCapsuleService interface {
	Create(userID uint64, req dto.CreateCapsuleRequest, ip, requestID string) (*model.TimeCapsule, error)
	ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error)
	GetByID(userID, capsuleID uint64) (*model.TimeCapsule, error)
	Delete(userID, capsuleID uint64, ip, requestID string) error
	UnlockDue() ([]model.TimeCapsule, error)
}

type timeCapsuleService struct {
	capsule repository.TimeCapsuleRepository
	audit   AuditService
	logger  *slog.Logger
}

// NewTimeCapsuleService 构造胶囊服务。
func NewTimeCapsuleService(capsule repository.TimeCapsuleRepository, audit AuditService, logger *slog.Logger) TimeCapsuleService {
	return &timeCapsuleService{capsule: capsule, audit: audit, logger: logger}
}

func (s *timeCapsuleService) Create(userID uint64, req dto.CreateCapsuleRequest, ip, requestID string) (*model.TimeCapsule, error) {
	if req.UnlockAt.Before(time.Now()) {
		return nil, util.NewAppError(constants.CodeBadRequest, "解锁时间必须晚于当前时间", errors.New("unlock_at in past"))
	}
	status := constants.CapsuleStatusLocked
	capsule := &model.TimeCapsule{
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		ImageURLs: util.EncodeStringArray(req.ImageURLs),
		AudioURL:  req.AudioURL,
		UnlockAt:  req.UnlockAt,
		Status:    status,
	}
	if err := s.capsule.Create(capsule); err != nil {
		s.logger.Error(constants.LogCapsuleCreateFail, "user_id", userID, "error", err)
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCapsuleCreated, "capsule_id", capsule.ID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "create_capsule", EntityType: "time_capsule", EntityID: u64str(capsule.ID),
		Detail: "封存时光胶囊", IP: ip, RequestID: requestID,
	})
	return capsule, nil
}

func (s *timeCapsuleService) ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	total, err := s.capsule.CountByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items, err := s.capsule.ListByUserID(userID, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	resps := make([]dto.CapsuleResponse, 0, len(items))
	for i := range items {
		resps = append(resps, dto.ToCapsuleResponse(&items[i], items[i].Status == constants.CapsuleStatusLocked))
	}
	return &dto.PageResult{Items: resps, Total: total, Page: page, PageSize: size}, nil
}

// GetByID 查看胶囊：未解锁时仅本人可见标题，内容打码。
func (s *timeCapsuleService) GetByID(userID, capsuleID uint64) (*model.TimeCapsule, error) {
	capsule, err := s.capsule.FindByID(capsuleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeCapsuleNotFound, "时光胶囊不存在", err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if capsule.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, constants.MsgNeedLogin+"：无权查看该胶囊", errors.New("capsule visibility forbidden"))
	}
	if capsule.Status == constants.CapsuleStatusLocked {
		now := time.Now()
		if !capsule.UnlockAt.After(now) {
			capsule.Status = constants.CapsuleStatusUnlocked
			ts := now
			capsule.UnlockedAt = &ts
			if err := s.capsule.Update(capsule); err != nil {
				return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
			}
			s.logger.Info(constants.LogCapsuleUnlocked, "capsule_id", capsule.ID, "user_id", userID)
		}
	}
	return capsule, nil
}

func (s *timeCapsuleService) Delete(userID, capsuleID uint64, ip, requestID string) error {
	capsule, err := s.capsule.FindByID(capsuleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeCapsuleNotFound, "时光胶囊不存在", err)
		}
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if capsule.UserID != userID {
		return util.NewAppError(constants.CodeForbidden, constants.MsgNeedLogin+"：无权删除该胶囊", errors.New("capsule owner mismatch"))
	}
	if err := s.capsule.Delete(capsuleID); err != nil {
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCapsuleDeleted, "capsule_id", capsuleID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "delete_capsule", EntityType: "time_capsule", EntityID: u64str(capsuleID),
		Detail: "删除时光胶囊", IP: ip, RequestID: requestID,
	})
	return nil
}

// UnlockDue 解锁所有到期胶囊（启动时与定时任务调用）。
func (s *timeCapsuleService) UnlockDue() ([]model.TimeCapsule, error) {
	items, err := s.capsule.UnlockDue(time.Now())
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	now := time.Now()
	for i := range items {
		items[i].Status = constants.CapsuleStatusUnlocked
		items[i].UnlockedAt = &now
		if err := s.capsule.Update(&items[i]); err != nil {
			return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		s.logger.Info(constants.LogCapsuleUnlocked, "capsule_id", items[i].ID)
	}
	return items, nil
}
