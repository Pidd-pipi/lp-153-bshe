package constants

// 验收状态机：pending(待验收) -> confirmed(验收通过) / rejected(已驳回，可重新送交)。
const (
	AcceptanceStatusPending   = "pending"
	AcceptanceStatusConfirmed = "confirmed"
	AcceptanceStatusRejected  = "rejected"
)

// ValidAcceptanceStatuses 验收状态白名单。
var ValidAcceptanceStatuses = []string{
	AcceptanceStatusPending,
	AcceptanceStatusConfirmed,
	AcceptanceStatusRejected,
}

// AcceptanceStatusText 验收状态文本（formatters 亦引用）。
func AcceptanceStatusText(status string) string {
	switch status {
	case AcceptanceStatusPending:
		return "待验收"
	case AcceptanceStatusConfirmed:
		return "验收通过"
	case AcceptanceStatusRejected:
		return "已驳回"
	default:
		return "未知"
	}
}
