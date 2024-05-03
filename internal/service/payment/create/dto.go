package create

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type CreatePaymentDTO struct {
	OrderId   string `json:"order_id" validate:"required,uuid4"`
	PaymentId string `json:"payment_id" validate:"required,uuid4"`

	TotalItems int     `json:"total_items" validate:"required,gte=1"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

func (dto *CreatePaymentDTO) Validate() error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
