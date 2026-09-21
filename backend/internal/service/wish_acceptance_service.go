package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// WishAcceptanceService 心愿验收服务：送交（事务+唯一索引）、发布者确认/驳回（事务+行锁）。
type WishAcceptanceService interface {
	Submit(ctx context.Context, userID, claimID uint64, req dto.SubmitAcceptanceRequest, ip, requestID string) (*model.WishAcceptance, error)
	Confirm(ctx context.Context, userID, acceptanceID uint64, ip, requestID string) (*model.WishAcceptance, error)
	Reject(ctx context.Context, userID, acceptanceID uint64, req dto.RejectAcceptanceRequest, ip, requestID string) (*model.WishAcceptance, error)
	GetByWishID(userID, wishID uint64) (*model.WishAcceptance, error)
}

type wishAcceptanceService struct {
	tx         repository.TxManager
	wish       repository.WishRepository
	claim      repository.WishClaimRepository
	acceptance repository.WishAcceptanceRepository
	badge      BadgeService
	audit      AuditService
	logger     *slog.Logger
}

// NewWishAcceptanceService 构造验收服务。
func NewWishAcceptanceService(
	tx repository.TxManager,
	wish repository.WishRepository,
	claim repository.WishClaimRepository,
	acceptance repository.WishAcceptanceRepository,
	badge BadgeService,
	audit AuditService,
	logger *slog.Logger,
) WishAcceptanceService {
	return &wishAcceptanceService{tx: tx, wish: wish, claim: claim, acceptance: acceptance, badge: badge, audit: audit, logger: logger}
}

// Submit 圆梦人送交验收：同一心愿仅一条记录；驳回后保留说明可重新送交；待验收期间禁止重复送交。
func (s *wishAcceptanceService) Submit(ctx context.Context, userID, claimID uint64, req dto.SubmitAcceptanceRequest, ip, requestID string) (*model.WishAcceptance, error) {
	if req.Note == "" {
		return nil, util.NewAppError(constants.CodeValidationFailed, "送交验收必须填写说明", errors.New("acceptance note empty"))
	}
	var saved *model.WishAcceptance
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		claim, err := s.claim.FindByID(claimID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
			}
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if claim.UserID != userID {
			return util.NewAppError(constants.CodeClaimNotOwner, "只有圆梦人才能送交验收", errors.New("claim owner mismatch"))
		}
		wish, err := s.wish.FindByID(claim.WishID)
		if err != nil {
			return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		if wish.Status == constants.WishStatusCompleted {
			return util.NewAppError(constants.CodeWishStatusInvalid, "心愿已完成，无需送交验收", errors.New("wish already completed"))
		}
		now := time.Now()
		existing, err := s.acceptance.FindByWishIDForUpdate(tx, claim.WishID)
		switch {
		case err == nil:
			switch existing.Status {
			case constants.AcceptanceStatusPending:
				return util.NewAppError(constants.CodeAcceptanceAlreadySubmitted, constants.MsgAcceptanceAlreadySubmitted, errors.New("acceptance already pending"))
			case constants.AcceptanceStatusConfirmed:
				return util.NewAppError(constants.CodeWishStatusInvalid, "心愿已完成验收", errors.New("acceptance already confirmed"))
			}
			// 已驳回记录重新送交：保留历史说明可覆盖，清空驳回原因，状态回到待验收。
			existing.Status = constants.AcceptanceStatusPending
			existing.Note = req.Note
			existing.RejectReason = ""
			existing.SubmittedAt = &now
			existing.ReviewedAt = nil
			existing.ReviewedBy = 0
			if err := s.acceptance.UpdateWithTx(tx, existing); err != nil {
				return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
			}
			saved = existing
		case errors.Is(err, repository.ErrNotFound):
			acc := &model.WishAcceptance{
				WishID: claim.WishID, ClaimID: claim.ID, UserID: claim.UserID,
				Note: req.Note, Status: constants.AcceptanceStatusPending, SubmittedAt: &now,
			}
			if err := s.acceptance.CreateWithTx(tx, acc); err != nil {
				if errors.Is(err, repository.ErrConflict) {
					return util.NewAppError(constants.CodeAcceptanceAlreadySubmitted, constants.MsgAcceptanceAlreadySubmitted, err)
				}
				return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
			}
			saved = acc
		default:
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		// 送交即视为圆梦进度 100%，心愿保持圆梦中等待发布者验收。
		claim.Progress = 100
		claim.LatestNote = req.Note
		if err := s.claim.UpdateWithTx(tx, claim); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogAcceptanceSubmitted, "acceptance_id", saved.ID, "wish_id", saved.WishID, "claim_id", saved.ClaimID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "submit_acceptance", EntityType: "wish_acceptance", EntityID: u64str(saved.ID),
		Detail: "送交验收：" + util.FormatAcceptanceStatus(saved.Status), IP: ip, RequestID: requestID,
	})
	return saved, nil
}

// Confirm 发布者确认验收：行锁保证并发验收只成功一次；确认后心愿完成并发放徽章。
func (s *wishAcceptanceService) Confirm(ctx context.Context, userID, acceptanceID uint64, ip, requestID string) (*model.WishAcceptance, error) {
	acc, err := s.review(ctx, userID, acceptanceID, constants.AcceptanceStatusConfirmed, "")
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogAcceptanceConfirmed, "acceptance_id", acc.ID, "wish_id", acc.WishID, "reviewer_id", userID)
	s.logger.Info(constants.LogClaimCompleted, "claim_id", acc.ClaimID, "user_id", acc.UserID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "confirm_acceptance", EntityType: "wish_acceptance", EntityID: u64str(acc.ID),
		Detail: "验收通过：" + util.FormatAcceptanceStatus(acc.Status), IP: ip, RequestID: requestID,
	})
	if err := s.badge.GrantCompletionBadges(acc.UserID); err != nil {
		s.logger.Warn("grant completion badge failed", "error", err)
	}
	return acc, nil
}

// Reject 发布者驳回验收：必须填写原因；记录保留说明回到圆梦中，可重新送交。
func (s *wishAcceptanceService) Reject(ctx context.Context, userID, acceptanceID uint64, req dto.RejectAcceptanceRequest, ip, requestID string) (*model.WishAcceptance, error) {
	if req.Reason == "" {
		return nil, util.NewAppError(constants.CodeValidationFailed, constants.MsgRejectReasonRequired, errors.New("reject reason empty"))
	}
	acc, err := s.review(ctx, userID, acceptanceID, constants.AcceptanceStatusRejected, req.Reason)
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogAcceptanceRejected, "acceptance_id", acc.ID, "wish_id", acc.WishID, "reviewer_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "reject_acceptance", EntityType: "wish_acceptance", EntityID: u64str(acc.ID),
		Detail: "驳回验收：" + util.TruncateString(req.Reason, 50), IP: ip, RequestID: requestID,
	})
	return acc, nil
}

// review 确认/驳回共用流程：FOR UPDATE 锁定验收记录，仅待验收状态可流转，失败不改动任何数据。
func (s *wishAcceptanceService) review(ctx context.Context, userID, acceptanceID uint64, target, reason string) (*model.WishAcceptance, error) {
	var reviewed *model.WishAcceptance
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		acc, err := s.acceptance.FindByIDForUpdate(tx, acceptanceID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeAcceptanceNotFound, constants.MsgAcceptanceNotFound, err)
			}
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		wish, err := s.wish.FindByID(acc.WishID)
		if err != nil {
			return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		if wish.UserID != userID {
			return util.NewAppError(constants.CodeAcceptanceNotPublisher, constants.MsgAcceptanceNotPublisher, errors.New("acceptance reviewer mismatch"))
		}
		if acc.Status != constants.AcceptanceStatusPending {
			return util.NewAppError(constants.CodeAcceptanceStatusInvalid, constants.MsgAcceptanceStatusInvalid, errors.New("acceptance status not pending"))
		}
		now := time.Now()
		acc.Status = target
		acc.ReviewedAt = &now
		acc.ReviewedBy = userID
		if target == constants.AcceptanceStatusRejected {
			// 驳回保留送交说明，记录回到圆梦中（心愿/认领保持 in_progress 不动）。
			acc.RejectReason = reason
		}
		if err := s.acceptance.UpdateWithTx(tx, acc); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if target == constants.AcceptanceStatusConfirmed {
			claim, err := s.claim.FindByID(acc.ClaimID)
			if err != nil {
				return util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
			}
			claim.Status = constants.WishStatusCompleted
			if err := s.claim.UpdateWithTx(tx, claim); err != nil {
				return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
			}
			wish.Status = constants.WishStatusCompleted
			wish.CompletionNote = acc.Note
			if err := s.wish.UpdateWithTx(tx, wish); err != nil {
				return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
			}
		}
		reviewed = acc
		return nil
	})
	if err != nil {
		return nil, err
	}
	return reviewed, nil
}

// GetByWishID 查询某心愿的验收记录（仅发布者与圆梦人可见）。
func (s *wishAcceptanceService) GetByWishID(userID, wishID uint64) (*model.WishAcceptance, error) {
	acc, err := s.acceptance.FindByWishID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeAcceptanceNotFound, constants.MsgAcceptanceNotFound, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
	}
	if wish.UserID != userID && acc.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, "只有心愿发布者或圆梦人才能查看验收记录", errors.New("acceptance visibility forbidden"))
	}
	return acc, nil
}
