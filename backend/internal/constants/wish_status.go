package constants

// 心愿状态机：pending -> claimed -> in_progress -> pending_acceptance -> completed。
// 验收驳回：pending_acceptance -> in_progress（圆梦人保留送交说明，可修改后重新送交）。
const (
	WishStatusPending           = "pending"
	WishStatusClaimed           = "claimed"
	WishStatusInProgress        = "in_progress"
	WishStatusPendingAcceptance = "pending_acceptance"
	WishStatusCompleted         = "completed"
)

// ValidWishStatuses 状态筛选白名单。
var ValidWishStatuses = []string{
	WishStatusPending,
	WishStatusClaimed,
	WishStatusInProgress,
	WishStatusPendingAcceptance,
	WishStatusCompleted,
}

// WishStatusText 状态文本（formatters 亦引用）。
func WishStatusText(status string) string {
	switch status {
	case WishStatusPending:
		return "待认领"
	case WishStatusClaimed:
		return "已被认领"
	case WishStatusInProgress:
		return "圆梦中"
	case WishStatusPendingAcceptance:
		return "待验收"
	case WishStatusCompleted:
		return "已完成"
	default:
		return "未知"
	}
}
