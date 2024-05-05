package payment_entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil if request is valid", func(t *testing.T) {
		// Arrange
		item := NewPaymentItem(uuid.NewString(), "item1", 1)

		// Act
		err := item.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error if request is invalid", func(t *testing.T) {
		// Arrange
		item := NewPaymentItem("", "item1", 1)

		// Act
		err := item.Validate()

		// Assert
		assert.Error(t, err)
	})
}
