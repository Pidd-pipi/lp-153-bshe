package dto

import (
	"github.com/wishwall/wishwall/internal/model"
)

// BadgeResponse 徽章返回结构。
type BadgeResponse struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	EarnedAt    string `json:"earned_at"`
}

// ToBadgeResponse 从模型构造返回结构。
func ToBadgeResponse(b *model.Badge) BadgeResponse {
	return BadgeResponse{
		ID:          b.ID,
		UserID:      b.UserID,
		Type:        b.Type,
		Title:       b.Title,
		Description: b.Description,
		Icon:        b.Icon,
		EarnedAt:    b.EarnedAt.Format("2006-01-02 15:04:05"),
	}
}
