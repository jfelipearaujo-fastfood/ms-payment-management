package update

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	t.Run("Should return nil if the request is valid", func(t *testing.T) {
		// Arrange
		dto := UpdatePaymentDTO{
			PaymentId: uuid.NewString(),
			Approved:  true,
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Nil(t, err)
	})

	t.Run("Should return error if the request is not valid", func(t *testing.T) {
		// Arrange
		dto := UpdatePaymentDTO{
			PaymentId: "",
			Approved:  true,
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NotNil(t, err)
	})
}
