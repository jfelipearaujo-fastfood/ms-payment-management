package payment_entity

import "time"

type Payment struct {
	OrderId   string `json:"order_id"`
	PaymentId string `json:"payment_id"`

	Items []PaymentItem `json:"items"`

	TotalItems int          `json:"total_items"`
	Amount     float64      `json:"amount"`
	State      PaymentState `json:"state"`
	StateTitle string       `json:"state_title"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPayment(orderId string, paymentId string, items []PaymentItem, totalItems int, amount float64, now time.Time) Payment {
	return Payment{
		OrderId:   orderId,
		PaymentId: paymentId,

		Items: items,

		TotalItems: totalItems,
		Amount:     amount,
		State:      WaitingForApproval,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Payment) IsInState(states ...PaymentState) bool {
	for _, state := range states {
		if p.State == state {
			return true
		}
	}

	return false
}

func (p *Payment) RefreshStateTitle() {
	p.StateTitle = p.State.String()
}

func (p *Payment) UpdateState(newState PaymentState, now time.Time) {
	p.State = newState
	p.StateTitle = p.State.String()

	p.UpdatedAt = now
}

func (p *Payment) Exists() bool {
	return p.OrderId != "" && p.PaymentId != ""
}
