package testutils

import (
	"encoding/json"
	"fmt"
	"go-manage-mysql/internal/models"
	"os"
)

func OpenMock(path string) models.User {
	var user models.User

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening json file")
		return models.User{}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&user)
	if err != nil {
		fmt.Println("Error unmarshalling json file")
		return models.User{}
	}

	return user
}
