package get_by_order_id

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/get_by_order_id"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return the payments", func(t *testing.T) {
		// Arrange
		getByOrderIdService := mocks.NewMockGetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO](t)

		getByOrderIdService.On("Handle", mock.Anything, mock.Anything).
			Return([]payment_entity.Payment{
				{
					OrderId: uuid.NewString(),
				},
			}, nil).
			Once()

		reqBody := get_by_order_id.GetByOrderIdDTO{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		service := NewHandler(getByOrderIdService)

		// Act
		err = service.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		getByOrderIdService.AssertExpectations(t)
	})

	t.Run("Should return no content if the payments are empty", func(t *testing.T) {
		// Arrange
		getByOrderIdService := mocks.NewMockGetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO](t)

		getByOrderIdService.On("Handle", mock.Anything, mock.Anything).
			Return([]payment_entity.Payment{}, nil).
			Once()

		reqBody := get_by_order_id.GetByOrderIdDTO{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		service := NewHandler(getByOrderIdService)

		// Act
		err = service.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.Code)
		getByOrderIdService.AssertExpectations(t)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		getByOrderIdService := mocks.NewMockGetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO](t)

		getByOrderIdService.On("Handle", mock.Anything, mock.Anything).
			Return([]payment_entity.Payment{}, custom_error.ErrRequestNotValid).
			Once()

		reqBody := get_by_order_id.GetByOrderIdDTO{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		service := NewHandler(getByOrderIdService)

		// Act
		err = service.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusUnprocessableEntity, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusUnprocessableEntity,
			Message: "validation error",
			Details: "request not valid, please check the fields",
		}, he.Message)

		getByOrderIdService.AssertExpectations(t)
	})

	t.Run("Should return internal server error when an unexpected error occurs", func(t *testing.T) {
		// Arrange
		getByOrderIdService := mocks.NewMockGetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO](t)

		getByOrderIdService.On("Handle", mock.Anything, mock.Anything).
			Return([]payment_entity.Payment{}, assert.AnError).
			Once()

		reqBody := get_by_order_id.GetByOrderIdDTO{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		service := NewHandler(getByOrderIdService)

		// Act
		err = service.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusInternalServerError, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Details: "assert.AnError general error for testing",
		}, he.Message)

		getByOrderIdService.AssertExpectations(t)
	})
}
