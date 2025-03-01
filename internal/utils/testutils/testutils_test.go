package testutils

import (
	"go-manage-mysql/internal/models"
	"os"
	"testing"
)

func TestOpenMock(t *testing.T) {
	t.Run("Archivo JSON válido", func(t *testing.T) {
		filePath := "test_user.json"
		jsonData := `{"id":"1","name":"John","surname":"Doe","username":"johndoe","phone":"123456789","email":"johndoe","password":"Password123456"}`
		if err := os.WriteFile(filePath, []byte(jsonData), 0644); err != nil {
			t.Fatalf("Error creando archivo de prueba: %v", err)
		}
		defer os.Remove(filePath)

		got := OpenMock(filePath)
		want := models.User{
			ID:       "1",
			Name:     "John",
			Surname:  "Doe",
			Username: "johndoe",
			Phone:    "123456789",
			Email:    "johndoe",
			Password: "Password123456",
		}

		if got != want {
			t.Errorf("OpenMock() = %+v, want %+v", got, want)
		}
	})

	t.Run("Archivo no existente", func(t *testing.T) {
		got := OpenMock("non_existent.json")
		want := models.User{}

		if got != want {
			t.Errorf("OpenMock() con archivo inexistente = %+v, want %+v", got, want)
		}
	})

	t.Run("Archivo JSON inválido", func(t *testing.T) {
		filePath := "invalid.json"
		invalidData := `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe",`
		if err := os.WriteFile(filePath, []byte(invalidData), 0644); err != nil {
			t.Fatalf("Error creando archivo de prueba: %v", err)
		}
		defer os.Remove(filePath)

		got := OpenMock(filePath)
		want := models.User{}

		if got != want {
			t.Errorf("OpenMock() con JSON inválido = %+v, want %+v", got, want)
		}
	})
}
