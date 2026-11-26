package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withGlobalOtelConfig points the global otelconf file at content, so tests do not depend on the user's file.
func withGlobalOtelConfig(t *testing.T, content string) {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "szgen"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "szgen", "opentelemetry.yaml"), []byte(content), 0o644))
}

func ptr[T any](v T) *T {
	return &v
}

func TestExecutorOverrides(t *testing.T) {
	fileConfig := func() config.ExecutorConfig {
		return config.NewExecutorConfig(
			config.WithExecutorStrategy(consts.ExecutorStrategyConcurrent),
			config.WithExecutorParams(map[string]any{consts.ParamMaxConcurrency: 4}),
		)
	}

	t.Run("no overrides keep config", func(t *testing.T) {
		ec := fileConfig()
		ExecutorOverrides{}.apply(&ec)

		assert.Equal(t, fileConfig(), ec)
	})

	t.Run("overrides replace config", func(t *testing.T) {
		ec := fileConfig()
		ExecutorOverrides{Strategy: ptr(consts.ExecutorStrategySerial), MaxConcurrency: ptr(0)}.apply(&ec)

		assert.Equal(t, consts.ExecutorStrategySerial, ec.Strategy)
		assert.Equal(t, 0, ec.Params[consts.ParamMaxConcurrency])
	})

	t.Run("nil params", func(t *testing.T) {
		ec := config.ExecutorConfig{Strategy: consts.ExecutorStrategySerial}
		ExecutorOverrides{MaxConcurrency: ptr(2)}.apply(&ec)

		assert.Equal(t, 2, ec.Params[consts.ParamMaxConcurrency])
	})
}

func TestLoad(t *testing.T) {
	withGlobalOtelConfig(t, "log_level: debug\n")

	cfg, err := Load(
		Options{Version: "1.2.3", Executor: ExecutorOverrides{Strategy: ptr(consts.ExecutorStrategyConcurrent)}},
		config.WithMetricTasks(*config.NewMetricTask()),
	)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.OpenTelemetry["log_level"], "global otelconf file is merged")
	assert.Contains(t, cfg.OpenTelemetry["resource"].(map[string]any)["attributes"],
		map[string]any{"name": "service.version", "value": "1.2.3", "type": "string"})
	assert.Equal(t, consts.ExecutorStrategyConcurrent, cfg.Executor.Strategy)
	assert.Len(t, cfg.Metrics.Tasks, 1)

	t.Run("layer error", func(t *testing.T) {
		_, err := Load(Options{}, config.WithSzgenConfigFile("/nonexist.yaml"))
		assert.ErrorContains(t, err, "failed to load config: failed to read szgen config file /nonexist.yaml")
	})
}

func TestExecute(t *testing.T) {
	withGlobalOtelConfig(t, "disabled: true\n")

	t.Run("invalid configuration", func(t *testing.T) {
		cfg, err := Load(Options{})
		require.NoError(t, err)

		assert.EqualError(t, Execute(cfg), "invalid configuration: no metrics defined in configuration")
	})

	t.Run("runs configured tasks", func(t *testing.T) {
		task := config.NewMetricTask(config.WithRate(1), config.WithCount(2))
		cfg, err := Load(Options{}, config.WithMetricTasks(*task))
		require.NoError(t, err)

		assert.NoError(t, Execute(cfg))
	})

	t.Run("task creation error", func(t *testing.T) {
		task := config.NewMetricTask(config.WithType(consts.ValueTypeInt64), config.WithValue("1.5"))
		cfg, err := Load(Options{}, config.WithMetricTasks(*task))
		require.NoError(t, err)

		assert.EqualError(t, Execute(cfg), "failed to create task 1: create counter iterator: invalid int64 value: 1.5")
	})
}

func TestExecuteTasks(t *testing.T) {
	withGlobalOtelConfig(t, "disabled: true\n")

	t.Run("runs given tasks without config tasks", func(t *testing.T) {
		cfg, err := Load(Options{})
		require.NoError(t, err)

		task := &mocks.MockTask{NameVal: "mock"}
		require.NoError(t, ExecuteTasks(cfg, task))
		assert.Equal(t, 1, task.ExecuteCalled)
	})

	t.Run("task error is returned", func(t *testing.T) {
		cfg, err := Load(Options{})
		require.NoError(t, err)

		err = ExecuteTasks(cfg, &mocks.MockTask{NameVal: "mock", ExecuteErr: assert.AnError})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("invalid opentelemetry configuration", func(t *testing.T) {
		cfg, err := Load(Options{})
		require.NoError(t, err)
		cfg.OpenTelemetry = nil

		assert.EqualError(t, ExecuteTasks(cfg), "invalid configuration: no opentelemetry configuration found")
	})
}

func TestSignalContextCancel(t *testing.T) {
	ctx, cancel := signalContext(context.Background())
	cancel()

	assert.ErrorIs(t, ctx.Err(), context.Canceled)
	_, bySignal := runner.SignalCause(ctx)
	assert.False(t, bySignal)
}
