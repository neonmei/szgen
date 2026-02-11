package metrictask

import (
	"context"
	"fmt"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/generator"
	"github.com/neonmei/szgen/internal/runner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func newInstrument[T int64 | float64](ctx context.Context, cfg config.MetricTask) (runner.Task, error) {
	attr := parseAttributes(cfg.Attributes)
	meter := otel.Meter(consts.DefaultMeterName)

	iter, err := generator.New[T](ctx, cfg.Generator, cfg.Value, cfg.Count)
	if err != nil {
		return nil, fmt.Errorf("create %s iterator: %w", cfg.Kind, err)
	}

	var rec any
	switch any(T(0)).(type) {
	case int64:
		rec, err = newInt64Recorder(meter, cfg, attr)
	case float64:
		rec, err = newFloat64Recorder(meter, cfg, attr)
	default:
		return nil, fmt.Errorf("unsupported numeric type")
	}
	if err != nil {
		return nil, err
	}

	return &metricTask[T]{
		taskName:    cfg.Name,
		genInterval: cfg.Rate,
		genIter:     iter,
		recorder:    rec.(valueRecorder[T]),
	}, nil
}

func newInt64Recorder(m metric.Meter, cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[int64], error) {
	desc := metric.WithDescription(cfg.Description)
	unit := metric.WithUnit(cfg.Unit)
	withAttr := metric.WithAttributes(attr...)

	switch cfg.Kind {
	case consts.MetricTypeCounter:
		c, err := m.Int64Counter(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create int64 counter %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v int64) { c.Add(ctx, v, withAttr) }, nil

	case consts.MetricTypeGauge:
		c, err := m.Int64Gauge(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create int64 gauge %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v int64) { c.Record(ctx, v, withAttr) }, nil

	case consts.MetricTypeHistogram:
		c, err := m.Int64Histogram(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create int64 histogram %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v int64) { c.Record(ctx, v, withAttr) }, nil

	case consts.MetricTypeUpDownCounter:
		c, err := m.Int64UpDownCounter(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create int64 updowncounter %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v int64) { c.Add(ctx, v, withAttr) }, nil

	default:
		return nil, fmt.Errorf("unsupported metric kind: %s", cfg.Kind)
	}
}

func newFloat64Recorder(m metric.Meter, cfg config.MetricTask, attr []attribute.KeyValue) (valueRecorder[float64], error) {
	desc := metric.WithDescription(cfg.Description)
	unit := metric.WithUnit(cfg.Unit)
	withAttr := metric.WithAttributes(attr...)

	switch cfg.Kind {
	case consts.MetricTypeCounter:
		c, err := m.Float64Counter(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create float64 counter %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v float64) { c.Add(ctx, v, withAttr) }, nil

	case consts.MetricTypeGauge:
		c, err := m.Float64Gauge(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create float64 gauge %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v float64) { c.Record(ctx, v, withAttr) }, nil

	case consts.MetricTypeHistogram:
		c, err := m.Float64Histogram(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create float64 histogram %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v float64) { c.Record(ctx, v, withAttr) }, nil

	case consts.MetricTypeUpDownCounter:
		c, err := m.Float64UpDownCounter(cfg.Name, desc, unit)
		if err != nil {
			return nil, fmt.Errorf("create float64 updowncounter %q: %w", cfg.Name, err)
		}
		return func(ctx context.Context, v float64) { c.Add(ctx, v, withAttr) }, nil

	default:
		return nil, fmt.Errorf("unsupported metric kind: %s", cfg.Kind)
	}
}
