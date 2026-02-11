package generator

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepGenerator(t *testing.T) {
	tests := []struct {
		name     string
		valueStr string
		count    int
		// check sequence
		expected []float64
		wantErr  bool
	}{
		{
			name:     "default step (1)",
			valueStr: "10", // initial=10, step=1 (default)
			count:    5,
			expected: []float64{10, 11, 12, 13, 14},
		},
		{
			name:     "custom step",
			valueStr: "0,5", // initial=0, step=5
			count:    5,
			expected: []float64{0, 5, 10, 15, 20},
		},
		{
			name:     "negative step",
			valueStr: "10,-1", // initial=10, step=-1
			count:    5,
			expected: []float64{10, 9, 8, 7, 6},
		},
		{
			name:     "empty value",
			valueStr: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			// Test Float64
			t.Run("float64", func(t *testing.T) {
				gen, err := newStepGenerator[float64](ctx, tt.valueStr, tt.count)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)

				var actual []float64
				for v := range gen {
					actual = append(actual, v)
				}
				assert.Equal(t, tt.expected, actual)
			})

			// Test Int64
			t.Run("int64", func(t *testing.T) {
				gen, err := newStepGenerator[int64](ctx, tt.valueStr, tt.count)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)

				var expectedInt []int64
				for _, v := range tt.expected {
					expectedInt = append(expectedInt, int64(v))
				}

				var actual []int64
				for v := range gen {
					actual = append(actual, v)
				}
				assert.Equal(t, expectedInt, actual)
			})
		})
	}
}
