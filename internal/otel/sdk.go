package otel

import (
	"context"
	"fmt"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/exporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// this is temporary until config fully lands
// see: https://opentelemetry.io/docs/specs/otel/configuration/sdk/
type otelSDK struct {
	meterProvider *sdkmetric.MeterProvider
}

func Start(cfg config.Config) (*otelSDK, error) {
	e, err := newExporter(cfg.Export)
	if err != nil {
		return nil, err
	}

	meterProvider, err := buildMetricsProvider(cfg, e)
	if err != nil {
		return nil, err
	}

	return &otelSDK{
		meterProvider: meterProvider,
	}, nil
}

func (s *otelSDK) ForceFlush(ctx context.Context) error {
	return s.meterProvider.ForceFlush(ctx)
}

func (s *otelSDK) Shutdown(ctx context.Context) error {
	return s.meterProvider.Shutdown(ctx)
}

func newExporter(exportCfg config.ExportConfig) (sdkmetric.Exporter, error) {
	ts := temporalitySelector(exportCfg)

	switch exportCfg.Mode {
	case consts.ExportModeExecute:
		return newGrpcExporter(exportCfg, ts)

	case consts.ExportModeExecuteAndSave:
		otlpExporter, err := newGrpcExporter(exportCfg, ts)
		if err != nil {
			return nil, err
		}

		fileExporter, err := exporter.NewFileExporter(exportCfg.File, ts)
		if err != nil {
			return nil, err
		}

		return exporter.NewCompositeExporter(otlpExporter, fileExporter), nil

	case consts.ExportModeSave:
		return exporter.NewFileExporter(exportCfg.File, ts)

	default:
		return nil, fmt.Errorf("invalid export mode: %s", exportCfg.Mode)
	}
}

func newGrpcExporter(exportCfg config.ExportConfig, t sdkmetric.TemporalitySelector) (sdkmetric.Exporter, error) {
	// TODO: Implement TLS

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(exportCfg.Endpoint),
		otlpmetricgrpc.WithTemporalitySelector(t),
		otlpmetricgrpc.WithInsecure(),
	}

	exporter, err := otlpmetricgrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	return exporter, nil
}
