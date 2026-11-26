package main

import (
	"context"
	"fmt"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
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

	metricsCmd.PersistentFlags().StringToStringP("attributes", "a", nil, "Comma-separated key=value pairs")
	metricsCmd.PersistentFlags().IntP("count", "c", consts.DefaultCount, "Number of data points to generate")
	metricsCmd.PersistentFlags().DurationP("rate", "r", consts.DefaultRate, "Time interval between each generated data point")
	metricsCmd.PersistentFlags().StringP("description", "d", "", "Metric description")
	metricsCmd.PersistentFlags().StringP("name", "n", consts.DefaultMetricName, "Metric name")
	metricsCmd.PersistentFlags().StringP("unit", "u", "", "Metric unit")
	metricsCmd.PersistentFlags().StringP("generator", "g", consts.DefaultGenerator, "Value generation pattern")
	metricsCmd.PersistentFlags().StringP("value", "v", consts.DefaultValue, "Static value or value range")
	metricsCmd.PersistentFlags().StringP("type", "t", consts.DefaultValueType, "Value type (int64, float64)")
}

func runMetricCommand(cmd *cobra.Command, metricType string) error {
	cfg, err := config.NewConfig(config.WithDefaultConfig(version), config.WithOtelConfigFile())
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	metricCfg, err := buildMetricConfig(cmd, metricType)
	if err != nil {
		return fmt.Errorf("failed to build metric config: %w", err)
	}

	ctx, cancelFn := setupSignalHandler(context.Background())
	defer cancelFn()

	// Append CLI task to config
	cfg.Metrics.Tasks = append(cfg.Metrics.Tasks, *metricCfg)
	applyExecutorFlags(cmd, &cfg.Executor)

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	task, err := metrictask.New(ctx, *metricCfg)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return runTasks(ctx, cfg, []runner.Task{task})
}
