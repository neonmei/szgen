package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/otel"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/executors"
	"github.com/neonmei/szgen/internal/runner/metrictask"
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
	rootCmd.PersistentFlags().String("config", "", "Configuration file path")
	rootCmd.AddCommand(runCmd)
}

func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	configPath, _ := cmd.Flags().GetString("config")
	if configPath == "" {
		return nil, fmt.Errorf("no config file path provided")
	}

	cfg, err := config.NewConfig(
		config.WithDefaultConfig(),
		config.WithOtelConfigFile(),
		config.WithSzgenConfigFile(configPath),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

func runConfigFile(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if len(cfg.Metrics.Tasks) == 0 {
		return fmt.Errorf("no metric tasks defined in configuration")
	}

	ctx, cancelFn := setupSignalHandler(context.Background())
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

	sdk, err := otel.NewSDK(cfg)
	if err != nil {
		return fmt.Errorf("failed to create sdk: %w", err)
	}

	if err := sdk.Start(); err != nil {
		return fmt.Errorf("failed to start sdk: %w", err)
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
