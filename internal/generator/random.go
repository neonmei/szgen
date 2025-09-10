package generator

import (
	"context"
	"fmt"
	"math/rand/v2"
)

func newRandomGenerator[T int64 | float64](ctx context.Context, valueStr string, count int) (ValueGenerator[T], error) {
	values, err := parseRange[T](valueStr)
	if err != nil {
		return nil, err
	}

	maxVal := values[0]
	minVal := T(0)

	if len(values) > 1 {
		minVal = values[1]
	}

	if minVal >= maxVal {
		return nil, fmt.Errorf("min value %v must be less than max value %v", minVal, maxVal)
	}

	return func(yield func(T) bool) {
		for range count {
			select {
			case <-ctx.Done():
				return
			default:
				var value T
				switch any(value).(type) {
				case int64:
					diff := int64(maxVal) - int64(minVal)
					val := int64(minVal) + rand.Int64N(diff+1)
					value = T(val)
				case float64:
					diff := float64(maxVal) - float64(minVal)
					val := float64(minVal) + rand.Float64()*diff
					value = T(val)
				}

				if !yield(value) {
					return
				}
			}
		}
	}, nil
}
