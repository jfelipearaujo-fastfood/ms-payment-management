package cloud

import "github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"

type UpdateOrderTopicPaymentContract struct {
	PaymentId string `json:"id"`
	State     string `json:"state"`
}

type UpdateOrderTopicContract struct {
	OrderId string                          `json:"order_id"`
	Payment UpdateOrderTopicPaymentContract `json:"payment"`
}

func NewUpdateOrderContractFromPayment(payment *payment_entity.Payment) *UpdateOrderTopicContract {
	payment.RefreshStateTitle()

	return &UpdateOrderTopicContract{
		OrderId: payment.OrderId,
		Payment: UpdateOrderTopicPaymentContract{
			PaymentId: payment.PaymentId,
			State:     payment.StateTitle,
		},
	}
}
