package get_by_id

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type GetByIdDTO struct {
	PaymentId string `param:"payment_id" validate:"required,uuid4"`
}

func (d *GetByIdDTO) Validate() error {
	validate := validator.New()

	if err := validate.Struct(d); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
