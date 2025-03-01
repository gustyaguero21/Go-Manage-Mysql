package apperror

import "fmt"

func AppError(msg string, err error) error {
	return fmt.Errorf("%s. Error: %w", msg, err)
}
