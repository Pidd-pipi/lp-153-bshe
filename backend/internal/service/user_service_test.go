package service

import (
	"context"
	"errors"
	"testing"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		repo        *mockUserRepo
		wantErrCode int
	}{
		{
			name: "register success",
			repo: &mockUserRepo{
				findByUsernameFn: func(username string) (*model.User, error) {
					return nil, repository.ErrNotFound
				},
				findByEmailFn: func(email string) (*model.User, error) {
					return nil, repository.ErrNotFound
				},
				createFn: func(user *model.User) error {
					user.ID = 1
					return nil
				},
			},
			wantErrCode: 0,
		},
		{
			name: "username already taken",
			repo: &mockUserRepo{
				findByUsernameFn: func(username string) (*model.User, error) {
					return &model.User{Username: username}, nil
				},
			},
			wantErrCode: constants.CodeUserExists,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUserService(tt.repo, testConfig(), &mockAudit{}, testLogger())
			user, err := svc.Register(context.Background(), dto.RegisterRequest{
				Username: "bob", Email: "bob@example.com", Password: "secret123",
			}, "127.0.0.1", "req-1")
			if tt.wantErrCode == 0 {
				if err != nil {
					t.Fatalf("register should succeed, got %v", err)
				}
				if user.Username != "bob" {
					t.Fatalf("unexpected username %s", user.Username)
				}
				if user.Role != constants.RoleUser {
					t.Fatalf("unexpected role %s", user.Role)
				}
				return
			}
			var appErr *util.AppError
			if !errors.As(err, &appErr) || appErr.Code != tt.wantErrCode {
				t.Fatalf("expected code %d, got %v", tt.wantErrCode, err)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()
	hash, _ := util.HashPassword("secret123")
	user := &model.User{ID: 2, Username: "alice", Email: "alice@example.com", PasswordHash: hash, Role: constants.RoleUser, Status: "active"}
	tests := []struct {
		name        string
		account     string
		password    string
		findFn      func(account string) (*model.User, error)
		wantErrCode int
		wantToken   bool
	}{
		{
			name:      "login success by username",
			account:   "alice",
			password:  "secret123",
			findFn:    func(account string) (*model.User, error) { return user, nil },
			wantToken: true,
		},
		{
			name:        "wrong password",
			account:     "alice",
			password:    "wrong-pass",
			findFn:      func(account string) (*model.User, error) { return user, nil },
			wantErrCode: constants.CodeInvalidCredential,
		},
		{
			name:        "user not found",
			account:     "ghost",
			password:    "secret123",
			findFn:      func(account string) (*model.User, error) { return nil, repository.ErrNotFound },
			wantErrCode: constants.CodeInvalidCredential,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockUserRepo{findByAccountFn: tt.findFn}
			svc := NewUserService(repo, testConfig(), &mockAudit{}, testLogger())
			resp, err := svc.Login(context.Background(), dto.LoginRequest{Account: tt.account, Password: tt.password}, "127.0.0.1", "req-2")
			if tt.wantErrCode != 0 {
				var appErr *util.AppError
				if !errors.As(err, &appErr) || appErr.Code != tt.wantErrCode {
					t.Fatalf("expected code %d, got %v", tt.wantErrCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("login should succeed, got %v", err)
			}
			if tt.wantToken && resp.Token == "" {
				t.Fatal("expected non-empty token")
			}
			if resp.User.ID != user.ID {
				t.Fatalf("expected user id %d, got %d", user.ID, resp.User.ID)
			}
		})
	}
}
