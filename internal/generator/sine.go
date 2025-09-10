package generator

import (
	"context"
	"math"

	"github.com/neonmei/szgen/internal/consts"
)

func newSineGenerator[T int64 | float64](ctx context.Context, valueStr string, count int) (ValueGenerator[T], error) {
	sineParams, err := parseRange[T](valueStr)
	if err != nil {
		return nil, err
	}

	amplitude := sineParams[0]
	b := T(consts.DefaultSineGeneratorB)
	verticalShift := T(1)
	horizontalShift := T(0)

	if len(sineParams) > 1 {
		b = sineParams[1]
	}

	if len(sineParams) > 2 {
		verticalShift = sineParams[2]
	}

	if len(sineParams) > 3 {
		horizontalShift = sineParams[3]
	}

	return func(yield func(T) bool) {
		for x := range count {
			select {
			case <-ctx.Done():
				return
			default:
				period := 2 * math.Pi / float64(b)
				angle := period * (float64(x) + float64(horizontalShift))
				sineValue := math.Sin(angle)
				result := T(float64(amplitude)*sineValue + float64(verticalShift))

				if !yield(result) {
					return
				}
			}
		}
	}, nil
}
