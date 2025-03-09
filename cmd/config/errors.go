package config

import "errors"

const (
	//service errors
	ErrCreatingUser  = "error creating user"
	ErrSearchingUser = "error searching user"
	ErrUpdatingUser  = "error updating user data"
	ErrDeletingUser  = "error deleting user data"
	ErrChangingPwd   = "error changing user password"
	ErrLoginUser     = "error login user"

	//errors
	ErrInvalidQueryParam    = "invalid query param"
	ErrInvalidBody          = "invalid body request"
	ErrAllFieldsAreRequired = "all fields are required"
	ErrUnauthorizedUser     = "invalid credentials. Please check username & password"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDbError           = errors.New("db error")
	ErrUserNotFound      = errors.New("user not found")
	ErrNoNewData         = errors.New("no new data to update")
	ErrPwdMatching       = errors.New("passwords doesnt match")
	ErrRecordNotFound    = errors.New("record not found")
)
