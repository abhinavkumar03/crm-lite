package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/repository"
)

type AuthService struct {
	repository *repository.AuthRepository
}

func New(
	repository *repository.AuthRepository,
) *AuthService {

	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) hashPassword(
	password string,
) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *AuthService) verifyPassword(
	hash string,
	password string,
) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*dto.UserResponse, error) {

	existingUser, err := s.repository.GetUserByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &auth.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	err = s.repository.CreateUser(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*auth.User, error) {

	user, err := s.repository.GetUserByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GetProfile(
	ctx context.Context,
	userID string,
) (*dto.UserResponse, error) {

	user, err := s.repository.GetUserByID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
