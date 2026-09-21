package model

import "time"

// Wish 心愿实体。状态机：pending -> claimed -> in_progress -> completed。
type Wish struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	UserID          uint64     `gorm:"index;not null" json:"user_id"`
	Title           string     `gorm:"size:100;not null" json:"title"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	ImageURLs       string     `gorm:"type:text;not null;default:'[]'" json:"-"`
	Category        string     `gorm:"size:30;index;not null;default:other" json:"category"`
	Visibility      string     `gorm:"size:20;index;not null;default:public" json:"visibility"`
	Difficulty      string     `gorm:"size:20;not null;default:medium" json:"difficulty"`
	ExpectedDeadline *time.Time `json:"expected_deadline"`
	Status          string     `gorm:"size:20;index;not null;default:pending" json:"status"`
	LikesCount      int        `gorm:"not null;default:0" json:"likes_count"`
	CompletionNote  string     `gorm:"type:text" json:"completion_note"`
	IsAnonymous     bool       `gorm:"not null;default:false" json:"is_anonymous"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
