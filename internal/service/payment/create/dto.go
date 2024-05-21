package create

import (
	"context"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type CreatePaymentItemDTO struct {
	Id       string `json:"id" validate:"required,uuid4"`
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gte=1"`
}

type CreatePaymentDTO struct {
	OrderId   string `json:"order_id" validate:"required,uuid4"`
	PaymentId string `json:"payment_id" validate:"required,uuid4"`

	Items []CreatePaymentItemDTO `json:"items" validate:"required,dive"`

	TotalItems int     `json:"total_items" validate:"required,gte=1"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

func (dto *CreatePaymentDTO) Validate(ctx context.Context) error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		slog.ErrorContext(ctx, "error validating payment", "error", err)
		return custom_error.ErrRequestNotValid
	}

	return nil
}
