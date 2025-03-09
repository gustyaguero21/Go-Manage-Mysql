package config

// messages from responses
const (
	//success messages
	CreatedUserMessage = "user created successfully"
	SearchUserMessage  = "user found successfully"
	UpdateUserMessage  = "user updated successfully"
	DeleteUserMessage  = "user deleted successfully"
	ChangePwdMessage   = "password changed successfully"

	//error messages

	ErrCreatingUser  = "error creating user"
	ErrSearchingUser = "error searching user"
	ErrUpdatingUser  = "error updating user data"
	ErrDeletingUser  = "error deleting user data"
	ErrChangingPwd   = "error changing user password"
	ErrLoginUser     = "error login user"
)
