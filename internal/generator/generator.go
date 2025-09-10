package generator

import (
	"context"
	"fmt"
	"iter"

	"github.com/neonmei/szgen/internal/consts"
)

type ValueGenerator[T ~int64 | ~float64] iter.Seq[T]

func New[T int64 | float64](ctx context.Context, pattern string, value string, count int) (ValueGenerator[T], error) {
	switch pattern {
	case consts.GeneratorConstant:
		return newConstantGenerator[T](ctx, value, count)
	case consts.GeneratorRandom:
		return newRandomGenerator[T](ctx, value, count)
	case consts.GeneratorStep:
		return newStepGenerator[T](ctx, value, count)
	case consts.GeneratorSine:
		return newSineGenerator[T](ctx, value, count)
	case consts.GeneratorSequence:
		return newSequenceGenerator[T](ctx, value, count)
	default:
		return nil, fmt.Errorf("unknown generator pattern: %s", pattern)
	}
}
