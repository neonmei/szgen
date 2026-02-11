package generator

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSineGenerator_Float64(t *testing.T) {
	tests := []struct {
		name     string
		valueStr string
		count    int
		// function to verify the generated values
		verify  func(t *testing.T, values []float64)
		wantErr bool
	}{
		{
			name:     "default sine wave (amplitude=100 from valueStr?)",
			valueStr: "100", // amplitude=100. b=10 (default), vShift=1, hShift=0
			count:    5,
			verify: func(t *testing.T, values []float64) {
				assert.Len(t, values, 5)
				// x=0: angle = (2pi/10)*(0) = 0. sin(0)=0. res = 100*0 + 1 = 1
				assert.InEpsilon(t, 1.0, values[0], 1e-9)
			},
		},
		{
			name:     "custom sine wave",
			valueStr: "10,5,0,0", // ampl=10, b=5, vShift=0, hShift=0
			count:    5,
			verify: func(t *testing.T, values []float64) {
				assert.Len(t, values, 5)
				// b=5 => period=2pi/5.
				// x=0: angle=0, sin=0, res=0
				assert.InDelta(t, 0.0, values[0], 1e-9)
				// x=1: angle=2pi/5 * 1. sin(72deg) approx 0.951. res=10*0.951=9.51
				expected := 10 * math.Sin(2*math.Pi/5*1)
				assert.InDelta(t, expected, values[1], 1e-9)
			},
		},
		{
			name:     "invalid value string",
			valueStr: "abc",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			gen, err := newSineGenerator[float64](ctx, tt.valueStr, tt.count)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			var values []float64
			for v := range gen {
				values = append(values, v)
			}
			tt.verify(t, values)
		})
	}
}

func TestSineGenerator_Int64(t *testing.T) {
	tests := []struct {
		name     string
		valueStr string
		count    int
		verify   func(t *testing.T, values []int64)
	}{
		{
			name:     "int64 sine",
			valueStr: "100,10,0,0",
			count:    5,
			verify: func(t *testing.T, values []int64) {
				assert.Len(t, values, 5)
				// x=0 => 0
				assert.Equal(t, int64(0), values[0])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			gen, err := newSineGenerator[int64](ctx, tt.valueStr, tt.count)
			require.NoError(t, err)

			var values []int64
			for v := range gen {
				values = append(values, v)
			}
			tt.verify(t, values)
		})
	}
}
