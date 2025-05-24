package usecase

import (
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserUsecase interface {
	Register(dto *user.CreateUserDTO) (*user.User, error)
	Login(dto *user.LoginDTO) (*user.TokenResponse, error)
	GetByID(id uint) (*user.User, error)
	GetAll(page, pageSize int) ([]*user.User, int64, error)
	Update(id uint, dto *user.UpdateUserDTO) (*user.User, error)
	Delete(id uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwtSvc   *jwt.Service
}

func NewUserUsecase(userRepo repository.UserRepository, jwtSvc *jwt.Service) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

func (u *userUsecase) Register(dto *user.CreateUserDTO) (*user.User, error) {
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

func (u *userUsecase) Login(dto *user.LoginDTO) (*user.TokenResponse, error) {
	// Get user by login
	existingUser, err := u.userRepo.GetByLogin(dto.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token pair
	tokenPair, err := u.jwtSvc.GenerateTokenPair(existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return &user.TokenResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (u *userUsecase) GetByID(id uint) (*user.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetAll(page, pageSize int) ([]*user.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.userRepo.GetAll(page, pageSize)
}

func (u *userUsecase) Update(id uint, dto *user.UpdateUserDTO) (*user.User, error) {
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

func (u *userUsecase) Delete(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
