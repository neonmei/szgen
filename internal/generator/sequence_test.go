package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSequenceGenerator(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		gen, err := newSequenceGenerator[int64](context.Background(), "1,2,3", 10)
		require.NoError(t, err)

		var values []int64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 3)
		assert.Equal(t, []int64{1, 2, 3}, values)
	})

	t.Run("float64", func(t *testing.T) {
		gen, err := newSequenceGenerator[float64](context.Background(), "1.1,2.2,3.3", 10)
		require.NoError(t, err)

		var values []float64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 3)
		assert.Equal(t, []float64{1.1, 2.2, 3.3}, values)
	})

	t.Run("count less than sequence length", func(t *testing.T) {
		gen, err := newSequenceGenerator[int64](context.Background(), "1,2,3,4,5", 3)
		require.NoError(t, err)

		var values []int64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 3)
		assert.Equal(t, []int64{1, 2, 3}, values)
	})

	t.Run("single value", func(t *testing.T) {
		gen, err := newSequenceGenerator[int64](context.Background(), "1", 10)
		require.NoError(t, err)

		var values []int64
		for v := range gen {
			values = append(values, v)
		}

		assert.Len(t, values, 1)
		assert.Equal(t, []int64{1}, values)
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := newSequenceGenerator[int64](context.Background(), "1,invalid,3", 10)
		assert.Error(t, err)
	})

	t.Run("empty value", func(t *testing.T) {
		_, err := newSequenceGenerator[int64](context.Background(), "", 10)
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		gen, err := newSequenceGenerator[int64](ctx, "1,2,3,4,5", 10)
		require.NoError(t, err)

		// Cancel immediately
		cancel()

		var values []int64
		for v := range gen {
			values = append(values, v)
		}

		// Should be less than full sequence
		assert.Less(t, len(values), 5)
	})
}
