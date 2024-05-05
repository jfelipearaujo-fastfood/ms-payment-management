package cloud

import "github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"

type OrderProductionTopicItemContract struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type OrderProductionTopicContract struct {
	OrderId string                             `json:"order_id"`
	Items   []OrderProductionTopicItemContract `json:"items"`
}

func NewOrderProductionContractFromPayment(payment *payment_entity.Payment) *OrderProductionTopicContract {
	items := make([]OrderProductionTopicItemContract, len(payment.Items))

	for i, item := range payment.Items {
		items[i] = OrderProductionTopicItemContract{
			Id:       item.Id,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}

	return &OrderProductionTopicContract{
		OrderId: payment.OrderId,
		Items:   items,
	}
}
