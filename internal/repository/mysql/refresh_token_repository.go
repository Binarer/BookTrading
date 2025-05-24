package mysql

import (
	"booktrading/internal/pkg/logger"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(userID uint, token string, expiresAt time.Time) error {
	refreshToken := &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(refreshToken).Error; err != nil {
		logger.Error("Failed to save refresh token", err)
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) Validate(token string) (uint, error) {
	var refreshToken RefreshToken
	if err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("refresh token not found or expired")
		}
		logger.Error("Failed to validate refresh token", err)
		return 0, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	return refreshToken.UserID, nil
}

func (r *RefreshTokenRepository) Delete(token string) error {
	if err := r.db.Where("token = ?", token).Delete(&RefreshToken{}).Error; err != nil {
		logger.Error("Failed to delete refresh token", err)
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	if err := r.db.Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}).Error; err != nil {
		logger.Error("Failed to delete expired refresh tokens", err)
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteUserTokens(userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		logger.Error("Failed to delete user's refresh tokens", err)
		return fmt.Errorf("failed to delete user's refresh tokens: %w", err)
	}

	return nil
}
