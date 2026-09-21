package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// AuditLogRepository 审计日志仓储接口。
type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	List(filters map[string]any, offset, limit int) ([]model.AuditLog, error)
	Count(filters map[string]any) (int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 构造审计日志仓储。
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *model.AuditLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepository) List(filters map[string]any, offset, limit int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	q := r.db.Model(&model.AuditLog{})
	if v, ok := filters["user_id"]; ok && v.(uint64) > 0 {
		q = q.Where("user_id = ?", v)
	}
	if v, ok := filters["action"]; ok && v != "" {
		q = q.Where("action = ?", v)
	}
	if v, ok := filters["entity_type"]; ok && v != "" {
		q = q.Where("entity_type = ?", v)
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, nil
}

func (r *auditLogRepository) Count(filters map[string]any) (int64, error) {
	var total int64
	q := r.db.Model(&model.AuditLog{})
	if v, ok := filters["user_id"]; ok && v.(uint64) > 0 {
		q = q.Where("user_id = ?", v)
	}
	if v, ok := filters["action"]; ok && v != "" {
		q = q.Where("action = ?", v)
	}
	if v, ok := filters["entity_type"]; ok && v != "" {
		q = q.Where("entity_type = ?", v)
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count audit logs: %w", err)
	}
	return total, nil
}
