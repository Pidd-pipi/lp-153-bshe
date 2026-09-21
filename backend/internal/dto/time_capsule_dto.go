package dto

import (
	"time"

	"github.com/wishwall/wishwall/internal/model"
)

// CreateCapsuleRequest 创建时光胶囊入参。
type CreateCapsuleRequest struct {
	Title     string    `json:"title" binding:"required,min=2,max=100"`
	Content   string    `json:"content" binding:"required,min=5,max=2000"`
	ImageURLs []string  `json:"image_urls"`
	AudioURL  string    `json:"audio_url" binding:"omitempty,max=255"`
	UnlockAt  time.Time `json:"unlock_at" binding:"required"`
}

// CapsuleResponse 胶囊返回结构。未解锁时 content 置空（内容不可见）。
type CapsuleResponse struct {
	ID         uint64   `json:"id"`
	UserID     uint64   `json:"user_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	ImageURLs  []string `json:"image_urls"`
	AudioURL   string   `json:"audio_url"`
	UnlockAt   string   `json:"unlock_at"`
	Status     string   `json:"status"`
	UnlockedAt *string  `json:"unlocked_at"`
	CreatedAt  string   `json:"created_at"`
}

// ToCapsuleResponse 从模型构造返回结构，maskContent 控制未解锁内容是否打码。
func ToCapsuleResponse(c *model.TimeCapsule, maskContent bool) CapsuleResponse {
	resp := CapsuleResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Title:     c.Title,
		AudioURL:  c.AudioURL,
		UnlockAt:  c.UnlockAt.Format("2006-01-02 15:04:05"),
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if c.UnlockedAt != nil {
		s := c.UnlockedAt.Format("2006-01-02 15:04:05")
		resp.UnlockedAt = &s
	}
	if maskContent {
		resp.Content = ""
		resp.ImageURLs = []string{}
	} else {
		resp.Content = c.Content
		resp.ImageURLs = decodeStrings(c.ImageURLs)
	}
	return resp
}
