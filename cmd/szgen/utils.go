package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/otel"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/executors"
	"github.com/spf13/cobra"
)

// runTasks starts the OpenTelemetry SDK, executes tasks and flushes pending telemetry before shutting down.
func runTasks(ctx context.Context, cfg *config.Config, tasks []runner.Task) error {
	exec, err := executors.New(cfg.Executor)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	sdk, err := otel.NewSDK(cfg)
	if err != nil {
		return fmt.Errorf("failed to create sdk: %w", err)
	}

	if err := sdk.Start(); err != nil {
		return fmt.Errorf("failed to start sdk: %w", err)
	}
	defer func() {
		// ctx may already be cancelled by a signal, so shutdown gets its own deadline.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), consts.DefaultFlushTimeout)
		defer cancel()
		_ = sdk.Shutdown(shutdownCtx)
	}()

	execErr := exec.Execute(ctx, tasks)

	slog.Info("Flushing metrics")
	flushCtx, cancel := context.WithTimeout(context.Background(), consts.DefaultFlushTimeout)
	defer cancel()

	if err := sdk.ForceFlush(flushCtx); err != nil {
		slog.Warn("Failed to flush metrics", "error", err)
	}

	return execErr
}

func buildMetricConfig(cmd *cobra.Command, metricType string) (*config.MetricTask, error) {
	options := []config.MetricTaskOption{
		config.WithKind(metricType),
	}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		options = append(options, config.WithName(name))
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		options = append(options, config.WithDescription(description))
	}
	if cmd.Flags().Changed("unit") {
		unit, _ := cmd.Flags().GetString("unit")
		options = append(options, config.WithUnit(unit))
	}
	if cmd.Flags().Changed("generator") {
		generator, _ := cmd.Flags().GetString("generator")
		options = append(options, config.WithGenerator(generator))
	}
	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetString("value")
		options = append(options, config.WithValue(value))
	}
	if cmd.Flags().Changed("type") {
		valueType, _ := cmd.Flags().GetString("type")
		options = append(options, config.WithType(valueType))
	}
	if cmd.Flags().Changed("count") {
		count, _ := cmd.Flags().GetInt("count")
		options = append(options, config.WithCount(count))
	}
	if cmd.Flags().Changed("rate") {
		rate, _ := cmd.Flags().GetDuration("rate")
		options = append(options, config.WithRate(rate))
	}
	if cmd.Flags().Changed("attributes") {
		strAttrs, _ := cmd.Flags().GetStringToString("attributes")
		attrs := make(map[string]any, len(strAttrs))
		for k, v := range strAttrs {
			attrs[k] = v
		}
		options = append(options, config.WithMetricAttributes(attrs))
	}

	mc := config.NewMetricTask(options...)

	if err := mc.Validate(); err != nil {
		return nil, err
	}

	slog.Debug(
		"loaded config",
		"metric", mc.Name,
		"kind", mc.Kind,
		"type", mc.Type,
		"rate", mc.Rate,
		"count", mc.Count,
		"value", mc.Value,
		"attributes", mc.Attributes,
		"generator", mc.Generator,
		"description", mc.Description,
		"unit", mc.Unit,
	)
	return mc, nil
}

func setupSignalHandler(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigCh:
			slog.Info("Received signal, shutting down gracefully", "signal", sig)
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
		close(sigCh)
	}()

	return ctx, cancel
}
