package create

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil if request is valid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		dto := CreatePaymentDTO{
			OrderId:   uuid.NewString(),
			PaymentId: uuid.NewString(),
			Items: []CreatePaymentItemDTO{
				{
					Id:       uuid.NewString(),
					Name:     "item1",
					Quantity: 1,
				},
				{
					Id:       uuid.NewString(),
					Name:     "item2",
					Quantity: 1,
				},
			},
			TotalItems: 1,
			Amount:     100,
		}

		// Act
		err := dto.Validate(ctx)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error if request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		dto := CreatePaymentDTO{
			OrderId:    uuid.NewString(),
			PaymentId:  uuid.NewString(),
			TotalItems: 0,
			Amount:     100,
		}

		// Act
		err := dto.Validate(ctx)

		// Assert
		assert.Error(t, err)
	})
}
