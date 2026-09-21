package constants

// 时光胶囊状态机：locked -> unlocked。
const (
	CapsuleStatusLocked   = "locked"
	CapsuleStatusUnlocked = "unlocked"
)

// ValidCapsuleStatuses 胶囊状态白名单。
var ValidCapsuleStatuses = []string{CapsuleStatusLocked, CapsuleStatusUnlocked}

// CapsuleStatusText 胶囊状态文本。
func CapsuleStatusText(s string) string {
	if s == CapsuleStatusUnlocked {
		return "已解锁"
	}
	return "未解锁"
}
