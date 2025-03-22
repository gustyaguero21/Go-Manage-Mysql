package mocks

const (
	CreateUser = `{
                "id":"1",
                "name": "John",
                "surname": "Doe",
                "username": "johndoe",
                "phone":"23456789",
                "email": "johndoe@example.com",
                "password": "Password1234"
            }`
	UpdateUser = `{
                "name": "Johncito",
                "surname": "Doecito",
                "phone":"23456789",
                "email": "johncitodoecito@example.com"
            }`

	InvalidJSON   = `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe","phone":"123456789", "email": "johndoe@example.com", "password": }`
	ValidateError = `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe","phone":"123456789"}`

	InvalidQueryParam = `{
                "surname": "Doecito",
                "phone":"23456789",
                "email": "johncitodoecito@example.com"
            }`

	ChangePwd = `{
                "username": "johndoe",
                "password": "Password1234"
            }`
	LoginUser = `{
                "username": "johndoe",
                "password":"Password1234"
            }`

	InvalidBody = `{
                "username": "johndoe",
                "password": "Password1234
            }`
	InvalidPasswordFormat = `{
                "username": "johndoe",
                "password": "Pass"
            }`
	InvalidRequest = `{
                "username": "johndoe",
                "password":"Password1234",
            }`
	UserNotFound = `{
                "username": "nonexistent",
                "password":"Password1234"
            }`
	Unauthorized = `{
                "username": "johndoe",
                "password":"Password1234"
            }`
)
