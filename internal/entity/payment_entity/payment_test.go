package payment_entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsInState(t *testing.T) {
	t.Run("Should return true if the payment is in one of the states", func(t *testing.T) {
		// Arrange
		now := time.Now()

		states := []PaymentState{WaitingForApproval, Approved, Rejected}

		for _, state := range states {
			payment := NewPayment("order_id", "payment_id", 1, 1.23, now)
			payment.State = state

			// Act
			res := payment.IsInState(state)

			// Assert
			assert.True(t, res)
		}
	})

	t.Run("Should return false if the payment is not in one of the states", func(t *testing.T) {
		// Arrange
		now := time.Now()

		states := []PaymentState{WaitingForApproval, Approved, Rejected}

		for _, state := range states {
			payment := NewPayment("order_id", "payment_id", 1, 1.23, now)
			payment.State = state

			// Act
			res := payment.IsInState(None)

			// Assert
			assert.False(t, res)
		}
	})
}

func TestRefreshStateTitle(t *testing.T) {
	t.Run("Should refresh the state title", func(t *testing.T) {
		// Arrange
		now := time.Now()

		payment := NewPayment("order_id", "payment_id", 1, 1.23, now)
		payment.State = Approved

		// Act
		payment.RefreshStateTitle()

		// Assert
		assert.Equal(t, "Approved", payment.StateTitle)
	})
}

func TestUpdateState(t *testing.T) {
	t.Run("Should update the state", func(t *testing.T) {
		// Arrange
		now := time.Now()

		payment := NewPayment("order_id", "payment_id", 1, 1.23, now)

		// Act
		payment.UpdateState(Approved, now)

		// Assert
		assert.Equal(t, Approved, payment.State)
		assert.Equal(t, "Approved", payment.StateTitle)
		assert.Equal(t, now, payment.UpdatedAt)
	})
}
