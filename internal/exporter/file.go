package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var _ sdkmetric.Exporter = (*fileExporter)(nil)

type fileExporter struct {
	filePath            string
	mutex               sync.Mutex
	temporalitySelector sdkmetric.TemporalitySelector
}

func (e *fileExporter) Export(_ context.Context, rm *metricdata.ResourceMetrics) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	data, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metrics data: %w", err)
	}

	file, err := os.OpenFile(e.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", e.filePath, err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", e.filePath, err)
	}

	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline to file %s: %w", e.filePath, err)
	}

	return nil
}

func (e *fileExporter) Aggregation(ik sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(ik)
}

func (e *fileExporter) Temporality(ik sdkmetric.InstrumentKind) metricdata.Temporality {
	return e.temporalitySelector(ik)
}

func (e *fileExporter) ForceFlush(_ context.Context) error {
	return nil
}

func (e *fileExporter) Shutdown(_ context.Context) error {
	return nil
}

func NewFileExporter(filePath string, ts sdkmetric.TemporalitySelector) (sdkmetric.Exporter, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	return &fileExporter{
		filePath:            filePath,
		temporalitySelector: ts,
	}, nil
}
