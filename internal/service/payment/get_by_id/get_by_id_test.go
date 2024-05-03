package get_by_id

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHanle(t *testing.T) {
	t.Run("Should return the payment", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		paymentId := uuid.NewString()

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{
				PaymentId: paymentId,
			}, nil).
			Once()

		service := NewService(repository)

		req := GetByIdDTO{
			PaymentId: paymentId,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, req.PaymentId, payment.PaymentId)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		service := NewService(repository)

		req := GetByIdDTO{
			PaymentId: "",
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, payment_entity.Payment{}, payment)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error if the payment is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockPaymentRepository(t)

		paymentId := uuid.NewString()

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetByIdDTO{
			PaymentId: paymentId,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, payment_entity.Payment{}, payment)
		repository.AssertExpectations(t)
	})
}
