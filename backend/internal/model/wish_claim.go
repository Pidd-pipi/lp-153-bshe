package model

import "time"

// WishClaim 心愿认领（圆梦人）实体。同一心愿仅允许一条有效认领（wish_id 唯一索引兜底并发）。
// 状态机：claimed -> in_progress -> pending_acceptance -> completed；
// 待验收被发布者驳回时 pending_acceptance -> in_progress（保留送交说明，可重新送交）。
type WishClaim struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	WishID         uint64 `gorm:"uniqueIndex;not null" json:"wish_id"`
	UserID         uint64 `gorm:"index;not null" json:"user_id"`
	Progress       int    `gorm:"not null;default:0" json:"progress"`
	LatestNote     string `gorm:"type:text" json:"latest_note"`
	Status         string `gorm:"size:20;index;not null;default:claimed" json:"status"`
	MilestoneCount int    `gorm:"not null;default:0" json:"milestone_count"`
	// SubmissionNote 圆梦人送交验收时填写的说明；驳回后保留，重新送交时覆盖。
	SubmissionNote string `gorm:"type:text;not null;default:''" json:"submission_note"`
	// SubmittedAt 最近一次送交验收的时间。
	SubmittedAt *time.Time `json:"submitted_at"`
	// RejectReason 发布者最近一次驳回原因；重新送交时清空。
	RejectReason string `gorm:"type:text;not null;default:''" json:"reject_reason"`
	// ReviewedAt 发布者最近一次验收（确认/驳回）时间。
	ReviewedAt *time.Time `json:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
