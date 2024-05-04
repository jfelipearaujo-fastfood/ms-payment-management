package gateway

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type CreatePaymentGatewayDTO struct {
	PaymentID string  `json:"payment_id" validate:"required,uuid4"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

func (d *CreatePaymentGatewayDTO) Validate() error {
	validator := validator.New()

	if err := validator.Struct(d); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
