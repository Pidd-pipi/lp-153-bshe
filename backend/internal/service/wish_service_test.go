package service

import (
	"errors"
	"testing"
	"time"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/util"
)

func TestWishService_Create(t *testing.T) {
	t.Parallel()
	granted := false
	wishRepo := &mockWishRepo{
		createFn:      func(w *model.Wish) error { w.ID = 10; return nil },
		countByUserFn: func(userID uint64) (int64, error) { return 1, nil },
	}
	badge := &mockBadge{grantFirstWishFn: func(userID uint64) error { granted = true; return nil }}
	svc := NewWishService(wishRepo, &mockClaimRepo{}, &mockBlessRepo{}, &mockUserRepo{}, badge, &mockAudit{}, testLogger())

	deadline := time.Now().Add(30 * 24 * time.Hour)
	wish, err := svc.Create(1, dto.CreateWishRequest{
		Title: "去冰岛看极光", Content: "想在极夜中看到绿色的极光",
		Category: constants.CategoryTravel, Visibility: constants.VisibilityPublic,
		Difficulty: constants.DifficultyHard, ExpectedDeadline: &deadline,
	}, "127.0.0.1", "req-3")
	if err != nil {
		t.Fatalf("create wish: %v", err)
	}
	if wish.Status != constants.WishStatusPending {
		t.Fatalf("expected pending status, got %s", wish.Status)
	}
	if !granted {
		t.Fatal("first wish badge should be granted")
	}
}

func TestWishService_Update_Forbidden(t *testing.T) {
	t.Parallel()
	wishRepo := &mockWishRepo{
		findByIDFn: func(id uint64) (*model.Wish, error) {
			return &model.Wish{ID: id, UserID: 99, Title: "别人", Content: "别人的心愿内容"}, nil
		},
	}
	svc := NewWishService(wishRepo, &mockClaimRepo{}, &mockBlessRepo{}, &mockUserRepo{}, &mockBadge{}, &mockAudit{}, testLogger())
	_, err := svc.Update(1, 5, dto.UpdateWishRequest{Title: "篡改"}, "127.0.0.1", "req-4")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeWishNotOwner {
		t.Fatalf("expected CodeWishNotOwner, got %v", err)
	}
}
