package generator

import (
	"context"
)

func newConstantGenerator[T int64 | float64](ctx context.Context, valueStr string, count int) (ValueGenerator[T], error) {
	value, err := parseValue[T](valueStr)
	if err != nil {
		return nil, err
	}

	return func(yield func(T) bool) {
		for range count {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(value) {
					return
				}
			}
		}
	}, nil
}
