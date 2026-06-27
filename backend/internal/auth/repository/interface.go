package repository

import "context"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}
