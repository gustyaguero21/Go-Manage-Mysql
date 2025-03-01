package apperror

import (
	"errors"
	"testing"
)

func TestAppError(t *testing.T) {
	originalErr := errors.New("original error")

	wrappedErr := AppError("an error occurred", originalErr)

	if wrappedErr.Error() != "an error occurred. Error: original error" {
		t.Errorf("Expected error message to be 'an error occurred. Error: original error', but got %s", wrappedErr.Error())
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("Expected wrapped error to be unwrapped to original error")
	}
}
