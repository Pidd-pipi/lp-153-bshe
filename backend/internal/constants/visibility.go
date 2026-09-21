package constants

// 可见范围枚举：公开/好友/匿名。
const (
	VisibilityPublic    = "public"
	VisibilityFriend    = "friend"
	VisibilityAnonymous = "anonymous"
)

// ValidVisibilities 可见范围白名单。
var ValidVisibilities = []string{VisibilityPublic, VisibilityFriend, VisibilityAnonymous}

// VisibilityText 可见范围文本。
func VisibilityText(v string) string {
	switch v {
	case VisibilityPublic:
		return "公开"
	case VisibilityFriend:
		return "好友可见"
	case VisibilityAnonymous:
		return "匿名"
	default:
		return "未知"
	}
}
