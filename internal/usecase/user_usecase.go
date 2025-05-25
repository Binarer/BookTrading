package usecase

import (
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/response"
	"booktrading/internal/domain/user"
	"errors"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserUseCase interface {
	Register(dto *user.CreateUserDTO) (*user.User, error)
	Login(dto *user.LoginDTO) (*response.TokenResponse, uint, error)
	Logout(refreshToken string) error
	GetByID(id uint) (*user.User, error)
	GetAll(page, pageSize int) ([]*user.User, int64, error)
	Update(id uint, dto *user.UpdateUserDTO) (*user.User, error)
	Delete(id uint) error
	GenerateTokenPair(user *user.User) (*response.TokenResponse, error)
	ValidateRefreshToken(token string) (*user.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
	jwtAuth  *jwtauth.JWTAuth
}

func NewUserUseCase(userRepo repository.UserRepository, jwtAuth *jwtauth.JWTAuth) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (u *userUseCase) Register(dto *user.CreateUserDTO) (*user.User, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := u.userRepo.GetByLogin(dto.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем нового пользователя
	newUser := dto.ToUser()
	newUser.Password = string(hashedPassword)

	// Если username не указан, используем login как username
	if newUser.Username == "" {
		newUser.Username = newUser.Login
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (u *userUseCase) Login(dto *user.LoginDTO) (*response.TokenResponse, uint, error) {
	// Get user by login
	existingUser, err := u.userRepo.GetByLogin(dto.Login)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return nil, 0, ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(dto.Password)); err != nil {
		return nil, 0, ErrInvalidCredentials
	}

	// Generate JWT token pair
	tokenPair, err := u.GenerateTokenPair(existingUser)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return tokenPair, existingUser.ID, nil
}

func (u *userUseCase) GetByID(id uint) (*user.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetAll(page, pageSize int) ([]*user.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.userRepo.GetAll(page, pageSize)
}

func (u *userUseCase) Update(id uint, dto *user.UpdateUserDTO) (*user.User, error) {
	existingUser, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Обновляем поля из DTO
	existingUser.UpdateFromDTO(dto)

	if err := u.userRepo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (u *userUseCase) Delete(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// Logout выполняет выход пользователя из системы
func (u *userUseCase) Logout(refreshToken string) error {
	// Валидируем токен перед инвалидацией
	token, err := u.jwtAuth.Decode(refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем тип токена
	tokenType, ok := token.PrivateClaims()["type"].(string)
	if !ok || tokenType != "refresh" {
		return errors.New("invalid token type")
	}

	return nil
}

// GenerateTokenPair генерирует пару токенов для пользователя
func (u *userUseCase) GenerateTokenPair(user *user.User) (*response.TokenResponse, error) {
	// Генерируем access token
	_, accessToken, err := u.jwtAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"login":   user.Login,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh token
	_, refreshToken, err := u.jwtAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"type":    "refresh",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenResponse := &response.TokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	return tokenResponse, nil
}

// ValidateRefreshToken проверяет refresh token
func (u *userUseCase) ValidateRefreshToken(token string) (*user.User, error) {
	// Валидируем токен
	jwtToken, err := u.jwtAuth.Decode(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	// Проверяем тип токена
	tokenType, ok := jwtToken.PrivateClaims()["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Получаем ID пользователя
	userID, ok := jwtToken.PrivateClaims()["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	// Получаем пользователя
	user, err := u.userRepo.GetByID(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
