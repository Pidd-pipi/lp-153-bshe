package dto

import (
	"github.com/wishwall/wishwall/internal/model"
)

// SubmitAcceptanceRequest 圆梦人送交验收入参（说明必填）。
type SubmitAcceptanceRequest struct {
	Note string `json:"note" binding:"required,min=1,max=1000"`
}

// RejectAcceptanceRequest 发布者驳回验收入参（原因必填）。
type RejectAcceptanceRequest struct {
	Reason string `json:"reason" binding:"required,min=1,max=500"`
}

// WishAcceptanceResponse 验收记录返回结构。
type WishAcceptanceResponse struct {
	ID            uint64 `json:"id"`
	WishID        uint64 `json:"wish_id"`
	ClaimID       uint64 `json:"claim_id"`
	UserID        uint64 `json:"user_id"`
	FulfillerName string `json:"fulfiller_name"`
	Note          string `json:"note"`
	Status        string `json:"status"`
	RejectReason  string `json:"reject_reason"`
	SubmittedAt   string `json:"submitted_at"`
	ReviewedAt    string `json:"reviewed_at"`
	ReviewedBy    uint64 `json:"reviewed_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ToWishAcceptanceResponse 从模型构造返回结构。
func ToWishAcceptanceResponse(a *model.WishAcceptance) WishAcceptanceResponse {
	resp := WishAcceptanceResponse{
		ID:           a.ID,
		WishID:       a.WishID,
		ClaimID:      a.ClaimID,
		UserID:       a.UserID,
		Note:         a.Note,
		Status:       a.Status,
		RejectReason: a.RejectReason,
		ReviewedBy:   a.ReviewedBy,
		CreatedAt:    a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if a.SubmittedAt != nil {
		resp.SubmittedAt = a.SubmittedAt.Format("2006-01-02 15:04:05")
	}
	if a.ReviewedAt != nil {
		resp.ReviewedAt = a.ReviewedAt.Format("2006-01-02 15:04:05")
	}
	return resp
}
