package executors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/runner"
)

type serialExecutor struct{}

func (e *serialExecutor) Execute(ctx context.Context, tasks []runner.Task) error {
	for i, task := range tasks {
		err := ctx.Err()
		if err == nil {
			err = task.Execute(ctx)
		}
		if errors.Is(err, context.Canceled) {
			slog.Info("serial executor interrupted", "completed", i, "tasks", len(tasks))
			return nil
		}
		if err != nil {
			return fmt.Errorf("task %d aborted: %w", i+1, err)
		}
	}

	slog.Info("serial executor finished", "tasks", len(tasks))
	return nil
}

func NewSerial() *serialExecutor {
	return &serialExecutor{}
}
