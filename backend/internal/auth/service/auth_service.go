package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
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
) (*auth.User, error) {

	exists, err := s.repository.ExistsByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	hash, err := s.hashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	user := &auth.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	}

	err = s.repository.Create(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*auth.User, error) {

	user, err := s.repository.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = s.verifyPassword(
		user.PasswordHash,
		password,
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.PasswordHash = ""

	return user, nil
}

func (s *AuthService) Profile(
	ctx context.Context,
	id string,
) (*auth.User, error) {

	return s.repository.FindByID(
		ctx,
		id,
	)
}
