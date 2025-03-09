package config

import "errors"

// service errors
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDbError           = errors.New("db error")
	ErrUserNotFound      = errors.New("user not found")
	ErrNoNewData         = errors.New("no new data to update")
	ErrPwdMatching       = errors.New("passwords doesnt match")
	ErrRecordNotFound    = errors.New("record not found")
)

// handler errors
var (
	ErrInvalidQueryParam    = "invalid query param"
	ErrInvalidBody          = "invalid body request"
	ErrAllFieldsAreRequired = "all fields are required"
	ErrUnauthorizedUser     = "invalid credentials. Please check username & password"
)
