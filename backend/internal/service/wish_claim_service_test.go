package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/util"
)

func TestWishClaimService_Claim(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		wishFn      func(tx *gorm.DB, id uint64) (*model.Wish, error)
		createFn    func(tx *gorm.DB, claim *model.WishClaim) error
		wantErrCode int
	}{
		{
			name: "claim success",
			wishFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
				return &model.Wish{ID: id, UserID: 1, Status: constants.WishStatusPending}, nil
			},
			createFn: func(tx *gorm.DB, claim *model.WishClaim) error {
				claim.ID = 7
				return nil
			},
			wantErrCode: 0,
		},
		{
			name: "wish already claimed",
			wishFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
				return &model.Wish{ID: id, UserID: 1, Status: constants.WishStatusClaimed}, nil
			},
			createFn:    func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
			wantErrCode: constants.CodeWishAlreadyClaimed,
		},
		{
			name: "self claim forbidden",
			wishFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
				return &model.Wish{ID: id, UserID: 2, Status: constants.WishStatusPending}, nil
			},
			createFn:    func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
			wantErrCode: constants.CodeWishStatusInvalid,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			wishRepo := &mockWishRepo{
				findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
					// 转换为真实签名：mock 忽略 tx
					return tt.wishFn(nil, id)
				},
				updateWithTxFn: func(tx *gorm.DB, wish *model.Wish) error { return nil },
			}
			claimRepo := &mockClaimRepo{
				createWithTxFn: func(tx *gorm.DB, claim *model.WishClaim) error {
					return tt.createFn(nil, claim)
				},
			}
			badge := &mockBadge{grantFirstClaimFn: func(userID uint64) error { return nil }}
			svc := NewWishClaimService(&mockTx{}, wishRepo, claimRepo, &mockUserRepo{}, badge, &mockAudit{}, testLogger())
			claim, err := svc.Claim(context.Background(), 2, 5, "127.0.0.1", "req-5")
			if tt.wantErrCode == 0 {
				if err != nil {
					t.Fatalf("claim should succeed, got %v", err)
				}
				if claim.WishID != 5 {
					t.Fatalf("unexpected wish id %d", claim.WishID)
				}
				return
			}
			var appErr *util.AppError
			if !errors.As(err, &appErr) || appErr.Code != tt.wantErrCode {
				t.Fatalf("expected code %d, got %v", tt.wantErrCode, err)
			}
		})
	}
}

// claimPair 构造共享同一条认领记录/心愿记录的 mock（模拟行锁后的串行效果）。
func claimPair(wish *model.Wish, claim *model.WishClaim) (*mockWishRepo, *mockClaimRepo) {
	var mu sync.Mutex
	wishRepo := &mockWishRepo{
		findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
			mu.Lock()
			defer mu.Unlock()
			cp := *wish
			return &cp, nil
		},
		updateWithTxFn: func(tx *gorm.DB, w *model.Wish) error {
			mu.Lock()
			defer mu.Unlock()
			*wish = *w
			return nil
		},
	}
	claimRepo := &mockClaimRepo{
		findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.WishClaim, error) {
			mu.Lock()
			defer mu.Unlock()
			cp := *claim
			return &cp, nil
		},
		updateWithTxFn: func(tx *gorm.DB, c *model.WishClaim) error {
			mu.Lock()
			defer mu.Unlock()
			*claim = *c
			return nil
		},
	}
	return wishRepo, claimRepo
}

func TestWishClaimService_SubmitForAcceptance(t *testing.T) {
	t.Parallel()
	wish := &model.Wish{ID: 5, UserID: 1, Status: constants.WishStatusInProgress}
	claim := &model.WishClaim{ID: 9, WishID: 5, UserID: 2, Progress: 80, Status: constants.WishStatusInProgress}
	wishRepo, claimRepo := claimPair(wish, claim)
	badge := &mockBadge{}
	svc := NewWishClaimService(&serialTx{}, wishRepo, claimRepo, &mockUserRepo{}, badge, &mockAudit{}, testLogger())

	got, err := svc.SubmitForAcceptance(context.Background(), 2, 9, dto.SubmitClaimRequest{Note: "全部完成，请验收"}, "127.0.0.1", "req-submit")
	if err != nil {
		t.Fatalf("submit should succeed, got %v", err)
	}
	if got.Status != constants.WishStatusPendingAcceptance {
		t.Fatalf("expected claim pending_acceptance, got %s", got.Status)
	}
	if got.Progress != 100 || got.SubmissionNote != "全部完成，请验收" || got.SubmittedAt == nil {
		t.Fatalf("unexpected claim after submit: %+v", got)
	}
	if wish.Status != constants.WishStatusPendingAcceptance {
		t.Fatalf("expected wish pending_acceptance, got %s", wish.Status)
	}

	// 待验收期间重复送交必须失败，且不得改动记录。
	_, err = svc.SubmitForAcceptance(context.Background(), 2, 9, dto.SubmitClaimRequest{Note: "再次送交"}, "127.0.0.1", "req-resubmit")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeClaimUnderReview {
		t.Fatalf("expected CodeClaimUnderReview, got %v", err)
	}
	if claim.SubmissionNote != "全部完成，请验收" {
		t.Fatalf("record must not change on rejected resubmit, got note %q", claim.SubmissionNote)
	}
}

func TestWishClaimService_UpdateProgress_BlockedUnderReview(t *testing.T) {
	t.Parallel()
	wish := &model.Wish{ID: 5, UserID: 1, Status: constants.WishStatusPendingAcceptance}
	claim := &model.WishClaim{ID: 9, WishID: 5, UserID: 2, Progress: 100, Status: constants.WishStatusPendingAcceptance}
	wishRepo, claimRepo := claimPair(wish, claim)
	svc := NewWishClaimService(&serialTx{}, wishRepo, claimRepo, &mockUserRepo{}, &mockBadge{}, &mockAudit{}, testLogger())

	_, err := svc.UpdateProgress(context.Background(), 2, 9, dto.UpdateProgressRequest{Progress: 90}, "127.0.0.1", "req-blocked")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeClaimUnderReview {
		t.Fatalf("expected CodeClaimUnderReview, got %v", err)
	}
	if claim.Progress != 100 || claim.Status != constants.WishStatusPendingAcceptance {
		t.Fatalf("record must not change, got %+v", claim)
	}
}

func TestWishClaimService_Review_Approve(t *testing.T) {
	t.Parallel()
	wish := &model.Wish{ID: 5, UserID: 1, Status: constants.WishStatusPendingAcceptance}
	claim := &model.WishClaim{ID: 9, WishID: 5, UserID: 2, Progress: 100,
		Status: constants.WishStatusPendingAcceptance, SubmissionNote: "完成说明"}
	wishRepo, claimRepo := claimPair(wish, claim)
	granted := false
	badge := &mockBadge{grantCompletionFn: func(userID uint64) error {
		if userID != 2 {
			t.Fatalf("badge should go to fulfiller 2, got %d", userID)
		}
		granted = true
		return nil
	}}
	svc := NewWishClaimService(&serialTx{}, wishRepo, claimRepo, &mockUserRepo{}, badge, &mockAudit{}, testLogger())

	// 非发布者不能验收。
	if _, err := svc.Review(context.Background(), 3, 9, dto.ReviewClaimRequest{Action: "approve"}, "ip", "r1"); err == nil {
		t.Fatal("non-publisher approve must fail")
	}
	// 圆梦人自己也不能验收。
	if _, err := svc.Review(context.Background(), 2, 9, dto.ReviewClaimRequest{Action: "approve"}, "ip", "r2"); err == nil {
		t.Fatal("fulfiller approve must fail")
	}

	got, err := svc.Review(context.Background(), 1, 9, dto.ReviewClaimRequest{Action: "approve"}, "127.0.0.1", "req-approve")
	if err != nil {
		t.Fatalf("approve should succeed, got %v", err)
	}
	if got.Status != constants.WishStatusCompleted || wish.Status != constants.WishStatusCompleted {
		t.Fatalf("expected both completed, claim=%s wish=%s", got.Status, wish.Status)
	}
	if wish.CompletionNote != "完成说明" {
		t.Fatalf("completion note should carry submission note, got %q", wish.CompletionNote)
	}
	if got.ReviewedAt == nil || !granted {
		t.Fatal("reviewed_at should be set and completion badge granted")
	}

	// 并发验收：第二次必须失败，且不改动记录/状态/徽章。
	granted = false
	_, err = svc.Review(context.Background(), 1, 9, dto.ReviewClaimRequest{Action: "approve"}, "127.0.0.1", "req-approve-2")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeClaimNotUnderReview {
		t.Fatalf("expected CodeClaimNotUnderReview on second approve, got %v", err)
	}
	if granted {
		t.Fatal("badge must not be granted on failed acceptance")
	}
	if claim.Status != constants.WishStatusCompleted || wish.Status != constants.WishStatusCompleted {
		t.Fatalf("failed acceptance must not mutate status, claim=%s wish=%s", claim.Status, wish.Status)
	}
}

func TestWishClaimService_Review_Reject(t *testing.T) {
	t.Parallel()
	wish := &model.Wish{ID: 5, UserID: 1, Status: constants.WishStatusPendingAcceptance}
	claim := &model.WishClaim{ID: 9, WishID: 5, UserID: 2, Progress: 100,
		Status: constants.WishStatusPendingAcceptance, SubmissionNote: "完成说明"}
	wishRepo, claimRepo := claimPair(wish, claim)
	granted := false
	badge := &mockBadge{grantCompletionFn: func(userID uint64) error { granted = true; return nil }}
	svc := NewWishClaimService(&serialTx{}, wishRepo, claimRepo, &mockUserRepo{}, badge, &mockAudit{}, testLogger())

	// 驳回不写原因必须失败，且不得改动记录。
	_, err := svc.Review(context.Background(), 1, 9, dto.ReviewClaimRequest{Action: "reject"}, "127.0.0.1", "r-no-reason")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeRejectReasonRequired {
		t.Fatalf("expected CodeRejectReasonRequired, got %v", err)
	}
	if claim.Status != constants.WishStatusPendingAcceptance {
		t.Fatalf("reject without reason must not mutate record, got %s", claim.Status)
	}

	got, err := svc.Review(context.Background(), 1, 9, dto.ReviewClaimRequest{Action: "reject", Reason: "缺少凭证"}, "127.0.0.1", "r-reject")
	if err != nil {
		t.Fatalf("reject should succeed, got %v", err)
	}
	if got.Status != constants.WishStatusInProgress || wish.Status != constants.WishStatusInProgress {
		t.Fatalf("expected back in_progress, claim=%s wish=%s", got.Status, wish.Status)
	}
	if got.RejectReason != "缺少凭证" || got.SubmissionNote != "完成说明" {
		t.Fatalf("reject reason should be saved and submission note retained, got %+v", got)
	}
	if granted {
		t.Fatal("badge must not be granted on rejection")
	}

	// 驳回后可重新送交。
	resubmitted, err := svc.SubmitForAcceptance(context.Background(), 2, 9, dto.SubmitClaimRequest{Note: "补充凭证，再次送交"}, "127.0.0.1", "r-resubmit")
	if err != nil {
		t.Fatalf("resubmit after rejection should succeed, got %v", err)
	}
	if resubmitted.Status != constants.WishStatusPendingAcceptance || resubmitted.RejectReason != "" {
		t.Fatalf("resubmit should clear reject reason, got %+v", resubmitted)
	}
}

// serialTx 串行执行事务函数，模拟数据库事务边界。
type serialTx struct{}

func (m *serialTx) Transaction(fn func(tx *gorm.DB) error) error { return fn(nil) }
