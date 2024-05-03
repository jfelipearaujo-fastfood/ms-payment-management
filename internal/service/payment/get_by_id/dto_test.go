package get_by_id

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil if the request is valid", func(t *testing.T) {
		// Arrange
		dto := GetByIdDTO{
			PaymentId: uuid.NewString(),
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		dto := GetByIdDTO{
			PaymentId: "",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
	})
}
