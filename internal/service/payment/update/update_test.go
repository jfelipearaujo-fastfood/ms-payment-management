package update

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	provider_mocks "github.com/jfelipearaujo-org/ms-payment-management/internal/provider/mocks"
	repository_mocks "github.com/jfelipearaujo-org/ms-payment-management/internal/repository/mocks"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should update the payment when was approved", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{
				State: payment_entity.WaitingForApproval,
			}, nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		repository.On("Update", ctx, mock.Anything).
			Return(nil).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdatePaymentDTO{
			PaymentId: uuid.NewString(),
			Approved:  true,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should update the payment when was rejected", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{
				State: payment_entity.WaitingForApproval,
			}, nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		repository.On("Update", ctx, mock.Anything).
			Return(nil).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdatePaymentDTO{
			PaymentId: uuid.NewString(),
			Approved:  false,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		req := UpdatePaymentDTO{
			PaymentId: "abc",
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error if the payment is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{}, custom_error.ErrPaymentNotFound).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdatePaymentDTO{
			PaymentId: uuid.NewString(),
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(payment_entity.Payment{}, nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		repository.On("Update", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdatePaymentDTO{
			PaymentId: uuid.NewString(),
			Approved:  true,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
