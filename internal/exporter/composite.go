package exporter

import (
	"context"
	"errors"
	"fmt"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var _ sdkmetric.Exporter = (*compositeExporter)(nil)

type compositeExporter struct {
	exporters []sdkmetric.Exporter
}

func NewCompositeExporter(exporters ...sdkmetric.Exporter) sdkmetric.Exporter {
	return &compositeExporter{exporters: exporters}
}

func (e *compositeExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) (errs error) {
	for _, exporter := range e.exporters {
		if err := exporter.Export(ctx, rm); err != nil {
			errs = errors.Join(errs, fmt.Errorf("composite exporter failed: %w", err))
		}
	}

	return errs
}

func (e *compositeExporter) ForceFlush(ctx context.Context) (errs error) {
	for _, exporter := range e.exporters {
		if err := exporter.ForceFlush(ctx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("composite flush failed: %w", err))
		}
	}

	return errs
}

func (e *compositeExporter) Shutdown(ctx context.Context) (errs error) {
	for _, exporter := range e.exporters {
		if err := exporter.Shutdown(ctx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("composite shutdown failed: %w", err))
		}
	}

	return errs
}

func (e *compositeExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	if len(e.exporters) == 0 {
		return sdkmetric.AggregationDefault{}
	}

	return e.exporters[0].Aggregation(kind)
}

func (e *compositeExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	if len(e.exporters) == 0 {
		return metricdata.DeltaTemporality
	}

	return e.exporters[0].Temporality(kind)
}
