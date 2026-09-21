package model

import "time"

// AuditLog 操作审计日志实体（横切关注点）。
type AuditLog struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:50;index;not null" json:"action"`
	EntityType string    `gorm:"size:50" json:"entity_type"`
	EntityID   string    `gorm:"size:64" json:"entity_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	IP         string    `gorm:"size:64" json:"ip"`
	RequestID  string    `gorm:"size:64;index" json:"request_id"`
	CreatedAt  time.Time `json:"created_at"`
}
