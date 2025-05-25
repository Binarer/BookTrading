package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(u *user.User) error {
	// Проверяем уникальность логина
	var count int64
	if err := r.db.Model(u).Where("login = ?", u.Login).Count(&count).Error; err != nil {
		logger.Error("Failed to check login uniqueness", err)
		return err
	}
	if count > 0 {
		return errors.New("login must be unique")
	}

	// Если username не указан, используем login
	if u.Username == "" {
		u.Username = u.Login
	}

	if err := r.db.Create(u).Error; err != nil {
		logger.Error("Failed to create user", err)
		return err
	}
	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(id uint) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("User not found", fmt.Errorf("user with ID %d not found", id))
			return nil, repository.ErrNotFound
		}
		logger.Error("Failed to get user by ID", err)
		return nil, err
	}

	// Загружаем только ID книг пользователя
	var bookIDs []uint
	if err := r.db.Model(&book.Book{}).
		Where("user_id = ?", id).
		Pluck("id", &bookIDs).Error; err != nil {
		logger.Error("Failed to get user's book IDs", err)
		return nil, err
	}
	u.BookIDs = bookIDs

	return &u, nil
}

// GetByLogin получает пользователя по логину
func (r *UserRepository) GetByLogin(login string) (*user.User, error) {
	var u user.User
	result := r.db.Where("login = ?", login).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &u, nil
}

// GetAll получает всех пользователей с пагинацией
func (r *UserRepository) GetAll(page, pageSize int) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	// Получаем общее количество пользователей
	if err := r.db.Model(&user.User{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count users", err)
		return nil, 0, err
	}

	// Получаем пользователей с пагинацией
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		logger.Error("Failed to get users", err)
		return nil, 0, err
	}

	// Загружаем только ID книг для каждого пользователя
	for _, u := range users {
		var bookIDs []uint
		if err := r.db.Model(&book.Book{}).
			Where("user_id = ?", u.ID).
			Pluck("id", &bookIDs).Error; err != nil {
			logger.Error("Failed to get user's book IDs", err)
			return nil, 0, err
		}
		u.BookIDs = bookIDs
	}

	return users, total, nil
}

// Update обновляет пользователя в базе данных
func (r *UserRepository) Update(user *user.User) error {
	updates := map[string]interface{}{}

	if user.Username != "" {
		updates["username"] = user.Username
	}

	if user.Password != "" {
		updates["password"] = user.Password
	}

	if user.Avatar != "" {
		updates["avatar"] = user.Avatar
	}

	return r.db.Model(user).Updates(updates).Error
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(id uint) error {
	// Проверяем наличие книг у пользователя
	var count int64
	if err := r.db.Model(&book.Book{}).
		Where("user_id = ?", id).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check user's books", err)
		return err
	}

	if count > 0 {
		return errors.New("cannot delete user: they have books")
	}

	result := r.db.Delete(&user.User{}, id)
	if result.Error != nil {
		logger.Error("Failed to delete user", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}
