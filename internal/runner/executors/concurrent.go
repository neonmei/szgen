package executors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
	"golang.org/x/sync/errgroup"
)

type concurrentExecutor struct {
	maxConcurrency int
}

func (e *concurrentExecutor) Execute(ctx context.Context, tasks []runner.Task) error {
	if len(tasks) == 0 {
		slog.Info("concurrent executor finished", "tasks", 0)
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	if e.maxConcurrency > 0 {
		g.SetLimit(e.maxConcurrency)
	}

	for _, task := range tasks {
		g.Go(func() error {
			return e.executeWithRecovery(ctx, task)
		})
	}

	if err := g.Wait(); err != nil {
		slog.Error("concurrent executor failed", "error", err)
		return err
	}

	slog.Info("concurrent executor finished successfully", "tasks", len(tasks))
	return nil
}

func (e *concurrentExecutor) executeWithRecovery(ctx context.Context, task runner.Task) (taskErr error) {
	defer func() {
		if r := recover(); r != nil {
			taskErr = fmt.Errorf("task panicked: %v", r)
		}
	}()

	if err := task.Execute(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("task aborted: %w", err)
	}
	return nil
}

func NewConcurrent(params map[string]any) *concurrentExecutor {
	return &concurrentExecutor{
		maxConcurrency: intParam(params, consts.ParamMaxConcurrency, 0),
	}
}

func intParam(params map[string]any, key string, fallback int) int {
	val, ok := params[key]
	if !ok {
		return fallback
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return fallback
	}
}
