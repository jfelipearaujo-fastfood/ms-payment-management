package payment_entity

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type PaymentItem struct {
	Id       string `json:"id" validate:"required,uuid4"`
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gte=1"`
}

func NewPaymentItem(id, name string, quantity int) PaymentItem {
	return PaymentItem{
		Id:       id,
		Name:     name,
		Quantity: quantity,
	}
}

func (p *PaymentItem) Validate() error {
	validator := validator.New()

	if err := validator.Struct(p); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
