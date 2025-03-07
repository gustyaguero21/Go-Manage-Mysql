package config

// router params
const (
	Port    = ":8080"
	BaseURL = "/api/go-manage"
)

//fields validating

var (
	Create_ValidateFields = []string{"name", "surname", "username", "phone", "email", "password"}
	Update_ValidateFields = []string{"name", "surname", "phone", "email"}
)
