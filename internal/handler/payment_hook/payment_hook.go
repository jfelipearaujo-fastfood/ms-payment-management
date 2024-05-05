package payment_hook

import (
	"log/slog"
	"net/http"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/update"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	updatePaymentService service.UpdatePaymentService[update.UpdatePaymentDTO]
	orderProductionTopic cloud.TopicService
	updateOrderTopic     cloud.TopicService
}

func NewHandler(
	updatePaymentService service.UpdatePaymentService[update.UpdatePaymentDTO],
	orderProductionTopic cloud.TopicService,
	updateOrderTopic cloud.TopicService,
) *Handler {
	return &Handler{
		updatePaymentService: updatePaymentService,
		orderProductionTopic: orderProductionTopic,
		updateOrderTopic:     updateOrderTopic,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request update.UpdatePaymentDTO

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	payment, err := h.updatePaymentService.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	if payment.State == payment_entity.Approved {
		slog.InfoContext(context, "payment approved, sending to production topic", "payment_id", payment.PaymentId)

		req := cloud.NewOrderProductionContractFromPayment(payment)

		messageId, err := h.orderProductionTopic.PublishMessage(context, req)
		if err != nil {
			slog.ErrorContext(context, "error publishing message to production topic", "error", err)
		}

		if messageId != nil {
			slog.InfoContext(context, "message published to production topic", "message_id", *messageId)
		}
	}

	slog.InfoContext(context, "payment updated, sending to update order topic", "payment_id", payment.PaymentId)

	req := cloud.NewUpdateOrderContractFromPayment(payment)

	messageId, err := h.updateOrderTopic.PublishMessage(context, req)
	if err != nil {
		slog.ErrorContext(context, "error publishing message to update order topic", "error", err)
	}

	if messageId != nil {
		slog.InfoContext(context, "message published to update order topic", "message_id", *messageId)
	}

	return ctx.JSON(http.StatusCreated, payment)
}
