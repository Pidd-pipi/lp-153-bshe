package model

import "time"

// TimeCapsule 时光胶囊实体。状态机：locked -> unlocked（到解锁时间自动解锁）。
type TimeCapsule struct {
	ID         uint64     `gorm:"primaryKey" json:"id"`
	UserID     uint64     `gorm:"index;not null" json:"user_id"`
	Title      string     `gorm:"size:100;not null" json:"title"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	ImageURLs  string     `gorm:"type:text;not null;default:'[]'" json:"-"`
	AudioURL   string     `gorm:"size:255" json:"audio_url"`
	UnlockAt   time.Time  `gorm:"not null" json:"unlock_at"`
	Status     string     `gorm:"size:20;index;not null;default:locked" json:"status"`
	UnlockedAt *time.Time `json:"unlocked_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
