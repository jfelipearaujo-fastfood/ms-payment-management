package create

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/provider/mocks"
	repository_mocks "github.com/jfelipearaujo-org/ms-payment-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestH(t *testing.T) {
	t.Run("Should create a payment", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := mocks.NewMockTimeProvider(t)

		repository.On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		req := CreatePaymentDTO{
			OrderId:    uuid.NewString(),
			PaymentId:  uuid.NewString(),
			TotalItems: 1,
			Amount:     100,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error if request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := mocks.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		req := CreatePaymentDTO{
			OrderId:    uuid.NewString(),
			PaymentId:  uuid.NewString(),
			TotalItems: -1,
			Amount:     100,
		}

		// Act
		payment, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error if repository returns error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockPaymentRepository(t)
		timeProvider := mocks.NewMockTimeProvider(t)

		repository.On("Create", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		req := CreatePaymentDTO{
			OrderId:    uuid.NewString(),
			PaymentId:  uuid.NewString(),
			TotalItems: 1,
			Amount:     100,
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
