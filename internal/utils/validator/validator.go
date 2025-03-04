package validator

import (
	"errors"
	"fmt"
	"go-manage-mysql/internal/models"

	"github.com/gustyaguero21/go-core/pkg/validator"
)

func ValidateData(user models.User, requiredFields []string) error {
	userMap := map[string]string{
		"name":     user.Name,
		"surname":  user.Surname,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"password": user.Password,
	}

	for _, field := range requiredFields {
		if value, exists := userMap[field]; exists && value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}

	if user.Email != "" && !validator.ValidateEmail(user.Email) {
		return errors.New("invalid email address")
	}

	if user.Password != "" && !validator.ValidatePassword(user.Password) {
		return errors.New("invalid password format")
	}

	return nil
}
