package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret     string
	Expiration time.Duration
}

type Service interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (uint, error)
}

type service struct {
	config Config
}

func NewService(config Config) Service {
	return &service{config: config}
}

func (s *service) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.config.Expiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *service) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return uint(sub), nil
		}
		return 0, errors.New("invalid subject claim")
	}

	return 0, errors.New("invalid token")
}
