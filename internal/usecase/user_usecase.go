package usecase

import (
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/repository"
	"errors"
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
	GetAll() ([]*user.User, error)
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
	// Check if user already exists
	existingUser, err := u.userRepo.GetByLogin(dto.Login)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	newUser := &user.User{
		Username:     dto.Username,
		Login:        dto.Login,
		PasswordHash: string(hashedPassword),
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *userUsecase) Login(dto *user.LoginDTO) (*user.TokenResponse, error) {
	// Get user by login
	existingUser, err := u.userRepo.GetByLogin(dto.Login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := u.jwtSvc.GenerateToken(existingUser)
	if err != nil {
		return nil, err
	}

	return &user.TokenResponse{Token: token}, nil
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

func (u *userUsecase) GetAll() ([]*user.User, error) {
	return u.userRepo.GetAll()
}

func (u *userUsecase) Update(id uint, dto *user.UpdateUserDTO) (*user.User, error) {
	existingUser, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if dto.Username != "" {
		existingUser.Username = dto.Username
	}
	if dto.Description != nil {
		existingUser.Description = dto.Description
	}
	if dto.Avatar != nil {
		existingUser.Avatar = dto.Avatar
	}

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
