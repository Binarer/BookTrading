package validator

import (
	"booktrading/internal/domain/tag"
	"encoding/base64"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate wraps the go-playground/validator.Validate
type Validate struct {
	*validator.Validate
}

// New creates a new validator instance
func New() *Validate {
	v := validator.New()
	
	// Register custom validation for base64
	_ = v.RegisterValidation("base64", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true // Empty string is considered valid
		}
		_, err := base64.StdEncoding.DecodeString(value)
		return err == nil
	})

	// Register custom validation for state
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

// ValidateStruct validates a struct using the validator
func (v *Validate) ValidateStruct(s interface{}) error {
	return v.Struct(s)
}

// ValidateTagName validates if a tag name is unique
func ValidateTagName(tagRepo interface {
	GetByName(name string) (*tag.Tag, error)
}, name string, excludeID int64) error {
	existingTag, err := tagRepo.GetByName(name)
	if err != nil {
		return err
	}
	
	if existingTag != nil && existingTag.ID != excludeID {
		return fmt.Errorf("tag name '%s' is not unique", name)
	}
	
	return nil
} 