package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neonmei/szgen/internal/runner/metrictask"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/otel"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/executors"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Execute configuration file",
	Long:    `Execute a YAML configuration file with multiple metric generation tasks.`,
	RunE:    runConfigFile,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runConfigFile(cmd *cobra.Command, _ []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	if configPath == "" {
		return fmt.Errorf("config file path is required. Use --config flag or set default location")
	}

	cfg, err := config.LoadConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	ctx, cancelFn := setupContext(context.Background())
	defer cancelFn()

	// Metrics tasks for now
	tasks := make([]runner.Task, 0, len(cfg.Metrics.Tasks))
	for i, metricCfg := range cfg.Metrics.Tasks {
		slog.Info("Queued task",
			"name", metricCfg.Name,
			"generator", metricCfg.Generator,
			"kind", metricCfg.Kind,
			"type", metricCfg.Type,
		)

		task, err := metrictask.New(ctx, metricCfg)
		if err != nil {
			return fmt.Errorf("failed to create task %d: %w", i+1, err)
		}

		tasks = append(tasks, task)
	}

	sdk, err := otel.Start(*cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry SDK: %w", err)
	}
	defer sdk.Shutdown(ctx)

	slog.Debug("Loaded configuration", "task_count", len(cfg.Metrics.Tasks))

	exec, err := executors.New(cfg.Executor)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	if err := exec.Execute(ctx, tasks); err != nil {
		return err
	}

	// Force flush metrics before shutdown
	slog.Info("Flushing metrics")
	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sdk.ForceFlush(flushCtx); err != nil {
		slog.Warn("Failed to flush metrics", "error", err)
	}

	return nil
}
