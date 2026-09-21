package constants

// 成就徽章类型枚举。
const (
	BadgeTypeFirstWish      = "first_wish"
	BadgeTypeFirstClaim     = "first_claim"
	BadgeTypeFirstBlessing  = "first_blessing"
	BadgeTypeTenCompletions = "ten_completions"
	BadgeTypeWishMaster     = "wish_master"
)

// ValidBadgeTypes 徽章类型白名单。
var ValidBadgeTypes = []string{
	BadgeTypeFirstWish,
	BadgeTypeFirstClaim,
	BadgeTypeFirstBlessing,
	BadgeTypeTenCompletions,
	BadgeTypeWishMaster,
}

// BadgeTypeText 徽章类型文本。
func BadgeTypeText(t string) string {
	switch t {
	case BadgeTypeFirstWish:
		return "首次许愿"
	case BadgeTypeFirstClaim:
		return "首次认领"
	case BadgeTypeFirstBlessing:
		return "首次祝福"
	case BadgeTypeTenCompletions:
		return "十次圆梦"
	case BadgeTypeWishMaster:
		return "圆梦大师"
	default:
		return "神秘徽章"
	}
}
