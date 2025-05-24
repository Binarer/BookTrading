package token

import "time"

type Repository interface {
	Save(userID uint, token string, expiresAt time.Time) error
	Validate(token string) (uint, error)
	Delete(token string) error
	DeleteExpired() error
	DeleteUserTokens(userID uint) error
}
