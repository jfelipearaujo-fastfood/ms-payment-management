package get_by_order_id

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type GetByOrderIdDTO struct {
	OrderId string `param:"order_id" validate:"required,uuid4"`
}

func (d *GetByOrderIdDTO) Validate() error {
	validator := validator.New()

	if err := validator.Struct(d); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
