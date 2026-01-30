package otel

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"go.opentelemetry.io/contrib/otelconf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"gopkg.in/yaml.v3"
)

type SDK struct {
	cfg *otelconf.OpenTelemetryConfiguration
	sdk *otelconf.SDK
}

func NewSDK(cfg *config.Config) (*SDK, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	bytes, err := yaml.Marshal(cfg.OpenTelemetry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal opentelemetry config: %w", err)
	}

	slog.Debug("Instantiated OpenTelemetry SDK config", "config", string(bytes))
	conf, err := otelconf.ParseYAML(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse opentelemetry config: %w", err)
	}

	return &SDK{cfg: conf}, nil
}

func (s *SDK) Start() error {
	sdk, err := otelconf.NewSDK(otelconf.WithOpenTelemetryConfiguration(*s.cfg))
	if err != nil {
		return fmt.Errorf("failed to create otel sdk: %w", err)
	}
	s.sdk = &sdk

	otel.SetTracerProvider(sdk.TracerProvider())
	otel.SetMeterProvider(sdk.MeterProvider())
	global.SetLoggerProvider(sdk.LoggerProvider())
	return nil
}

func (s *SDK) Shutdown(ctx context.Context) error {
	if s.sdk != nil {
		return s.sdk.Shutdown(ctx)
	}
	return nil
}

func (s *SDK) ForceFlush(ctx context.Context) error {
	if s.sdk == nil {
		return nil
	}
	if provider, ok := s.sdk.MeterProvider().(interface{ ForceFlush(context.Context) error }); ok {
		return provider.ForceFlush(ctx)
	}
	if provider, ok := s.sdk.TracerProvider().(interface{ ForceFlush(context.Context) error }); ok {
		return provider.ForceFlush(ctx)
	}
	if provider, ok := s.sdk.LoggerProvider().(interface{ ForceFlush(context.Context) error }); ok {
		return provider.ForceFlush(ctx)
	}

	return nil
}
