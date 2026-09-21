package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// gormclauseLock 返回 SELECT ... FOR UPDATE 行级锁子句（并发认领）。
func gormclauseLock() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}

// applyWishFilters 按筛选条件组装查询。
func applyWishFilters(q *gorm.DB, filters map[string]any) *gorm.DB {
	if v, ok := filters["keyword"]; ok && v != "" {
		kw := "%" + v.(string) + "%"
		q = q.Where("title ILIKE ? OR content ILIKE ?", kw, kw)
	}
	if v, ok := filters["category"]; ok && v != "" {
		q = q.Where("category = ?", v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		q = q.Where("status = ?", v)
	}
	if v, ok := filters["visibility"]; ok && v != "" {
		q = q.Where("visibility = ?", v)
	}
	if v, ok := filters["user_id"]; ok && v.(uint64) > 0 {
		q = q.Where("user_id = ?", v)
	}
	return q
}
