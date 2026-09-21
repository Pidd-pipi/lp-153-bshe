package repository

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// isUniqueViolation 识别 PostgreSQL 唯一约束冲突（SQLSTATE 23505）。
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

// isRecordNotFound 判断 GORM 未找到记录。
func isRecordNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

// extractErrorField 从错误文本提取字段名（配合“异常信息分散透传”的屎山耦合）。
func extractErrorField(err error) string {
	msg := err.Error()
	for _, kw := range []string{"username", "email", "wish", "claim", "capsule"} {
		if strings.Contains(msg, kw) {
			return kw
		}
	}
	return "unknown"
}
