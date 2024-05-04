package gateway

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	t.Run("Should send the request to the gateway", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		request := CreatePaymentGatewayDTO{
			PaymentID: uuid.NewString(),
			Amount:    100,
		}

		service := NewService()

		// Act
		err := service.Handle(ctx, request)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		request := CreatePaymentGatewayDTO{
			PaymentID: "invalid",
			Amount:    100,
		}

		service := NewService()

		// Act
		err := service.Handle(ctx, request)

		// Assert
		assert.NotNil(t, err)
	})
}
