package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/telemetry/metrics"
)

// buildTasks creates a runnable task for every task defined in the configuration, across all signals.
// The context here allows cancelling generation at the producer (i.e: value generator) level.
func buildTasks(ctx context.Context, cfg *config.Config) ([]runner.Task, error) {
	var tasks []runner.Task

	if cfg.Metrics != nil {
		for _, metricCfg := range cfg.Metrics.Tasks {
			slog.Info(
				"Queued task",
				"metric", metricCfg.Name,
				"generator", metricCfg.Generator,
				"kind", metricCfg.Kind,
				"type", metricCfg.Type,
			)

			task, err := metrics.NewTask(ctx, metricCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to create task %d: %w", len(tasks)+1, err)
			}

			tasks = append(tasks, task)
		}
	}

	slog.Debug("Loaded configuration", "task_count", len(tasks))

	return tasks, nil
}
