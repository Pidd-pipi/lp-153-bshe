package constants

// 角色枚举：JWT + RBAC 依赖该枚举。
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// ValidRoles 用于 DTO/handler 校验。
var ValidRoles = []string{RoleUser, RoleAdmin}
