package model

import "time"

// Blessing 祝福留言实体，挂在心愿留言板上，可携带虚拟礼物表情。
type Blessing struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	WishID        uint64    `gorm:"index;not null" json:"wish_id"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	GiftEmoji     string    `gorm:"size:50" json:"gift_emoji"`
	IsCelebrating bool      `gorm:"not null;default:false" json:"is_celebrating"`
	CreatedAt     time.Time `json:"created_at"`
}
