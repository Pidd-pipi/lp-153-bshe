package util

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wishwall/wishwall/internal/constants"
)

// FormatTime 格式化时间（date 风格）。
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 格式化时间（datetime 风格）。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatWishStatus 心愿状态文本（耦合 constants.WishStatusText）。
func FormatWishStatus(status string) string { return constants.WishStatusText(status) }

// FormatVisibility 可见范围文本（耦合 constants.VisibilityText）。
func FormatVisibility(v string) string { return constants.VisibilityText(v) }

// FormatDifficulty 难度文本（耦合 constants.DifficultyText）。
func FormatDifficulty(d string) string { return constants.DifficultyText(d) }

// FormatCategory 分类文本（耦合 constants.CategoryText）。
func FormatCategory(c string) string { return constants.CategoryText(c) }

// FormatCapsuleStatus 胶囊状态文本（耦合 constants.CapsuleStatusText）。
func FormatCapsuleStatus(s string) string { return constants.CapsuleStatusText(s) }

// FormatBadgeType 徽章类型文本（耦合 constants.BadgeTypeText）。
func FormatBadgeType(t string) string { return constants.BadgeTypeText(t) }

// FormatRole 角色文本。
func FormatRole(role string) string {
	if role == constants.RoleAdmin {
		return "管理员"
	}
	return "普通用户"
}

// TruncateString 截断字符串。
func TruncateString(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

// EncodeStringArray 将 []string 编码为 JSON 文本存储。
func EncodeStringArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(items)
	return string(b)
}

// DecodeStringArray 将 JSON 文本解码为 []string。
func DecodeStringArray(s string) []string {
	if s == "" {
		return []string{}
	}
	var items []string
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return []string{}
	}
	return items
}

// BuildUploadURL 拼接文件访问 URL（耦合 MinIO 配置）。
func BuildUploadURL(baseURL, bucket, objectKey string) string {
	base := strings.TrimRight(baseURL, "/")
	key := strings.TrimLeft(objectKey, "/")
	return base + "/" + bucket + "/" + key
}

// ParseTags 把逗号/空格分隔的标签字符串拆成数组。
func ParseTags(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '，' || r == ' ' || r == '\n'
	})
	var out []string
	seen := map[string]bool{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}
