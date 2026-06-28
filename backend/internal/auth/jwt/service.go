package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{
		secret: secret,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwtv5.RegisteredClaims
}

func (j *Service) GenerateToken(
	userID string,
	email string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp": time.Now().
			Add(24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(j.secret),
	)
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {

	token, err := jwtv5.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtv5.Token) (interface{}, error) {
			return []byte(s.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
