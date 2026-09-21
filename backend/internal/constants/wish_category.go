package constants

// 心愿分类枚举。
const (
	CategoryStudy   = "study"
	CategoryTravel  = "travel"
	CategoryEmotion = "emotion"
	CategoryCareer  = "career"
	CategoryLife    = "life"
	CategoryOther   = "other"
)

// ValidCategories 分类白名单。
var ValidCategories = []string{CategoryStudy, CategoryTravel, CategoryEmotion, CategoryCareer, CategoryLife, CategoryOther}

// CategoryText 分类文本。
func CategoryText(c string) string {
	switch c {
	case CategoryStudy:
		return "学习成长"
	case CategoryTravel:
		return "旅行探险"
	case CategoryEmotion:
		return "情感陪伴"
	case CategoryCareer:
		return "职业发展"
	case CategoryLife:
		return "生活小确幸"
	case CategoryOther:
		return "其他"
	default:
		return "其他"
	}
}
