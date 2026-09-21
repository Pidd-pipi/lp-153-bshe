package dto

import "github.com/wishwall/wishwall/internal/model"

// RegisterRequest 注册入参。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=120"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=50"`
}

// LoginRequest 登录入参。
type LoginRequest struct {
	Account  string `json:"account" binding:"required,max=120"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// UpdateProfileRequest 更新个人资料入参。
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Bio      string `json:"bio" binding:"omitempty,max=500"`
}

// UserResponse 用户返回结构。
type UserResponse struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ToUserResponse 从模型构造返回结构。
func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// LoginResponse 登录返回（含 JWT）。
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
