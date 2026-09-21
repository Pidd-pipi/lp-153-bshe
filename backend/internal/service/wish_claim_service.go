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

// WishClaimService 心愿认领服务：认领（事务+行锁）、进度更新、送交验收、发布者验收、我的认领。
type WishClaimService interface {
	Claim(ctx context.Context, userID, wishID uint64, ip, requestID string) (*model.WishClaim, error)
	UpdateProgress(ctx context.Context, userID, claimID uint64, req dto.UpdateProgressRequest, ip, requestID string) (*model.WishClaim, error)
	SubmitForAcceptance(ctx context.Context, userID, claimID uint64, req dto.SubmitClaimRequest, ip, requestID string) (*model.WishClaim, error)
	Review(ctx context.Context, publisherID, claimID uint64, req dto.ReviewClaimRequest, ip, requestID string) (*model.WishClaim, error)
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

// UpdateProgress 更新圆梦进度。
// 待验收（pending_acceptance）期间禁止更新进度；进度到达 100% 也不直接完成，
// 必须由圆梦人送交验收、发布者确认后才完成。
func (s *wishClaimService) UpdateProgress(ctx context.Context, userID, claimID uint64, req dto.UpdateProgressRequest, ip, requestID string) (*model.WishClaim, error) {
	var updated *model.WishClaim
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		wish, claim, err := s.loadWishAndClaimForUpdate(tx, claimID)
		if err != nil {
			return err
		}
		if claim.UserID != userID {
			return util.NewAppError(constants.CodeClaimNotOwner, "圆梦人字段 user_id 不匹配：只有圆梦人才能更新进度", errors.New("claim owner mismatch"))
		}
		if claim.Status == constants.WishStatusPendingAcceptance {
			return util.NewAppError(constants.CodeClaimUnderReview, constants.MsgClaimUnderReview, errors.New("claim pending acceptance"))
		}
		if claim.Status == constants.WishStatusCompleted {
			return util.NewAppError(constants.CodeClaimStatusInvalid, "认领记录已完成，不能再更新进度", errors.New("claim already completed"))
		}
		claim.Progress = req.Progress
		if req.Note != "" {
			claim.LatestNote = req.Note
		}
		if req.IsMilestone {
			claim.MilestoneCount++
		}
		if claim.Status == constants.WishStatusClaimed {
			claim.Status = constants.WishStatusInProgress
		}
		// 进度即使达到 100% 也保持圆梦中，等待圆梦人送交验收。
		wish.Status = constants.WishStatusInProgress
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
	return updated, nil
}

// SubmitForAcceptance 圆梦人送交验收：填写送交说明，认领记录与心愿一并进入待验收。
// 待验收期间重复送交被拒绝；驳回后（记录回到圆梦中、说明保留）可重新送交。
func (s *wishClaimService) SubmitForAcceptance(ctx context.Context, userID, claimID uint64, req dto.SubmitClaimRequest, ip, requestID string) (*model.WishClaim, error) {
	note := req.Note
	if note == "" {
		// binding 标签之外的二次兜底，错误信息携带实体名与字段名。
		return nil, util.NewAppError(constants.CodeSubmissionNoteRequired, constants.MsgClaimSubmissionNote, errors.New("wish_claim.note is required"))
	}
	var submitted *model.WishClaim
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		wish, claim, err := s.loadWishAndClaimForUpdate(tx, claimID)
		if err != nil {
			return err
		}
		if claim.UserID != userID {
			return util.NewAppError(constants.CodeClaimNotOwner, "圆梦人字段 user_id 不匹配：只有圆梦人才能送交验收", errors.New("claim owner mismatch"))
		}
		if claim.Status == constants.WishStatusPendingAcceptance {
			return util.NewAppError(constants.CodeClaimUnderReview, constants.MsgClaimUnderReview, errors.New("claim already pending acceptance"))
		}
		if claim.Status == constants.WishStatusCompleted {
			return util.NewAppError(constants.CodeClaimStatusInvalid, "认领记录已完成，不能重复送交", errors.New("claim already completed"))
		}
		now := time.Now()
		claim.Status = constants.WishStatusPendingAcceptance
		claim.SubmissionNote = note
		claim.SubmittedAt = &now
		claim.RejectReason = ""
		claim.ReviewedAt = nil
		if claim.Progress < 100 {
			claim.Progress = 100
		}
		wish.Status = constants.WishStatusPendingAcceptance
		if err := s.claim.UpdateWithTx(tx, claim); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if err := s.wish.UpdateWithTx(tx, wish); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		submitted = claim
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogClaimSubmitted, "claim_id", claimID, "wish_id", submitted.WishID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "submit_for_acceptance", EntityType: "wish_claim", EntityID: u64str(claimID),
		Detail: "送交验收：" + util.TruncateString(note, 30), IP: ip, RequestID: requestID,
	})
	return submitted, nil
}

// Review 发布者验收：只有心愿发布者能确认（approve）或驳回（reject）。
// 事务内对认领行与心愿行加 FOR UPDATE 锁并二次校验状态，并发验收只能成功一次：
// 第二个事务阻塞后看到的状态已不是 pending_acceptance，直接失败且不改动任何记录、状态或徽章。
func (s *wishClaimService) Review(ctx context.Context, publisherID, claimID uint64, req dto.ReviewClaimRequest, ip, requestID string) (*model.WishClaim, error) {
	isReject := req.Action == "reject"
	if isReject && req.Reason == "" {
		return nil, util.NewAppError(constants.CodeRejectReasonRequired, constants.MsgRejectReason, errors.New("wish_claim.reject_reason is required"))
	}
	var reviewed *model.WishClaim
	var accepted bool
	err := s.tx.Transaction(func(tx *gorm.DB) error {
		wish, claim, err := s.loadWishAndClaimForUpdate(tx, claimID)
		if err != nil {
			return err
		}
		if wish.UserID != publisherID {
			return util.NewAppError(constants.CodeClaimNotPublisher, constants.MsgClaimNotPublisher+"（角色：发布者）", errors.New("wish publisher mismatch"))
		}
		if claim.Status != constants.WishStatusPendingAcceptance || wish.Status != constants.WishStatusPendingAcceptance {
			return util.NewAppError(constants.CodeClaimNotUnderReview, constants.MsgClaimNotUnderReview, errors.New("claim not pending acceptance"))
		}
		now := time.Now()
		claim.ReviewedAt = &now
		if isReject {
			// 驳回：记录回到圆梦中，保留送交说明，记录驳回原因，之后可重新送交。
			claim.Status = constants.WishStatusInProgress
			claim.RejectReason = req.Reason
			wish.Status = constants.WishStatusInProgress
		} else {
			// 确认：心愿完成，完成说明取送交说明，徽章在事务提交后发放。
			claim.Status = constants.WishStatusCompleted
			wish.Status = constants.WishStatusCompleted
			wish.CompletionNote = claim.SubmissionNote
			accepted = true
		}
		if err := s.claim.UpdateWithTx(tx, claim); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		if err := s.wish.UpdateWithTx(tx, wish); err != nil {
			return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		reviewed = claim
		return nil
	})
	if err != nil {
		return nil, err
	}
	if isReject {
		s.logger.Info(constants.LogClaimRejected, "claim_id", claimID, "wish_id", reviewed.WishID, "publisher_id", publisherID)
		_ = s.audit.Record(&model.AuditLog{
			UserID: publisherID, Action: "reject_claim", EntityType: "wish_claim", EntityID: u64str(claimID),
			Detail: "驳回验收：" + util.TruncateString(req.Reason, 30), IP: ip, RequestID: requestID,
		})
		return reviewed, nil
	}
	s.logger.Info(constants.LogClaimAccepted, "claim_id", claimID, "wish_id", reviewed.WishID, "publisher_id", publisherID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: publisherID, Action: "accept_claim", EntityType: "wish_claim", EntityID: u64str(claimID),
		Detail: "验收通过，心愿完成", IP: ip, RequestID: requestID,
	})
	if accepted {
		// 徽章仅在事务提交成功后发放，验收失败的请求不会触发任何徽章变更。
		if err := s.badge.GrantCompletionBadges(reviewed.UserID); err != nil {
			s.logger.Warn("grant completion badge failed", "error", err)
		}
	}
	return reviewed, nil
}

// loadWishAndClaimForUpdate 事务内按统一顺序（先心愿后认领，FOR UPDATE）加锁，
// 避免并发验收/送交/进度更新之间出现死锁与脏写。
func (s *wishClaimService) loadWishAndClaimForUpdate(tx *gorm.DB, claimID uint64) (*model.Wish, *model.WishClaim, error) {
	claim, err := s.claim.FindByIDForUpdate(tx, claimID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
		}
		return nil, nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	wish, err := s.wish.FindByIDForUpdate(tx, claim.WishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return nil, nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return wish, claim, nil
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

// GetByWishID 查询某心愿的认领（心愿详情页圆梦人/发布者模块复用）。
// 圆梦人本人与心愿发布者均可查看（详情页展示验收状态与处理入口）。
func (s *wishClaimService) GetByWishID(userID, wishID uint64) (*model.WishClaim, error) {
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	claim, err := s.claim.FindByWishID(wishID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeClaimNotFound, constants.MsgClaimNotFound, err)
	}
	if claim.UserID != userID && wish.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, constants.MsgNeedLogin+"：无权限查看该认领", errors.New("claim visibility forbidden"))
	}
	return claim, nil
}
