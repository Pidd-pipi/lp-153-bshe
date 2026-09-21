package dto

import "github.com/wishwall/wishwall/internal/util"

// decodeStrings 解码 JSON 文本为字符串数组。
func decodeStrings(s string) []string { return util.DecodeStringArray(s) }
