package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyData struct {
	Val int
}

func TestState(t *testing.T) {
	t.Run("Should be able to retrieve data", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		var key CtxKeyType = "test"

		myData := MyData{
			Val: 1,
		}

		state := NewState[MyData](key)

		// Act
		ctx = state.enrich(ctx, &myData)

		// Assert
		res := state.retrieve(ctx)
		assert.Equal(t, myData.Val, res.Val)
	})

	t.Run("Should return nil if data is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		var key CtxKeyType = "test"

		state := NewState[MyData](key)

		// Act
		res := state.retrieve(ctx)

		// Assert
		assert.Empty(t, res)
	})
}
