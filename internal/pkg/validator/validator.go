package validator

import (
	"booktrading/internal/domain/tag"
	"encoding/base64"
	"fmt"

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
		_, err := base64.StdEncoding.DecodeString(value)
		return err == nil
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