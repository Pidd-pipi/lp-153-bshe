package model

// FulfillerStat 圆梦人排行榜统计行（排行榜与成就复用）。
type FulfillerStat struct {
	UserID         uint64 `gorm:"column:user_id" json:"user_id"`
	Nickname       string `gorm:"column:nickname" json:"nickname"`
	Avatar         string `gorm:"column:avatar" json:"avatar"`
	CompletedCount int64  `gorm:"column:completed_count" json:"completed_count"`
}
