package custom_error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBusinessErr(t *testing.T) {
	t.Run("Should return true when error is a business error", func(t *testing.T) {
		// Arrange
		err := New(123, "error", "error")

		// Act
		result := IsBusinessErr(err)

		// Assert
		assert.True(t, result)
	})

	t.Run("Should return false when error is not a business error", func(t *testing.T) {
		// Arrange
		err := errors.New("error")

		// Act
		result := IsBusinessErr(err)

		// Assert
		assert.False(t, result)
	})

	t.Run("Should return false when error is nil", func(t *testing.T) {
		// Arrange

		// Act
		result := IsBusinessErr(nil)

		// Assert
		assert.False(t, result)
	})
}
