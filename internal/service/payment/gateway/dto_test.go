package gateway

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil if the request is valid", func(t *testing.T) {
		// Arrange
		request := CreatePaymentGatewayDTO{
			PaymentID: uuid.NewString(),
			Amount:    100,
		}

		// Act
		err := request.Validate()

		// Assert
		assert.Nil(t, err)
	})

	t.Run("Should return an error if the payment id is not a valid uuid", func(t *testing.T) {
		// Arrange
		request := CreatePaymentGatewayDTO{
			PaymentID: "invalid",
			Amount:    100,
		}

		// Act
		err := request.Validate()

		// Assert
		assert.NotNil(t, err)
	})

	t.Run("Should return an error if the amount is invalid", func(t *testing.T) {
		// Arrange
		request := CreatePaymentGatewayDTO{
			PaymentID: uuid.NewString(),
			Amount:    0,
		}

		// Act
		err := request.Validate()

		// Assert
		assert.NotNil(t, err)
	})
}
