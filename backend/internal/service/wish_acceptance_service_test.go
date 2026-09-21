package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

func acceptanceTestRepos() (*mockWishRepo, *mockClaimRepo) {
	wishRepo := &mockWishRepo{
		findByIDFn: func(id uint64) (*model.Wish, error) {
			return &model.Wish{ID: id, UserID: 1, Status: constants.WishStatusInProgress}, nil
		},
		updateWithTxFn: func(tx *gorm.DB, wish *model.Wish) error { return nil },
	}
	claimRepo := &mockClaimRepo{
		findByIDFn: func(id uint64) (*model.WishClaim, error) {
			return &model.WishClaim{ID: id, WishID: 5, UserID: 2, Progress: 90, Status: constants.WishStatusInProgress}, nil
		},
		updateWithTxFn: func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
	}
	return wishRepo, claimRepo
}

func TestWishAcceptanceService_Submit(t *testing.T) {
	t.Parallel()
	t.Run("first submit creates pending record", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		acceptanceRepo := &mockAcceptanceRepo{
			findByWishForUpdateFn: func(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error) {
				return nil, repository.ErrNotFound
			},
			createWithTxFn: func(tx *gorm.DB, acc *model.WishAcceptance) error {
				acc.ID = 11
				return nil
			},
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())
		acc, err := svc.Submit(context.Background(), 2, 9, dto.SubmitAcceptanceRequest{Note: "极光照片已送达"}, "127.0.0.1", "req-a1")
		if err != nil {
			t.Fatalf("submit: %v", err)
		}
		if acc.Status != constants.AcceptanceStatusPending {
			t.Fatalf("expected pending, got %s", acc.Status)
		}
		if acc.WishID != 5 || acc.ClaimID != 9 || acc.UserID != 2 {
			t.Fatalf("unexpected acceptance linkage: %+v", acc)
		}
	})

	t.Run("duplicate submit while pending forbidden", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		acceptanceRepo := &mockAcceptanceRepo{
			findByWishForUpdateFn: func(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error) {
				return &model.WishAcceptance{ID: 11, WishID: wishID, Status: constants.AcceptanceStatusPending}, nil
			},
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Submit(context.Background(), 2, 9, dto.SubmitAcceptanceRequest{Note: "重复送交"}, "127.0.0.1", "req-a2")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceAlreadySubmitted {
			t.Fatalf("expected CodeAcceptanceAlreadySubmitted, got %v", err)
		}
	})

	t.Run("unique index conflict maps to duplicate submit", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		acceptanceRepo := &mockAcceptanceRepo{
			findByWishForUpdateFn: func(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error) {
				return nil, repository.ErrNotFound
			},
			createWithTxFn: func(tx *gorm.DB, acc *model.WishAcceptance) error {
				return repository.ErrConflict
			},
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Submit(context.Background(), 2, 9, dto.SubmitAcceptanceRequest{Note: "并发送交"}, "127.0.0.1", "req-a3")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceAlreadySubmitted {
			t.Fatalf("expected CodeAcceptanceAlreadySubmitted, got %v", err)
		}
	})

	t.Run("resubmit after rejection keeps record and clears reason", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		existing := &model.WishAcceptance{
			ID: 11, WishID: 5, ClaimID: 9, UserID: 2,
			Note: "旧说明", Status: constants.AcceptanceStatusRejected, RejectReason: "照片不清晰",
		}
		acceptanceRepo := &mockAcceptanceRepo{
			findByWishForUpdateFn: func(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error) {
				return existing, nil
			},
			updateWithTxFn: func(tx *gorm.DB, acc *model.WishAcceptance) error { return nil },
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())
		acc, err := svc.Submit(context.Background(), 2, 9, dto.SubmitAcceptanceRequest{Note: "已补拍清晰照片"}, "127.0.0.1", "req-a4")
		if err != nil {
			t.Fatalf("resubmit: %v", err)
		}
		if acc.ID != 11 {
			t.Fatalf("resubmit must reuse the same record, got id %d", acc.ID)
		}
		if acc.Status != constants.AcceptanceStatusPending {
			t.Fatalf("expected pending after resubmit, got %s", acc.Status)
		}
		if acc.Note != "已补拍清晰照片" {
			t.Fatalf("expected note updated, got %s", acc.Note)
		}
		if acc.RejectReason != "" {
			t.Fatalf("expected reject reason cleared, got %s", acc.RejectReason)
		}
	})

	t.Run("submit requires note", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, &mockAcceptanceRepo{}, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Submit(context.Background(), 2, 9, dto.SubmitAcceptanceRequest{Note: ""}, "127.0.0.1", "req-a5")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeValidationFailed {
			t.Fatalf("expected CodeValidationFailed, got %v", err)
		}
	})

	t.Run("only fulfiller can submit", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, &mockAcceptanceRepo{}, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Submit(context.Background(), 99, 9, dto.SubmitAcceptanceRequest{Note: "越权送交"}, "127.0.0.1", "req-a6")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeClaimNotOwner {
			t.Fatalf("expected CodeClaimNotOwner, got %v", err)
		}
	})
}

func TestWishAcceptanceService_Confirm(t *testing.T) {
	t.Parallel()
	newPendingRepo := func(acc *model.WishAcceptance) *mockAcceptanceRepo {
		return &mockAcceptanceRepo{
			findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.WishAcceptance, error) {
				return acc, nil
			},
			updateWithTxFn: func(tx *gorm.DB, a *model.WishAcceptance) error { return nil },
		}
	}

	t.Run("confirm completes wish and grants badges", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		var savedWish *model.Wish
		var savedClaim *model.WishClaim
		wishRepo.updateWithTxFn = func(tx *gorm.DB, wish *model.Wish) error { savedWish = wish; return nil }
		claimRepo.updateWithTxFn = func(tx *gorm.DB, claim *model.WishClaim) error { savedClaim = claim; return nil }
		acc := &model.WishAcceptance{ID: 11, WishID: 5, ClaimID: 9, UserID: 2, Note: "完成了", Status: constants.AcceptanceStatusPending}
		badgeGranted := false
		badge := &mockBadge{grantCompletionFn: func(userID uint64) error { badgeGranted = true; return nil }}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, newPendingRepo(acc), badge, &mockAudit{}, testLogger())

		got, err := svc.Confirm(context.Background(), 1, 11, "127.0.0.1", "req-c1")
		if err != nil {
			t.Fatalf("confirm: %v", err)
		}
		if got.Status != constants.AcceptanceStatusConfirmed {
			t.Fatalf("expected confirmed, got %s", got.Status)
		}
		if got.ReviewedBy != 1 || got.ReviewedAt == nil {
			t.Fatalf("expected reviewer recorded, got %+v", got)
		}
		if savedWish == nil || savedWish.Status != constants.WishStatusCompleted {
			t.Fatalf("wish should be completed, got %+v", savedWish)
		}
		if savedWish.CompletionNote != "完成了" {
			t.Fatalf("completion note should come from acceptance note, got %s", savedWish.CompletionNote)
		}
		if savedClaim == nil || savedClaim.Status != constants.WishStatusCompleted {
			t.Fatalf("claim should be completed, got %+v", savedClaim)
		}
		if !badgeGranted {
			t.Fatal("completion badges should be granted after confirm")
		}
	})

	t.Run("only publisher can confirm", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		acc := &model.WishAcceptance{ID: 11, WishID: 5, ClaimID: 9, UserID: 2, Status: constants.AcceptanceStatusPending}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, newPendingRepo(acc), &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Confirm(context.Background(), 2, 11, "127.0.0.1", "req-c2")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceNotPublisher {
			t.Fatalf("expected CodeAcceptanceNotPublisher, got %v", err)
		}
	})

	t.Run("concurrent confirm only succeeds once without side effects", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		wishTouched := false
		claimTouched := false
		wishRepo.updateWithTxFn = func(tx *gorm.DB, wish *model.Wish) error { wishTouched = true; return nil }
		claimRepo.updateWithTxFn = func(tx *gorm.DB, claim *model.WishClaim) error { claimTouched = true; return nil }
		// 模拟并发下第二个请求：行锁后读到状态已非 pending。
		acc := &model.WishAcceptance{ID: 11, WishID: 5, ClaimID: 9, UserID: 2, Status: constants.AcceptanceStatusConfirmed}
		badgeGranted := false
		badge := &mockBadge{grantCompletionFn: func(userID uint64) error { badgeGranted = true; return nil }}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, newPendingRepo(acc), badge, &mockAudit{}, testLogger())

		_, err := svc.Confirm(context.Background(), 1, 11, "127.0.0.1", "req-c3")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceStatusInvalid {
			t.Fatalf("expected CodeAcceptanceStatusInvalid, got %v", err)
		}
		if wishTouched || claimTouched {
			t.Fatal("failed confirm must not modify wish or claim")
		}
		if badgeGranted {
			t.Fatal("failed confirm must not grant badges")
		}
	})
}

func TestWishAcceptanceService_Reject(t *testing.T) {
	t.Parallel()
	t.Run("reject keeps note and records reason", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		wishTouched := false
		wishRepo.updateWithTxFn = func(tx *gorm.DB, wish *model.Wish) error { wishTouched = true; return nil }
		acc := &model.WishAcceptance{ID: 11, WishID: 5, ClaimID: 9, UserID: 2, Note: "送交说明", Status: constants.AcceptanceStatusPending}
		acceptanceRepo := &mockAcceptanceRepo{
			findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.WishAcceptance, error) { return acc, nil },
			updateWithTxFn:      func(tx *gorm.DB, a *model.WishAcceptance) error { return nil },
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())

		got, err := svc.Reject(context.Background(), 1, 11, dto.RejectAcceptanceRequest{Reason: "照片不清晰"}, "127.0.0.1", "req-r1")
		if err != nil {
			t.Fatalf("reject: %v", err)
		}
		if got.Status != constants.AcceptanceStatusRejected {
			t.Fatalf("expected rejected, got %s", got.Status)
		}
		if got.RejectReason != "照片不清晰" {
			t.Fatalf("expected reject reason recorded, got %s", got.RejectReason)
		}
		if got.Note != "送交说明" {
			t.Fatalf("note must be preserved after reject, got %s", got.Note)
		}
		if wishTouched {
			t.Fatal("reject must not change wish status (stays 圆梦中)")
		}
	})

	t.Run("reject requires reason", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, &mockAcceptanceRepo{}, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Reject(context.Background(), 1, 11, dto.RejectAcceptanceRequest{Reason: ""}, "127.0.0.1", "req-r2")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeValidationFailed {
			t.Fatalf("expected CodeValidationFailed, got %v", err)
		}
	})

	t.Run("only publisher can reject", func(t *testing.T) {
		t.Parallel()
		wishRepo, claimRepo := acceptanceTestRepos()
		acc := &model.WishAcceptance{ID: 11, WishID: 5, ClaimID: 9, UserID: 2, Status: constants.AcceptanceStatusPending}
		acceptanceRepo := &mockAcceptanceRepo{
			findByIDForUpdateFn: func(tx *gorm.DB, id uint64) (*model.WishAcceptance, error) { return acc, nil },
		}
		svc := NewWishAcceptanceService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())
		_, err := svc.Reject(context.Background(), 2, 11, dto.RejectAcceptanceRequest{Reason: "越权驳回"}, "127.0.0.1", "req-r3")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceNotPublisher {
			t.Fatalf("expected CodeAcceptanceNotPublisher, got %v", err)
		}
	})
}

func TestWishAcceptanceService_GetByWishID(t *testing.T) {
	t.Parallel()
	wishRepo, _ := acceptanceTestRepos()
	now := time.Now()
	acceptanceRepo := &mockAcceptanceRepo{
		findByWishFn: func(wishID uint64) (*model.WishAcceptance, error) {
			return &model.WishAcceptance{ID: 11, WishID: wishID, UserID: 2, Status: constants.AcceptanceStatusPending, SubmittedAt: &now}, nil
		},
	}
	svc := NewWishAcceptanceService(&mockTx{}, wishRepo, &mockClaimRepo{}, acceptanceRepo, &mockBadge{}, &mockAudit{}, testLogger())

	if _, err := svc.GetByWishID(1, 5); err != nil {
		t.Fatalf("publisher should read acceptance: %v", err)
	}
	if _, err := svc.GetByWishID(2, 5); err != nil {
		t.Fatalf("fulfiller should read acceptance: %v", err)
	}
	if _, err := svc.GetByWishID(99, 5); err == nil {
		t.Fatal("stranger should not read acceptance")
	}
}
