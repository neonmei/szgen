package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/otel"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/executors"
	"github.com/neonmei/szgen/internal/runner/metrictask"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:     "metrics",
	Aliases: []string{"m"},
	Short:   "Generate OpenTelemetry metrics",
	Long:    `Generate various types of OpenTelemetry metrics including counters, gauges, histograms, and updowncounters.`,
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	metricsCmd.PersistentFlags().StringP("attributes", "a", "", "Comma-separated key=value pairs")
	metricsCmd.PersistentFlags().IntP("count", "c", consts.DefaultCount, "Number of data points to generate")
	metricsCmd.PersistentFlags().DurationP("rate", "r", consts.DefaultRate, "Time interval between each generated data point")
	metricsCmd.PersistentFlags().StringP("description", "d", consts.DefaultDescription, "Metric description")
	metricsCmd.PersistentFlags().StringP("name", "n", consts.DefaultMetricName, "Metric name")
	metricsCmd.PersistentFlags().StringP("unit", "u", "", "Metric unit")
	metricsCmd.PersistentFlags().StringP("generator", "g", consts.DefaultGenerator, "Value generation pattern")
	metricsCmd.PersistentFlags().StringP("value", "v", consts.DefaultValue, "Static value or value range")
	metricsCmd.PersistentFlags().StringP("type", "t", "", "Value type (int64, float64) - smart defaults: counter=int64, others=float64")
}

func runMetricCommand(cmd *cobra.Command, metricType string) error {
	exportCfg, err := parseExportConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to build export config: %w", err)
	}

	resourceCfg, err := parseResourceConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to build resource config: %w", err)
	}

	metricCfg, err := buildMetricConfig(cmd, metricType)
	if err != nil {
		return fmt.Errorf("failed to build metric config: %w", err)
	}

	ctx, cancelFn := setupContext(context.Background())
	defer cancelFn()

	cfg := config.Config{
		Metrics: &config.MetricsConfig{
			Tasks: []config.MetricTask{*metricCfg},
		},
		Export:   *exportCfg,
		Resource: *resourceCfg,
	}

	task, err := metrictask.New(ctx, *metricCfg)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	sdk, err := otel.Start(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry SDK: %w", err)
	}
	defer sdk.Shutdown(ctx)

	executorConfig, err := parseExecutorConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse executor config: %w", err)
	}

	exec, err := executors.New(*executorConfig)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	if err := exec.Execute(ctx, []runner.Task{task}); err != nil {
		return err
	}

	slog.Info("Flushing metrics")
	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sdk.ForceFlush(flushCtx); err != nil {
		slog.Warn("Failed to flush metrics", "error", err)
	}

	return nil
}

func runMetricCommandWithViews(cmd *cobra.Command, metricType string, views []config.MetricView) error {
	exportCfg, err := parseExportConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to build export config: %w", err)
	}

	resourceCfg, err := parseResourceConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to build resource config: %w", err)
	}

	metricCfg, err := buildMetricConfig(cmd, metricType)
	if err != nil {
		return fmt.Errorf("failed to build metric config: %w", err)
	}

	ctx, cancelFn := setupContext(context.Background())
	defer cancelFn()

	cfg := config.Config{
		Metrics: &config.MetricsConfig{
			Tasks: []config.MetricTask{*metricCfg},
			Views: views,
		},
		Export:   *exportCfg,
		Resource: *resourceCfg,
	}

	task, err := metrictask.New(ctx, *metricCfg)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	sdk, err := otel.Start(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry SDK: %w", err)
	}
	defer sdk.Shutdown(ctx)

	executorConfig, err := parseExecutorConfigFromCli(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse executor config: %w", err)
	}

	exec, err := executors.New(*executorConfig)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	if err := exec.Execute(ctx, []runner.Task{task}); err != nil {
		return err
	}

	slog.Info("Flushing metrics")
	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sdk.ForceFlush(flushCtx); err != nil {
		slog.Warn("Failed to flush metrics", "error", err)
	}

	return nil
}
