package model

import "time"

// WishAcceptance 心愿验收记录。同一心愿仅允许一条（wish_id 唯一索引兜底并发送交）。
// 状态机：pending(待验收) -> confirmed(验收通过) / rejected(已驳回，可重新送交)。
type WishAcceptance struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	WishID       uint64     `gorm:"uniqueIndex;not null" json:"wish_id"`
	ClaimID      uint64     `gorm:"index;not null" json:"claim_id"`
	UserID       uint64     `gorm:"index;not null" json:"user_id"` // 圆梦人
	Note         string     `gorm:"type:text" json:"note"`         // 送交说明（驳回后保留）
	Status       string     `gorm:"size:20;index;not null;default:pending" json:"status"`
	RejectReason string     `gorm:"type:text" json:"reject_reason"` // 驳回原因（驳回必填）
	SubmittedAt  *time.Time `json:"submitted_at"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
	ReviewedBy   uint64     `gorm:"not null;default:0" json:"reviewed_by"` // 验收的发布者
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
