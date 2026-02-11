package generator

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomGenerator(t *testing.T) {
	tests := []struct {
		name     string
		valueStr string
		count    int
		min      float64
		max      float64
		wantErr  bool
	}{
		{
			name:     "valid range",
			valueStr: "50,100", // min=50, max=100
			count:    100,
			min:      50,
			max:      100,
		},
		{
			name:     "single value (min only)",
			valueStr: "10", // min=10, max=0 (default) -> error: min >= max
			count:    50,
			wantErr:  true,
		},
		{
			name:     "min >= max error",
			valueStr: "100,50", // min=100, max=50 -> error
			count:    10,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			t.Run("float64", func(t *testing.T) {
				genHelper[float64](t, ctx, tt.valueStr, tt.count, tt.min, tt.max, tt.wantErr)
			})

			t.Run("int64", func(t *testing.T) {
				genHelper[int64](t, ctx, tt.valueStr, tt.count, tt.min, tt.max, tt.wantErr)
			})
		})
	}
}

func genHelper[T int64 | float64](t *testing.T, ctx context.Context, valueStr string, count int, min, max float64, wantErr bool) {
	gen, err := newRandomGenerator[T](ctx, valueStr, count)
	if wantErr {
		assert.Error(t, err)
		return
	}
	require.NoError(t, err)

	for v := range gen {
		val := float64(v)
		assert.GreaterOrEqual(t, val, min, "value %v < min %v", val, min)
		assert.LessOrEqual(t, val, max, "value %v > max %v", val, max)
	}
}
