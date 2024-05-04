package payment_hook

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/service"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/update"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.UpdatePaymentService[update.UpdatePaymentDTO]
}

func NewHandler(
	service service.UpdatePaymentService[update.UpdatePaymentDTO],
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request update.UpdatePaymentDTO

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	payment, err := h.service.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	return ctx.JSON(http.StatusCreated, payment)
}
