package config

//service errors

const (
	ErrUserAlreadyExists    = "user already exists"
	ErrUserNotFound         = "user not found"
	ErrRecordNotFound       = "record not found"
	ErrCreatingUser         = "error creating user"
	ErrSearchingUser        = "error searching user"
	ErrUpdatingUser         = "error updating user data"
	ErrDeletingUser         = "error deleting user data"
	ErrChangingPwd          = "error changing user password"
	ErrPwdMatching          = "passwords doesnt match"
	ErrDbError              = "db error"
	ErrInvalidQueryParam    = "invalid query param"
	ErrAllFieldsAreRequired = "all fields are required"
	ErrUnauthorizedUser     = "invalid credentials. Please check username & password"
)
