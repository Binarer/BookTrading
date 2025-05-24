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
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrInvalidClaims  = errors.New("invalid token claims")
	ErrInvalidRefresh = errors.New("invalid refresh token")
	ErrExpiredRefresh = errors.New("refresh token has expired")
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Service struct {
	secretKey     []byte
	refreshSecret []byte
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	refreshRepo   token.Repository
	userRepo      interface {
		GetByID(id uint) (*user.User, error)
	}
}

func NewService(secretKey string, refreshSecret string, issuer string, refreshRepo token.Repository, userRepo interface {
	GetByID(id uint) (*user.User, error)
}) *Service {
	if len(secretKey) < 32 {
		logger.Error("Secret key is too short", fmt.Errorf("minimum length is 32 bytes, got %d", len(secretKey)))
	}
	if len(refreshSecret) < 32 {
		logger.Error("Refresh secret key is too short", fmt.Errorf("minimum length is 32 bytes, got %d", len(refreshSecret)))
	}
	return &Service{
		secretKey:     []byte(secretKey),
		refreshSecret: []byte(refreshSecret),
		issuer:        issuer,
		accessTTL:     24 * time.Hour,
		refreshTTL:    30 * 24 * time.Hour, // 30 дней
		refreshRepo:   refreshRepo,
		userRepo:      userRepo,
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

	now := time.Now()
	expiresAt := now.Add(s.accessTTL)

	// Генерируем access token
	claims := Claims{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.secretKey)
	if err != nil {
		logger.Error("Failed to sign access token", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		logger.Error("Failed to generate refresh token", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем refresh token в базе данных
	if err := s.refreshRepo.Save(user.ID, refreshToken, now.Add(s.refreshTTL)); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresAt.Unix(),
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		logger.Error("Token validation failed", err)
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Проверяем обязательные поля
	if claims.UserID == 0 {
		return nil, ErrInvalidClaims
	}

	if claims.Login == "" {
		return nil, ErrInvalidClaims
	}

	// Проверяем issuer
	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
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

func (s *Service) GetJWTAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", s.secretKey, nil)
}

// ValidateRefreshToken проверяет refresh token и возвращает пользователя
func (s *Service) ValidateRefreshToken(token string) (*user.User, error) {
	// Проверяем подпись токена
	claims := &jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Получаем ID пользователя из токена
	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Получаем пользователя из базы данных
	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
