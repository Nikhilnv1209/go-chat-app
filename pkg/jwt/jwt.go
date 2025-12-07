package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	Secret     string
	Expiration time.Duration
}

type Service interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type service struct {
	config Config
}

func NewService(config Config) Service {
	return &service{config: config}
}

func (s *service) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(s.config.Expiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *service) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			id, err := uuid.Parse(sub)
			if err != nil {
				return uuid.Nil, errors.New("invalid subject UUID")
			}
			return id, nil
		}
		return uuid.Nil, errors.New("invalid subject claim")
	}

	return uuid.Nil, errors.New("invalid token")
}
