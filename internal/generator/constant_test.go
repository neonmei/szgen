package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConstantGenerator(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		gen, err := newConstantGenerator[int64](context.Background(), "42", 5)
		require.NoError(t, err)

		var values []int64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 5)
		for _, v := range values {
			assert.Equal(t, int64(42), v)
		}
	})

	t.Run("float64", func(t *testing.T) {
		gen, err := newConstantGenerator[float64](context.Background(), "3.14", 3)
		require.NoError(t, err)

		var values []float64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 3)
		for _, v := range values {
			assert.Equal(t, 3.14, v)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := newConstantGenerator[int64](context.Background(), "invalid", 5)
		assert.Error(t, err)
	})

	t.Run("empty value", func(t *testing.T) {
		_, err := newConstantGenerator[int64](context.Background(), "", 5)
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Create a large count
		gen, err := newConstantGenerator[int64](ctx, "1", 100000)
		require.NoError(t, err)

		// Consume one value then cancel
		count := 0
		cancel() // Cancel immediately for this test to ensure it stops

		for range gen {
			count++
		}

		// Should stop significantly before 100000
		assert.Less(t, count, 100000)
	})
}
