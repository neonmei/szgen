package config

import (
	"fmt"

	"github.com/neonmei/szgen/internal/consts"
)

type MetricViewOption func(*MetricView)

type MetricView struct {
	Instrument ViewInstrument `yaml:"instrument"`
	Stream     ViewStream     `yaml:"stream"`
}

type ViewInstrument struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
}

type ViewStream struct {
	Aggregation ViewAggregation `yaml:"aggregation"`
}

type ViewAggregation struct {
	Kind       string    `yaml:"kind"`
	MaxScale   int       `yaml:"max_scale,omitempty"`
	MaxSize    int       `yaml:"max_size,omitempty"`
	NoMinMax   bool      `yaml:"no_min_max,omitempty"`
	Boundaries []float64 `yaml:"buckets,omitempty"`
}

func (mv *MetricView) Validate() error {
	if mv.Instrument.Name == "" && mv.Instrument.Kind == "" {
		return fmt.Errorf("instrument matching empty, should use either kind or name")
	}

	if err := ValidateInstrumentKind(mv.Instrument.Kind); err != nil {
		return fmt.Errorf("invalid instrument: %w", err)
	}

	switch mv.Stream.Aggregation.Kind {
	case consts.AggregationExponentialHistogram:
		if mv.Stream.Aggregation.MaxScale < 0 || mv.Stream.Aggregation.MaxScale > 20 {
			return fmt.Errorf("max_scale must be between 0 and 20")
		}
		if mv.Stream.Aggregation.MaxSize <= 0 {
			return fmt.Errorf("max_size must be greater than 0")
		}
	case consts.AggregationExplicitBucketHistogram:
		if len(mv.Stream.Aggregation.Boundaries) == 0 {
			return fmt.Errorf("buckets cannot be empty for explicit bucket histogram")
		}
	default:
		return fmt.Errorf("unsupported aggregation kind: %s", mv.Stream.Aggregation.Kind)
	}

	return nil
}

func NewMetricView(opts ...MetricViewOption) *MetricView {
	mv := &MetricView{
		Instrument: ViewInstrument{
			Name: "*",
			Kind: consts.InstrumentKindUndefined,
		},
		Stream: ViewStream{
			Aggregation: ViewAggregation{},
		},
	}

	for _, opt := range opts {
		opt(mv)
	}

	return mv
}

func WithInstrumentName(name string) MetricViewOption {
	return func(mv *MetricView) {
		mv.Instrument.Name = name
	}
}

func WithInstrumentKind(kind string) MetricViewOption {
	return func(mv *MetricView) {
		mv.Instrument.Kind = kind
	}
}

func WithExplicitBuckets(buckets []float64) MetricViewOption {
	return func(mv *MetricView) {
		mv.Stream.Aggregation.Kind = consts.AggregationExplicitBucketHistogram
		mv.Stream.Aggregation.Boundaries = buckets
	}
}

func WithExponentialHistogram(maxScale, maxSize int, noMinMax bool) MetricViewOption {
	return func(mv *MetricView) {
		mv.Stream.Aggregation.Kind = consts.AggregationExponentialHistogram
		mv.Stream.Aggregation.MaxScale = maxScale
		mv.Stream.Aggregation.MaxSize = maxSize
		mv.Stream.Aggregation.NoMinMax = noMinMax
	}
}
