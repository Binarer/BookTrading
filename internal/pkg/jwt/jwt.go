package jwt

import (
	"booktrading/internal/domain/token"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrInvalidClaims  = errors.New("invalid token claims")
	ErrInvalidRefresh = errors.New("invalid refresh token")
	ErrExpiredRefresh = errors.New("refresh token has expired")
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Service struct {
	tokenAuth   *jwtauth.JWTAuth
	refreshRepo token.Repository
	userRepo    interface {
		GetByID(id uint) (*user.User, error)
	}
}

func NewService(secretKey string, refreshRepo token.Repository, userRepo interface {
	GetByID(id uint) (*user.User, error)
}) *Service {
	return &Service{
		tokenAuth:   jwtauth.New("HS256", []byte(secretKey), nil),
		refreshRepo: refreshRepo,
		userRepo:    userRepo,
	}
}

// generateRefreshToken генерирует криптографически безопасный refresh token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Service) GenerateTokenPair(user *user.User) (*TokenPair, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if user.Login == "" {
		return nil, fmt.Errorf("invalid user login")
	}

	// Generate access token
	_, accessToken, err := s.tokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	_, refreshToken, err := s.tokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour * 7).Unix(), // 7 days
	})
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh token в базе данных
	if err := s.refreshRepo.Save(user.ID, refreshToken, time.Now().Add(24*time.Hour*7)); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(24 * time.Hour * 7).Unix(),
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := s.tokenAuth.Decode(tokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Проверяем обязательные поля
	userID, ok := token.Get("user_id")
	if !ok {
		return nil, ErrInvalidClaims
	}

	login, ok := token.Get("login")
	if !ok {
		return nil, ErrInvalidClaims
	}

	return map[string]interface{}{
		"user_id": userID,
		"login":   login,
	}, nil
}

func (s *Service) RefreshTokenPair(refreshToken string) (*TokenPair, error) {
	// Проверяем refresh token в базе данных
	userID, err := s.refreshRepo.Validate(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Удаляем использованный refresh token
	if err := s.refreshRepo.Delete(refreshToken); err != nil {
		logger.Error("Failed to delete used refresh token", err)
	}

	// Генерируем новую пару токенов
	return s.GenerateTokenPair(user)
}

func (s *Service) RevokeRefreshToken(refreshToken string) error {
	return s.refreshRepo.Delete(refreshToken)
}

func (s *Service) RevokeAllUserTokens(userID uint) error {
	return s.refreshRepo.DeleteUserTokens(userID)
}

func (s *Service) CleanupExpiredTokens() error {
	return s.refreshRepo.DeleteExpired()
}

func (s *Service) GetTokenAuth() *jwtauth.JWTAuth {
	return s.tokenAuth
}

// ValidateRefreshToken проверяет refresh token и возвращает пользователя
func (s *Service) ValidateRefreshToken(tokenString string) (*user.User, error) {
	token, err := s.tokenAuth.Decode(tokenString)
	if err != nil {
		return nil, err
	}

	userID, ok := token.Get("user_id")
	if !ok {
		return nil, ErrInvalidToken
	}

	return &user.User{
		ID: uint(userID.(float64)),
	}, nil
}
