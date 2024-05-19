package get_by_order_id

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/service"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/get_by_order_id"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	getByOrderId service.GetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO]
}

func NewHandler(getByOrderId service.GetPaymentsByOrderIDService[get_by_order_id.GetByOrderIdDTO]) *Handler {
	return &Handler{
		getByOrderId: getByOrderId,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request get_by_order_id.GetByOrderIdDTO

	if err := ctx.Bind(&request); err != nil {
		return err
	}

	context := ctx.Request().Context()

	payments, err := h.getByOrderId.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	if len(payments) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	return ctx.JSON(200, payments)
}
