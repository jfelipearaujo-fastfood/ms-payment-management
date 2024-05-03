package get_by_order_id

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil if the request is valid", func(t *testing.T) {
		// Arrange
		dto := GetByOrderIdDTO{
			OrderId: uuid.NewString(),
		}

		// Act
		err := dto.Validate()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Should return an error if the request is invalid", func(t *testing.T) {
		// Arrange
		dto := GetByOrderIdDTO{
			OrderId: "",
		}

		// Act
		err := dto.Validate()

		// Assert
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
