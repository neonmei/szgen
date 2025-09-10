package metrictask

import (
	"context"
	"fmt"

	"github.com/neonmei/szgen/internal/runner"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/generator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func newInstrument[T int64 | float64](ctx context.Context, cfg config.MetricTask) (runner.Task, error) {
	attr := parseAttributes(cfg.Attributes)

	iter, err := generator.New[T](ctx, cfg.Generator, cfg.Value, cfg.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to create int64 iterator: %w", err)
	}

	var recorder valueRecorder[T]
	switch cfg.Kind {
	case consts.MetricTypeCounter:
		recorder, err = newCounter[T](cfg, attr)
	case consts.MetricTypeGauge:
		recorder, err = newGauge[T](cfg, attr)
	case consts.MetricTypeHistogram:
		recorder, err = newHistogram[T](cfg, attr)
	case consts.MetricTypeUpDownCounter:
		recorder, err = newUpDownCounter[T](cfg, attr)
	default:
		return nil, fmt.Errorf("unsupported metric kind: %s", cfg.Kind)
	}

	if err != nil {
		return nil, err
	}

	return &metricTask[T]{
		taskName:    cfg.Name,
		genInterval: cfg.Rate,
		genIter:     iter,
		recorder:    recorder,
	}, nil
}

func newCounter[T int64 | float64](cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[T], error) {
	meter := otel.Meter("szgen")

	switch any(T(0)).(type) {
	case int64:
		c, err := meter.Int64Counter(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Add(ctx, int64(value), metric.WithAttributes(attr...)) }, nil

	case float64:
		c, err := meter.Float64Counter(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Add(ctx, float64(value), metric.WithAttributes(attr...)) }, nil

	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func newGauge[T int64 | float64](cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[T], error) {
	meter := otel.Meter("szgen")

	switch any(T(0)).(type) {
	case int64:
		c, err := meter.Int64Gauge(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Record(ctx, int64(value), metric.WithAttributes(attr...)) }, nil

	case float64:
		c, err := meter.Float64Gauge(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Record(ctx, float64(value), metric.WithAttributes(attr...)) }, nil

	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func newUpDownCounter[T int64 | float64](cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[T], error) {
	meter := otel.Meter("szgen")

	switch any(T(0)).(type) {
	case int64:
		c, err := meter.Int64UpDownCounter(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Add(ctx, int64(value), metric.WithAttributes(attr...)) }, nil

	case float64:
		c, err := meter.Float64UpDownCounter(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Add(ctx, float64(value), metric.WithAttributes(attr...)) }, nil

	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func newHistogram[T int64 | float64](cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[T], error) {
	meter := otel.Meter("szgen")

	switch any(T(0)).(type) {
	case int64:
		c, err := meter.Int64Histogram(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Record(ctx, int64(value), metric.WithAttributes(attr...)) }, nil

	case float64:
		c, err := meter.Float64Histogram(cfg.Name, metric.WithDescription(cfg.Description), metric.WithUnit(cfg.Unit))
		if err != nil {
			return nil, fmt.Errorf("failed to create int64 counter instrument: %w", err)
		}

		return func(ctx context.Context, value T) { c.Record(ctx, float64(value), metric.WithAttributes(attr...)) }, nil

	default:
		return nil, fmt.Errorf("invalid type")
	}
}
