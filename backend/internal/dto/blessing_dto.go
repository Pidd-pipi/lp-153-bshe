package dto

import (
	"github.com/wishwall/wishwall/internal/model"
)

// CreateBlessingRequest 发送祝福入参。
type CreateBlessingRequest struct {
	Content   string `json:"content" binding:"required,min=1,max=500"`
	GiftEmoji string `json:"gift_emoji" binding:"omitempty,max=50"`
}

// BlessingResponse 祝福返回结构。
type BlessingResponse struct {
	ID            uint64 `json:"id"`
	WishID        uint64 `json:"wish_id"`
	UserID        uint64 `json:"user_id"`
	SenderName    string `json:"sender_name"`
	Content       string `json:"content"`
	GiftEmoji     string `json:"gift_emoji"`
	IsCelebrating bool   `json:"is_celebrating"`
	CreatedAt     string `json:"created_at"`
}

// ToBlessingResponse 从模型构造返回结构。
func ToBlessingResponse(b *model.Blessing, senderName string) BlessingResponse {
	return BlessingResponse{
		ID:            b.ID,
		WishID:        b.WishID,
		UserID:        b.UserID,
		SenderName:    senderName,
		Content:       b.Content,
		GiftEmoji:     b.GiftEmoji,
		IsCelebrating: b.IsCelebrating,
		CreatedAt:     b.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
