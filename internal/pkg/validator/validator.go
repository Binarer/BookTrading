package validator

import (
	"booktrading/internal/domain/tag"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate обертка для go-playground/validator.Validate
type Validate struct {
	*validator.Validate
}

// New создает новый экземпляр валидатора
func New() *Validate {
	v := validator.New()

	// Регистрация пользовательской валидации для base64
	_ = v.RegisterValidation("base64", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true // Пустая строка считается валидной
		}

		// Проверяем, является ли строка data URL
		if !strings.HasPrefix(value, "data:image/") {
			return false
		}

		// Проверяем формат изображения
		contentType := strings.TrimPrefix(value, "data:")
		if !strings.HasPrefix(contentType, "image/jpeg;base64,") && !strings.HasPrefix(contentType, "image/png;base64,") {
			return false
		}

		// Извлекаем base64 часть
		parts := strings.Split(value, ",")
		if len(parts) != 2 {
			return false
		}

		// Проверяем размер (5MB = 5 * 1024 * 1024 байт)
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return false
		}

		if len(decoded) > 5*1024*1024 {
			return false
		}

		return true
	})

	// Регистрация пользовательской валидации для состояния
	_ = v.RegisterValidation("state", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		validStates := map[string]bool{
			"available": true,
			"trading":   true,
			"traded":    true,
		}
		return validStates[value]
	})

	return &Validate{v}
}

// ValidateStruct валидирует структуру с помощью валидатора
func (v *Validate) ValidateStruct(s interface{}) error {
	return v.Struct(s)
}

// ValidateTagName проверяет, является ли имя тега уникальным
func ValidateTagName(tagRepo interface {
	GetByName(name string) (*tag.Tag, error)
}, name string, excludeID uint) error {
	existingTag, err := tagRepo.GetByName(name)
	if err != nil {
		return err
	}

	if existingTag != nil && existingTag.ID != excludeID {
		return fmt.Errorf("tag name '%s' is not unique", name)
	}

	return nil
}
