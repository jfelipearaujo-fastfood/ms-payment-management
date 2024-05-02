package payment_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaymentState(t *testing.T) {
	t.Run("Should return the correct state", func(t *testing.T) {
		// Arrange
		cases := []struct {
			title    string
			expected PaymentState
		}{
			{"WaitingForApproval", WaitingForApproval},
			{"Approved", Approved},
			{"Rejected", Rejected},
		}

		for _, c := range cases {
			// Act
			res := NewPaymentState(c.title)

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})

	t.Run("Should return None when state is invalid", func(t *testing.T) {
		// Arrange
		title := "Invalid"

		// Act
		res := NewPaymentState(title)

		// Assert
		assert.Equal(t, None, res)
	})
}

func TestString(t *testing.T) {
	t.Run("Should return the string representation of the state", func(t *testing.T) {
		// Arrange
		cases := []struct {
			state    PaymentState
			expected string
		}{
			{None, "None"},
			{WaitingForApproval, "WaitingForApproval"},
			{Approved, "Approved"},
			{Rejected, "Rejected"},
			{PaymentState(100), "Unknown"},
		}

		for _, c := range cases {
			// Act
			res := c.state.String()

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})
}

func TestCanTransitionTo(t *testing.T) {
	t.Run("Should return true when the transition is valid", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from     PaymentState
			to       PaymentState
			expected bool
		}{
			{None, WaitingForApproval, true},
			{WaitingForApproval, Approved, true},
			{WaitingForApproval, Rejected, true},
			{Approved, Approved, false},
			{Approved, Rejected, false},
			{Rejected, Approved, false},
			{Rejected, Rejected, false},
		}

		for _, c := range cases {
			// Act
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})

	t.Run("Should return false when the transition is invalid", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from     PaymentState
			to       PaymentState
			expected bool
		}{
			{None, Approved, false},
			{None, Rejected, false},
			{WaitingForApproval, None, false},
			{Approved, None, false},
			{Rejected, None, false},
		}

		for _, c := range cases {
			// Act
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})
}
