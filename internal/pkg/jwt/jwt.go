package jwt

import (
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

type Service struct {
	secretKey []byte
}

func NewService(secretKey string) *Service {
	return &Service{
		secretKey: []byte(secretKey),
	}
}

func (s *Service) GenerateToken(user *user.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	logger.Info("Starting token validation")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		logger.Info("Using secret key for validation")
		return s.secretKey, nil
	})

	if err != nil {
		logger.Error("Token validation failed", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		logger.Info(fmt.Sprintf("Token validation successful for user ID: %d", claims.UserID))
		return claims, nil
	}

	logger.Error("Token claims validation failed", nil)
	return nil, ErrInvalidToken
}
