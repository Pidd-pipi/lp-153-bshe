package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestWishClaimRepository_FindByIDForUpdate 验收入口必须以 SELECT ... FOR UPDATE
// 锁定认领行，配合状态二次校验保证并发验收只能成功一次。
func TestWishClaimRepository_FindByIDForUpdate(t *testing.T) {
	t.Parallel()
	gdb, mock := newMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "wish_id", "user_id", "progress", "status"}).
		AddRow(9, 5, 2, 100, "pending_acceptance")
	// 断言生成 SQL 带 FOR UPDATE 行锁。
	mock.ExpectQuery(regexp.QuoteMeta(`FOR UPDATE`)).
		WithArgs(uint64(9), 1).
		WillReturnRows(rows)

	repo := NewWishClaimRepository(gdb)
	claim, err := repo.FindByIDForUpdate(gdb, 9)
	if err != nil {
		t.Fatalf("find for update: %v", err)
	}
	if claim.Status != "pending_acceptance" {
		t.Fatalf("unexpected status %s", claim.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
