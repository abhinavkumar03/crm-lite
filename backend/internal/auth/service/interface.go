package service

import "context"

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	UserID      string
	Name        string
	Email       string
	AccessToken string
}

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	Profile(ctx context.Context, userID string) (*AuthResult, error)
}
