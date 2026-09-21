package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWishRepository_Count(t *testing.T) {
	t.Parallel()
	gdb, mock := newMockDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "wishes"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	repo := NewWishRepository(gdb)
	total, err := repo.Count(map[string]any{"status": "pending", "keyword": "旅行"})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected 3, got %d", total)
	}
}

func TestWishRepository_List(t *testing.T) {
	t.Parallel()
	gdb, mock := newMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "title", "status"}).
		AddRow(1, "去冰岛看极光", "pending").
		AddRow(2, "学会吉他", "in_progress")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "wishes"`)).
		WillReturnRows(rows)
	repo := NewWishRepository(gdb)
	items, err := repo.List(map[string]any{"sort": ""}, "hot", 0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}
