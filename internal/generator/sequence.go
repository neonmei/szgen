package generator

import (
	"context"
	"fmt"
)

func newSequenceGenerator[T int64 | float64](ctx context.Context, valueStr string, count int) (ValueGenerator[T], error) {
	values, err := parseRange[T](valueStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sequence values: %w", err)
	}

	return func(yield func(T) bool) {
		maxIterations := min(count, len(values))
		for i := range maxIterations {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(values[i]) {
					return
				}
			}
		}
	}, nil
}
