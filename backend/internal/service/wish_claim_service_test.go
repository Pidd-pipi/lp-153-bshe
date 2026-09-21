package service

import (
	"context"
	"errors"
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
			svc := NewWishClaimService(&mockTx{}, wishRepo, claimRepo, &mockAcceptanceRepo{}, &mockUserRepo{}, badge, &mockAudit{}, testLogger())
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

func TestWishClaimService_UpdateProgress(t *testing.T) {
	t.Parallel()
	newSvc := func(acceptanceRepo *mockAcceptanceRepo) (WishClaimService, *mockWishRepo, *mockClaimRepo) {
		wishRepo := &mockWishRepo{
			findByIDFn: func(id uint64) (*model.Wish, error) {
				return &model.Wish{ID: id, UserID: 1, Status: constants.WishStatusInProgress}, nil
			},
			updateWithTxFn: func(tx *gorm.DB, wish *model.Wish) error { return nil },
		}
		claimRepo := &mockClaimRepo{
			findByIDFn: func(id uint64) (*model.WishClaim, error) {
				return &model.WishClaim{ID: id, WishID: 5, UserID: 2, Progress: 40, Status: constants.WishStatusInProgress}, nil
			},
			updateWithTxFn: func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
		}
		svc := NewWishClaimService(&mockTx{}, wishRepo, claimRepo, acceptanceRepo, &mockUserRepo{}, &mockBadge{}, &mockAudit{}, testLogger())
		return svc, wishRepo, claimRepo
	}

	t.Run("progress update below 100 succeeds", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := newSvc(&mockAcceptanceRepo{})
		claim, err := svc.UpdateProgress(context.Background(), 2, 9, dto.UpdateProgressRequest{Progress: 80, Note: "快到终点了", IsMilestone: true}, "127.0.0.1", "req-6")
		if err != nil {
			t.Fatalf("update progress: %v", err)
		}
		if claim.Progress != 80 {
			t.Fatalf("expected progress 80, got %d", claim.Progress)
		}
		if claim.MilestoneCount != 1 {
			t.Fatalf("expected milestone count 1, got %d", claim.MilestoneCount)
		}
		if claim.Status != constants.WishStatusInProgress {
			t.Fatalf("expected in_progress, got %s", claim.Status)
		}
	})

	t.Run("progress 100 requires acceptance instead of completing", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := newSvc(&mockAcceptanceRepo{})
		_, err := svc.UpdateProgress(context.Background(), 2, 9, dto.UpdateProgressRequest{Progress: 100, Note: "做到了！"}, "127.0.0.1", "req-7")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptanceRequired {
			t.Fatalf("expected CodeAcceptanceRequired, got %v", err)
		}
	})

	t.Run("progress update forbidden while acceptance pending", func(t *testing.T) {
		t.Parallel()
		acceptanceRepo := &mockAcceptanceRepo{
			findByWishFn: func(wishID uint64) (*model.WishAcceptance, error) {
				return &model.WishAcceptance{ID: 3, WishID: wishID, Status: constants.AcceptanceStatusPending}, nil
			},
		}
		svc, _, _ := newSvc(acceptanceRepo)
		_, err := svc.UpdateProgress(context.Background(), 2, 9, dto.UpdateProgressRequest{Progress: 60, Note: "想改进度"}, "127.0.0.1", "req-8")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.Code != constants.CodeAcceptancePending {
			t.Fatalf("expected CodeAcceptancePending, got %v", err)
		}
	})
}
