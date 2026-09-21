package dto

import (
	"github.com/wishwall/wishwall/internal/model"
)

// UpdateProgressRequest 更新圆梦进度入参（百分比+文字+里程碑打卡）。
// 待验收（pending_acceptance）期间禁止更新进度，由 service 状态机拦截。
type UpdateProgressRequest struct {
	Progress    int    `json:"progress" binding:"required,min=0,max=100"`
	Note        string `json:"note" binding:"omitempty,max=1000"`
	IsMilestone bool   `json:"is_milestone"`
}

// SubmitClaimRequest 送交验收入参：圆梦人填写送交说明后进入待验收。
type SubmitClaimRequest struct {
	Note string `json:"note" binding:"required,min=2,max=1000"`
}

// ReviewClaimRequest 发布者验收入参：action=approve 确认完成；action=reject 驳回且 reason 必填。
type ReviewClaimRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Reason string `json:"reason" binding:"omitempty,max=1000"`
}

// WishClaimResponse 认领记录返回结构（含验收状态/送交说明/驳回原因/验收时间）。
type WishClaimResponse struct {
	ID             uint64  `json:"id"`
	WishID         uint64  `json:"wish_id"`
	WishTitle      string  `json:"wish_title"`
	UserID         uint64  `json:"user_id"`
	FulfillerName  string  `json:"fulfiller_name"`
	Progress       int     `json:"progress"`
	LatestNote     string  `json:"latest_note"`
	Status         string  `json:"status"`
	MilestoneCount int     `json:"milestone_count"`
	SubmissionNote string  `json:"submission_note"`
	SubmittedAt    *string `json:"submitted_at"`
	RejectReason   string  `json:"reject_reason"`
	ReviewedAt     *string `json:"reviewed_at"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ToWishClaimResponse 从模型构造返回结构。
func ToWishClaimResponse(c *model.WishClaim, wishTitle, fulfillerName string) WishClaimResponse {
	resp := WishClaimResponse{
		ID:             c.ID,
		WishID:         c.WishID,
		WishTitle:      wishTitle,
		UserID:         c.UserID,
		FulfillerName:  fulfillerName,
		Progress:       c.Progress,
		LatestNote:     c.LatestNote,
		Status:         c.Status,
		MilestoneCount: c.MilestoneCount,
		SubmissionNote: c.SubmissionNote,
		RejectReason:   c.RejectReason,
		CreatedAt:      c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if c.SubmittedAt != nil {
		s := c.SubmittedAt.Format("2006-01-02 15:04:05")
		resp.SubmittedAt = &s
	}
	if c.ReviewedAt != nil {
		s := c.ReviewedAt.Format("2006-01-02 15:04:05")
		resp.ReviewedAt = &s
	}
	return resp
}
