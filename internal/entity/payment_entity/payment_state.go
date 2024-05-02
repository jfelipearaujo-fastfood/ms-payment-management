package payment_entity

type PaymentState int

const (
	None               PaymentState = iota
	WaitingForApproval              // When the payment request is sent to the payment gateway
	Approved                        // When the payment is approved by the payment gateway
	Rejected                        // When the payment is rejected by the payment gateway
)

var (
	payment_state_machine = map[PaymentState][]PaymentState{
		None:               {WaitingForApproval},
		WaitingForApproval: {Approved, Rejected},
		Approved:           {},
		Rejected:           {},
	}
)

func NewPaymentState(title string) PaymentState {
	state, ok := map[string]PaymentState{
		"WaitingForApproval": WaitingForApproval,
		"Approved":           Approved,
		"Rejected":           Rejected,
	}[title]
	if !ok {
		return None
	}

	return state
}

func (s PaymentState) CanTransitionTo(to PaymentState) bool {
	for _, allowed := range payment_state_machine[s] {
		if to == allowed {
			return true
		}
	}
	return false
}

func (s PaymentState) String() string {
	text, ok := map[PaymentState]string{
		None:               "None",
		WaitingForApproval: "WaitingForApproval",
		Approved:           "Approved",
		Rejected:           "Rejected",
	}[s]
	if !ok {
		return "Unknown"
	}

	return text
}
