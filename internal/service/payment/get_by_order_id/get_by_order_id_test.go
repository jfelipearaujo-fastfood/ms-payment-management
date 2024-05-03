package get_by_order_id

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return the payments", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		orderId := uuid.NewString()

		repository.On("GetByOrderID", ctx, mock.Anything).
			Return([]payment_entity.Payment{
				{
					OrderId: orderId,
				},
			}, nil).
			Once()

		service := NewService(repository)

		req := GetByOrderIdDTO{
			OrderId: orderId,
		}

		// Act
		payments, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, req.OrderId, payments[0].OrderId)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		service := NewService(repository)

		req := GetByOrderIdDTO{
			OrderId: "",
		}

		// Act
		payments, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payments)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an empty result if no payments are found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		orderId := uuid.NewString()

		repository.On("GetByOrderID", ctx, mock.Anything).
			Return([]payment_entity.Payment{}, nil).
			Once()

		service := NewService(repository)

		req := GetByOrderIdDTO{
			OrderId: orderId,
		}

		// Act
		payments, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, payments)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error if repository return an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		orderId := uuid.NewString()

		repository.On("GetByOrderID", ctx, mock.Anything).
			Return([]payment_entity.Payment{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetByOrderIdDTO{
			OrderId: orderId,
		}

		// Act
		payments, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payments)
		repository.AssertExpectations(t)
	})
}
