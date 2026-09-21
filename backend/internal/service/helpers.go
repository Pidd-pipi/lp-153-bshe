package service

import "strconv"

// u64str 无符号整数转字符串。
func u64str(v uint64) string { return strconv.FormatUint(v, 10) }

// itoa 整数转字符串。
func itoa(v int) string { return strconv.Itoa(v) }
