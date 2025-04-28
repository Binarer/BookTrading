package validator

import (
	"booktrading/internal/domain/tag"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom validation for base64
	_ = validate.RegisterValidation("base64", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		
		// Check if it's a data URL
		if strings.HasPrefix(value, "data:") {
			parts := strings.Split(value, ",")
			if len(parts) != 2 {
				return false
			}
			value = parts[1]
		}
		
		_, err := base64.StdEncoding.DecodeString(value)
		return err == nil
	})
}

// ValidateStruct validates a struct using the custom validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
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