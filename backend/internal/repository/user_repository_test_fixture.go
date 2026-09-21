package repository

import "github.com/wishwall/wishwall/internal/model"

// userFixture 测试用户样本。
var userFixture = model.User{
	Username: "alice", Email: "alice@example.com",
	PasswordHash: "hash", Nickname: "Alice",
}
