package otel

import (
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"go.opentelemetry.io/otel"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func buildMetricsProvider(cfg config.Config, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, error) {
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(NewResource(&cfg.Resource)),
	}

	for _, viewCfg := range cfg.Metrics.Views {
		if err := viewCfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid view configuration: %w", err)
		}

		view, err := newMetricView(viewCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create view: %w", err)
		}

		opts = append(opts, sdkmetric.WithView(view))
	}

	slog.Info("Configuring OTLP exporter",
		"temporality", cfg.Export.Temporality,
		"endpoint", cfg.Export.Endpoint,
		"insecure", cfg.Export.Insecure,
		"interval", cfg.Export.Interval,
		"mode", cfg.Export.Mode,
	)

	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(cfg.Export.Interval),
	)

	opts = append(opts, sdkmetric.WithReader(reader))

	provider := sdkmetric.NewMeterProvider(opts...)
	otel.SetMeterProvider(provider)

	return provider, nil
}

func newMetricView(viewCfg config.MetricView) (sdkmetric.View, error) {
	kind, err := parseInstrumentKind(viewCfg.Instrument.Kind)
	if err != nil {
		return nil, fmt.Errorf("invalid instrument kind: %w", err)
	}

	instrument := sdkmetric.Instrument{
		Name: viewCfg.Instrument.Name,
		Kind: kind,
	}

	var aggregation sdkmetric.Aggregation
	switch viewCfg.Stream.Aggregation.Kind {
	case consts.AggregationExponentialHistogram:
		aggregation = sdkmetric.AggregationBase2ExponentialHistogram{
			MaxSize:  int32(viewCfg.Stream.Aggregation.MaxSize),
			MaxScale: int32(viewCfg.Stream.Aggregation.MaxScale),
			NoMinMax: viewCfg.Stream.Aggregation.NoMinMax,
		}
	case consts.AggregationExplicitBucketHistogram:
		aggregation = sdkmetric.AggregationExplicitBucketHistogram{
			Boundaries: viewCfg.Stream.Aggregation.Boundaries,
			NoMinMax:   viewCfg.Stream.Aggregation.NoMinMax,
		}
	default:
		return nil, fmt.Errorf("unsupported aggregation kind: %s", viewCfg.Stream.Aggregation.Kind)
	}

	stream := sdkmetric.Stream{
		Aggregation: aggregation,
	}

	return sdkmetric.NewView(instrument, stream), nil
}

func temporalitySelector(exportCfg config.ExportConfig) sdkmetric.TemporalitySelector {
	return func(_ sdkmetric.InstrumentKind) metricdata.Temporality {
		switch exportCfg.Temporality {
		case consts.TemporalityCumulative:
			return metricdata.CumulativeTemporality
		case consts.TemporalityDelta:
			return metricdata.DeltaTemporality
		default:
			return metricdata.DeltaTemporality
		}
	}
}

func parseInstrumentKind(kind string) (sdkmetric.InstrumentKind, error) {
	switch kind {
	case consts.InstrumentKindUndefined:
		return sdkmetric.InstrumentKind(0), nil
	case consts.InstrumentKindCounter:
		return sdkmetric.InstrumentKindCounter, nil
	case consts.InstrumentKindUpDownCounter:
		return sdkmetric.InstrumentKindUpDownCounter, nil
	case consts.InstrumentKindHistogram:
		return sdkmetric.InstrumentKindHistogram, nil
	case consts.InstrumentKindGauge:
		return sdkmetric.InstrumentKindGauge, nil
	case consts.InstrumentKindObservableCounter:
		return sdkmetric.InstrumentKindObservableCounter, nil
	case consts.InstrumentKindObservableUpDown:
		return sdkmetric.InstrumentKindObservableUpDownCounter, nil
	case consts.InstrumentKindObservableGauge:
		return sdkmetric.InstrumentKindObservableGauge, nil
	default:
		return sdkmetric.InstrumentKind(0), fmt.Errorf("unknown instrument kind: %s", kind)
	}
}
