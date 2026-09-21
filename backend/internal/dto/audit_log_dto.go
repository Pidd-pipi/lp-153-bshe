package dto

import (
	"github.com/wishwall/wishwall/internal/model"
)

// AuditLogQuery 审计日志查询参数。
type AuditLogQuery struct {
	PageQuery
	UserID     uint64 `form:"user_id" json:"user_id"`
	Action     string `form:"action" json:"action"`
	EntityType string `form:"entity_type" json:"entity_type"`
}

// AuditLogResponse 审计日志返回结构。
type AuditLogResponse struct {
	ID         uint64 `json:"id"`
	UserID     uint64 `json:"user_id"`
	Username   string `json:"username"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Detail     string `json:"detail"`
	IP         string `json:"ip"`
	RequestID  string `json:"request_id"`
	CreatedAt  string `json:"created_at"`
}

// ToAuditLogResponse 从模型构造返回结构。
func ToAuditLogResponse(a *model.AuditLog, username string) AuditLogResponse {
	return AuditLogResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		Username:   username,
		Action:     a.Action,
		EntityType: a.EntityType,
		EntityID:   a.EntityID,
		Detail:     a.Detail,
		IP:         a.IP,
		RequestID:  a.RequestID,
		CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
