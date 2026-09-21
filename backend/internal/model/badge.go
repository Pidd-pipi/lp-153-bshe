package model

import "time"

// Badge 成就徽章实体，同一用户同一类型仅一枚（唯一索引兜底并发）。
type Badge struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	UserID      uint64    `gorm:"index:idx_user_type,unique;not null" json:"user_id"`
	Type        string    `gorm:"size:50;index:idx_user_type,unique;not null" json:"type"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"size:255" json:"description"`
	Icon        string    `gorm:"size:50" json:"icon"`
	EarnedAt    time.Time `json:"earned_at"`
}
