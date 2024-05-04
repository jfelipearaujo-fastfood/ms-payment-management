package payment_hook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/update"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should create a payment gateway", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdatePaymentService[update.UpdatePaymentDTO](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(&payment_entity.Payment{}, nil).
			Once()

		reqBody := update.UpdatePaymentDTO{
			Approved: true,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetParamNames("payment_id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdatePaymentService[update.UpdatePaymentDTO](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, custom_error.ErrRequestNotValid).
			Once()

		reqBody := update.UpdatePaymentDTO{
			Approved: true,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetParamNames("payment_id")
		ctx.SetParamValues("invalid-payment-id")

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

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

		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error when an unexpected error occurs", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdatePaymentService[update.UpdatePaymentDTO](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		reqBody := update.UpdatePaymentDTO{
			Approved: true,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetParamNames("payment_id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

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

		service.AssertExpectations(t)
	})
}
