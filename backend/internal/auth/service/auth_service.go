package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/entity"
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
) (*dto.ProfileResponse, error) {

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

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	err = s.repository.Create(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*entity.User, error) {

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
) (*entity.User, error) {

	return s.repository.FindByID(
		ctx,
		id,
	)
}
