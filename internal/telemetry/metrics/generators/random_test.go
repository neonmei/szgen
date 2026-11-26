package generators

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
			name:     "single value is the range from zero",
			valueStr: "10", // bounds 10 and 0 (default) -> [0, 10]
			count:    50,
			min:      0,
			max:      10,
		},
		{
			name:     "single negative value",
			valueStr: "-10", // bounds -10 and 0 (default) -> [-10, 0]
			count:    50,
			min:      -10,
			max:      0,
		},
		{
			name:     "reversed bounds",
			valueStr: "100,50", // max,min -> [50, 100]
			count:    100,
			min:      50,
			max:      100,
		},
		{
			name:     "equal bounds error",
			valueStr: "5,5",
			count:    10,
			wantErr:  true,
		},
		{
			name:     "zero error",
			valueStr: "0",
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
