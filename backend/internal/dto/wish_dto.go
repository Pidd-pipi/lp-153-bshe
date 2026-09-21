package dto

import (
	"time"

	"github.com/wishwall/wishwall/internal/model"
)

// CreateWishRequest 发布心愿入参。
type CreateWishRequest struct {
	Title           string    `json:"title" binding:"required,min=2,max=100"`
	Content         string    `json:"content" binding:"required,min=5,max=2000"`
	ImageURLs       []string  `json:"image_urls"`
	Category        string    `json:"category" binding:"required,oneof=study travel emotion career life other"`
	Visibility      string    `json:"visibility" binding:"required,oneof=public friend anonymous"`
	Difficulty      string    `json:"difficulty" binding:"required,oneof=easy medium hard"`
	ExpectedDeadline *time.Time `json:"expected_deadline"`
	IsAnonymous     bool      `json:"is_anonymous"`
}

// UpdateWishRequest 更新心愿入参（仅心愿发布者可用）。
type UpdateWishRequest struct {
	Title           string    `json:"title" binding:"omitempty,min=2,max=100"`
	Content         string    `json:"content" binding:"omitempty,min=5,max=2000"`
	ImageURLs       []string  `json:"image_urls"`
	Category        string    `json:"category" binding:"omitempty,oneof=study travel emotion career life other"`
	Visibility      string    `json:"visibility" binding:"omitempty,oneof=public friend anonymous"`
	Difficulty      string    `json:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	ExpectedDeadline *time.Time `json:"expected_deadline"`
	IsAnonymous     *bool     `json:"is_anonymous"`
}

// WishQuery 心愿列表查询参数。
type WishQuery struct {
	PageQuery
	Keyword    string `form:"keyword" json:"keyword"`
	Category   string `form:"category" json:"category"`
	Status     string `form:"status" json:"status"`
	Visibility string `form:"visibility" json:"visibility"`
	Sort       string `form:"sort" json:"sort" binding:"omitempty,oneof=hot new"`
	UserID     uint64 `form:"user_id" json:"user_id"`
}

// WishResponse 心愿返回结构。
type WishResponse struct {
	ID               uint64    `json:"id"`
	UserID           uint64    `json:"user_id"`
	AuthorNickname   string    `json:"author_nickname"`
	AuthorAvatar     string    `json:"author_avatar"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	ImageURLs        []string  `json:"image_urls"`
	Category         string    `json:"category"`
	Visibility       string    `json:"visibility"`
	Difficulty       string    `json:"difficulty"`
	ExpectedDeadline *string   `json:"expected_deadline"`
	Status           string    `json:"status"`
	LikesCount       int       `json:"likes_count"`
	CompletionNote   string    `json:"completion_note"`
	IsAnonymous      bool      `json:"is_anonymous"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

// ToWishResponse 从模型构造返回结构。
func ToWishResponse(w *model.Wish) WishResponse {
	resp := WishResponse{
		ID:             w.ID,
		UserID:         w.UserID,
		Title:          w.Title,
		Content:        w.Content,
		ImageURLs:      decodeStrings(w.ImageURLs),
		Category:       w.Category,
		Visibility:     w.Visibility,
		Difficulty:     w.Difficulty,
		Status:         w.Status,
		LikesCount:     w.LikesCount,
		CompletionNote: w.CompletionNote,
		IsAnonymous:    w.IsAnonymous,
		CreatedAt:      w.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      w.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if w.ExpectedDeadline != nil {
		s := w.ExpectedDeadline.Format("2006-01-02")
		resp.ExpectedDeadline = &s
	}
	return resp
}

// WishDetailResponse 心愿详情返回结构（含认领摘要、验收记录与祝福数）。
type WishDetailResponse struct {
	WishResponse
	Claim         *WishClaimResponse      `json:"claim"`
	Acceptance    *WishAcceptanceResponse `json:"acceptance"`
	BlessingCount int64                   `json:"blessing_count"`
}
