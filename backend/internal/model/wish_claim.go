package model

import "time"

// WishClaim 心愿认领（圆梦人）实体。同一心愿仅允许一条有效认领（wish_id 唯一索引兜底并发）。
type WishClaim struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	WishID         uint64    `gorm:"uniqueIndex;not null" json:"wish_id"`
	UserID         uint64    `gorm:"index;not null" json:"user_id"`
	Progress       int       `gorm:"not null;default:0" json:"progress"`
	LatestNote     string    `gorm:"type:text" json:"latest_note"`
	Status         string    `gorm:"size:20;index;not null;default:claimed" json:"status"`
	MilestoneCount int       `gorm:"not null;default:0" json:"milestone_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
