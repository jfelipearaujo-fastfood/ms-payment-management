package custom_error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpAppError(t *testing.T) {
	t.Run("Should return an HTTP error", func(t *testing.T) {
		// Arrange

		// Act
		err := NewHttpAppError(123, "error", errors.New("my error"))

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, 123, err.Code)
		assert.Equal(t, AppError{
			Code:    123,
			Message: "error",
			Details: "my error",
		}, err.Message)
	})
}

func TestNewHttpAppErrorFromBusinessError(t *testing.T) {
	t.Run("Should return an HTTP error from business error", func(t *testing.T) {
		// Arrange
		buErr := New(123, "error", "error")

		// Act
		err := NewHttpAppErrorFromBusinessError(buErr)

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, 123, err.Code)
		assert.Equal(t, AppError{
			Code:    123,
			Message: "error",
			Details: "error",
		}, err.Message)
	})
}
