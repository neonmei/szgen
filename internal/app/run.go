package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/executors"
	"github.com/neonmei/szgen/internal/telemetry"
)

// Execute validates the configuration and runs every task it defines.
func Execute(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return run(cfg, func(ctx context.Context) ([]runner.Task, error) {
		return buildTasks(ctx, cfg)
	})
}

// ExecuteTasks runs tasks built by the caller. The configuration only sets up the SDK and the executor.
func ExecuteTasks(cfg *config.Config, tasks ...runner.Task) error {
	if err := cfg.ValidateOpenTelemetry(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return run(cfg, func(context.Context) ([]runner.Task, error) {
		return tasks, nil
	})
}

// run builds tasks with a context cancelled by SIGINT/SIGTERM and runs them.
func run(cfg *config.Config, build func(context.Context) ([]runner.Task, error)) error {
	ctx, cancelFn := signalContext(context.Background())
	defer cancelFn()

	tasks, err := build(ctx)
	if err != nil {
		return err
	}

	return runTasks(ctx, cfg, tasks)
}

// runTasks starts the OpenTelemetry SDK, executes tasks and flushes pending telemetry before shutting down.
func runTasks(ctx context.Context, cfg *config.Config, tasks []runner.Task) error {
	exec, err := executors.New(cfg.Executor)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	sdk, err := telemetry.NewSDK(cfg)
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

	slog.Info("Flushing telemetry")
	flushCtx, cancel := context.WithTimeout(context.Background(), consts.DefaultFlushTimeout)
	defer cancel()

	if err := sdk.ForceFlush(flushCtx); err != nil {
		slog.Warn("Failed to flush telemetry", "error", err)
	}

	return execErr
}
