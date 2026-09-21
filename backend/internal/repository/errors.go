package repository

import "errors"

// 仓储层哨兵错误，上层用 errors.Is 判断。
var (
	ErrNotFound   = errors.New("record not found")
	ErrConflict   = errors.New("record conflict")
	ErrDuplicate  = errors.New("duplicate record")
)
