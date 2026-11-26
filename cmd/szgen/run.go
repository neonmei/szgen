package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/runner"
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
	runCmd.Flags().String("config", "", "Configuration file path")
	rootCmd.AddCommand(runCmd)
}

func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	configPath, _ := cmd.Flags().GetString("config")
	if configPath == "" {
		return nil, errors.New("no config file path provided")
	}

	return config.NewConfig(
		config.WithDefaultConfig(version),
		config.WithOtelConfigFile(),
		config.WithSzgenConfigFile(configPath),
	)
}

func runConfigFile(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	applyExecutorFlags(cmd, &cfg.Executor)

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancelFn := setupSignalHandler(context.Background())
	defer cancelFn()

	tasks := make([]runner.Task, 0, len(cfg.Metrics.Tasks))
	for i, metricCfg := range cfg.Metrics.Tasks {
		slog.Info(
			"Queued task",
			"metric", metricCfg.Name,
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

	slog.Debug("Loaded configuration", "task_count", len(cfg.Metrics.Tasks))

	return runTasks(ctx, cfg, tasks)
}
