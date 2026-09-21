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
			createFn: func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
			wantErrCode: constants.CodeWishAlreadyClaimed,
		},
		{
			name: "self claim forbidden",
			wishFn: func(tx *gorm.DB, id uint64) (*model.Wish, error) {
				return &model.Wish{ID: id, UserID: 2, Status: constants.WishStatusPending}, nil
			},
			createFn: func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
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

func TestWishClaimService_UpdateProgress_Complete(t *testing.T) {
	t.Parallel()
	wishRepo := &mockWishRepo{
		findByIDFn: func(id uint64) (*model.Wish, error) {
			return &model.Wish{ID: id, UserID: 1, Status: constants.WishStatusClaimed}, nil
		},
		updateWithTxFn: func(tx *gorm.DB, wish *model.Wish) error { return nil },
		countCompletedFn: func(userID uint64) (int64, error) { return 11, nil },
	}
	claimRepo := &mockClaimRepo{
		findByIDFn: func(id uint64) (*model.WishClaim, error) {
			return &model.WishClaim{ID: id, WishID: 5, UserID: 2, Progress: 40, Status: constants.WishStatusInProgress}, nil
		},
		updateWithTxFn: func(tx *gorm.DB, claim *model.WishClaim) error { return nil },
	}
	badge := &mockBadge{grantCompletionFn: func(userID uint64) error { return nil }}
	svc := NewWishClaimService(&mockTx{}, wishRepo, claimRepo, &mockUserRepo{}, badge, &mockAudit{}, testLogger())

	claim, err := svc.UpdateProgress(context.Background(), 2, 9, dto.UpdateProgressRequest{Progress: 100, Note: "做到了！", IsMilestone: true}, "127.0.0.1", "req-6")
	if err != nil {
		t.Fatalf("update progress: %v", err)
	}
	if claim.Status != constants.WishStatusCompleted {
		t.Fatalf("expected completed, got %s", claim.Status)
	}
	if claim.MilestoneCount != 1 {
		t.Fatalf("expected milestone count 1, got %d", claim.MilestoneCount)
	}
}
