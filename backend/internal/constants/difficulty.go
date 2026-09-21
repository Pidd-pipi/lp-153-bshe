package constants

// 难度标签枚举。
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)

// ValidDifficulties 难度白名单。
var ValidDifficulties = []string{DifficultyEasy, DifficultyMedium, DifficultyHard}

// DifficultyText 难度文本。
func DifficultyText(d string) string {
	switch d {
	case DifficultyEasy:
		return "简单"
	case DifficultyMedium:
		return "中等"
	case DifficultyHard:
		return "困难"
	default:
		return "未知"
	}
}
