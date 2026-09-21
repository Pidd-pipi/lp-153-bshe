package service

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// WishClaimService 心愿认领服务：认领（事务+行锁）、进度更新、完成、我的认领。
type WishClaimService interface {
	Claim(ctx context.Context, userID, wishID uint64, ip, requestID string) (*model.WishClaim, error)
	UpdateProgress(ctx context.Context, userID, claimID uint64, req dto.UpdateProgressRequest, ip, requestID string) (*model.WishClaim, error)
	Complete(ctx context.Context, userID, claimID uint64, req dto.CompleteClaimRequest, ip, requestID string) (*model.WishClaim, error)
	ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error)
	GetByWishID(userID, wishID uint64) (*model.WishClaim, error)
}

type wishClaimService struct {
	tx     repository.TxManager
	wish   repository.WishRepository
	claim  repository.WishClaimRepository
	user   repository.UserRepository
	badge  BadgeService
	audit  AuditService
	logger *slog.Logger
}

// NewWishClaimService 构造认领服务。
func NewWishClaimService(
	tx repository.TxManager,
	wish repository.WishRepository,
	claim repository.WishClaimRepository,
	user repository.UserRepository,
	badge BadgeService,
	audit AuditService,
	logger *slog.Logger,
) WishClaimService {
	return &wishClaimService{tx: tx, wish: wish, claim: claim, user: user, badge: badge, audit: audit, logger: logger}
}

// Claim 认领心愿：SELECT ... FOR UPDATE 锁定心愿行防止并发重复认领。
func (s *wishClaimService) Claim(ctx context.Context, userID, wishID uint64, ip, requestID string) (*model.WishClaim, error) {
	var created *model.WishClaim
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		wish, err := s.wish.FindByIDForUpdate(tx, wishID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
			}
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if wish.UserID == userID {
			return util.NewAppError(constants.CodeWishStatusInvalid, "不能认领自己发布的心愿", errors.New("self claim forbidden"))
		}
		if wish.Status != constants.WishStatusPending {
			return util.NewAppError(constants.CodeWishAlreadyClaimed, constants.MsgWishAlreadyClaimed, errors.New("wish status not pending"))
		}
		claim := &model.WishClaim{
			WishID: wishID, UserID: userID,
			Progress: 0, Status: constants.WishStatusClaimed,
		}
		if err := s.claim.CreateWithTx(tx, claim); err != nil {
			if errors.Is(err, repository.ErrConflict) {
				return util.NewAppError(constants.CodeWishAlreadyClaimed, constants.MsgWishAlreadyClaimed, err)
			}
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		wish.Status = constants.WishStatusClaimed
		if err := s.wish.UpdateWithTx(tx, wish); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		created = claim
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogWishClaimed, "wish_id", wishID, "claim_id", created.ID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "claim_wish", EntityType: "wish_claim", EntityID: u64str(created.ID),
		Detail: "认领心愿，成为圆梦人", IP: ip, RequestID: requestID,
	})
	if err := s.badge.GrantFirstClaim(userID); err != nil {
		s.logger.Warn("grant first claim badge failed", "error", err)
	}
	return created, nil
}

// UpdateProgress 更新进度：progress>=100 时状态机流转为 completed。
func (s *wishClaimService) UpdateProgress(ctx context.Context, userID, claimID uint64, req dto.UpdateProgressRequest, ip, requestID string) (*model.WishClaim, error) {
	var updated *model.WishClaim
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		claim, err := s.claim.FindByID(claimID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
			}
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if claim.UserID != userID {
			return util.NewAppError(constants.CodeClaimNotOwner, "只有圆梦人才能更新进度", errors.New("claim owner mismatch"))
		}
		wish, err := s.wish.FindByID(claim.WishID)
		if err != nil {
			return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		claim.Progress = req.Progress
		if req.Note != "" {
			claim.LatestNote = req.Note
		}
		if req.IsMilestone {
			claim.MilestoneCount++
		}
		if claim.Progress >= 100 {
			claim.Status = constants.WishStatusCompleted
			wish.Status = constants.WishStatusCompleted
			wish.CompletionNote = req.Note
		} else {
			if claim.Status == constants.WishStatusClaimed {
				claim.Status = constants.WishStatusInProgress
			}
			wish.Status = constants.WishStatusInProgress
		}
		if err := s.claim.UpdateWithTx(tx, claim); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if err := s.wish.UpdateWithTx(tx, wish); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		updated = claim
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogClaimProgress, "claim_id", claimID, "progress", updated.Progress, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "update_progress", EntityType: "wish_claim", EntityID: u64str(claimID),
		Detail: "更新圆梦进度至 " + itoa(updated.Progress) + "%", IP: ip, RequestID: requestID,
	})
	if updated.Status == constants.WishStatusCompleted {
		s.logger.Info(constants.LogClaimCompleted, "claim_id", claimID, "user_id", userID)
		_ = s.audit.Record(&model.AuditLog{
			UserID: userID, Action: "complete_wish", EntityType: "wish_claim", EntityID: u64str(claimID),
			Detail: "心愿达成，进入庆祝时刻", IP: ip, RequestID: requestID,
		})
		if err := s.badge.GrantCompletionBadges(userID); err != nil {
			s.logger.Warn("grant completion badge failed", "error", err)
		}
	}
	return updated, nil
}

// Complete 直接完成心愿（进度置 100）。
func (s *wishClaimService) Complete(ctx context.Context, userID, claimID uint64, req dto.CompleteClaimRequest, ip, requestID string) (*model.WishClaim, error) {
	return s.UpdateProgress(ctx, userID, claimID, dto.UpdateProgressRequest{Progress: 100, Note: req.Note, IsMilestone: true}, ip, requestID)
}

func (s *wishClaimService) ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	total, err := s.claim.CountByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	claims, err := s.claim.ListByUserID(userID, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items := make([]dto.WishClaimResponse, 0, len(claims))
	for i := range claims {
		wishTitle := ""
		if wish, werr := s.wish.FindByID(claims[i].WishID); werr == nil {
			wishTitle = wish.Title
		}
		name := ""
		if u, uerr := s.user.FindByID(userID); uerr == nil {
			name = u.Nickname
		}
		items = append(items, dto.ToWishClaimResponse(&claims[i], wishTitle, name))
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

// GetByWishID 查询某心愿的认领（心愿详情页圆梦人模块复用）。
func (s *wishClaimService) GetByWishID(userID, wishID uint64) (*model.WishClaim, error) {
	claim, err := s.claim.FindByWishID(wishID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
	}
	if claim.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, constants.MsgNeedLogin+"：无权限查看该认领", errors.New("claim visibility forbidden"))
	}
	return claim, nil
}
