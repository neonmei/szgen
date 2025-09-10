package generator

import (
	"context"
)

func newStepGenerator[T int64 | float64](ctx context.Context, valueStr string, count int) (ValueGenerator[T], error) {
	values, err := parseRange[T](valueStr)
	if err != nil {
		return nil, err
	}

	initial := values[0]
	step := T(1)

	if len(values) > 1 {
		step = values[1]
	}

	return func(yield func(T) bool) {
		current := initial
		for range count {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(current) {
					return
				}
				current += step
			}
		}
	}, nil
}
