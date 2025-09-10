package executors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/runner"
)

type serialExecutor struct{}

func (e *serialExecutor) Execute(ctx context.Context, tasks []runner.Task) error {
	for i, task := range tasks {
		if err := task.Execute(ctx); err != nil {
			return fmt.Errorf("task %d aborted: %w", i+1, err)
		}
	}

	slog.Info("serial executor finished", "tasks", len(tasks))
	return nil
}

func NewSerial() *serialExecutor {
	return &serialExecutor{}
}
