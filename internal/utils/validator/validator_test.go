package validator

import (
	"go-manage-mysql/internal/models"
	"testing"
)

func TestValidateData(t *testing.T) {
	tests := []struct {
		name           string
		user           models.User
		requiredFields []string
		expectError    bool
	}{
		{
			name: "Valid create user",
			user: models.User{
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Phone:    "123456789",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			requiredFields: []string{"name", "surname", "username", "phone", "email", "password"},
			expectError:    false,
		},
		{
			name: "Missing required field",
			user: models.User{
				Name:    "John",
				Surname: "Doe",
				Email:   "johndoe@example.com",
				Phone:   "123456789",
			},
			requiredFields: []string{"name", "surname", "username", "phone", "email", "password"},
			expectError:    true,
		},
		{
			name: "Invalid email",
			user: models.User{
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Phone:    "123456789",
				Email:    "invalid-email", // Email inválido para prueba
				Password: "StrongP@ssw0rd2024!",
			},
			requiredFields: []string{"name", "surname", "username", "phone", "email", "password"},
			expectError:    true,
		},
		{
			name: "Invalid password",
			user: models.User{
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Phone:    "123456789",
				Email:    "johndoe@example.com",
				Password: "123", // Contraseña débil para la prueba
			},
			requiredFields: []string{"name", "surname", "username", "phone", "email", "password"},
			expectError:    true,
		},
		{
			name: "Valid update user (only name, surname, email, phone)",
			user: models.User{
				Name:    "John",
				Surname: "Doe",
				Phone:   "123456789",
				Email:   "johndoe@example.com",
			},
			requiredFields: []string{"name", "surname", "email", "phone"},
			expectError:    false,
		},
		{
			name: "Missing required field in update",
			user: models.User{
				Name:    "John",
				Surname: "Doe",
				Phone:   "123456789",
			},
			requiredFields: []string{"name", "surname", "email", "phone"},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateData(tt.user, tt.requiredFields)
			if (err != nil) != tt.expectError {
				t.Errorf("Test failed for case: %s. Expected error: %v, got: %v", tt.name, tt.expectError, err)
			}
		})
	}
}
