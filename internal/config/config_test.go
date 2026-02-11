package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		cfg, err := NewConfig(WithDefaultConfig("test"))
		require.NoError(t, err)
		assert.NotNil(t, cfg.Metrics)
		assert.Empty(t, cfg.Metrics.Tasks)
		assert.NotNil(t, cfg.OpenTelemetry)
		assert.Equal(t, consts.DefaultExecutorStrategy, cfg.Executor.Strategy)
	})

	t.Run("with metrics config", func(t *testing.T) {
		metrics := &MetricsConfig{
			Tasks: []MetricTask{{Name: "test"}},
		}
		cfg, err := NewConfig(WithMetricsConfig(metrics))
		require.NoError(t, err)
		assert.Equal(t, metrics, cfg.Metrics)
	})

	t.Run("with otel config", func(t *testing.T) {
		otel := map[string]any{"key": "val"}
		cfg, err := NewConfig(WithOpenTelemetryConfig(otel))
		require.NoError(t, err)
		assert.Equal(t, otel, cfg.OpenTelemetry)
	})

	t.Run("with executor config", func(t *testing.T) {
		exec := ExecutorConfig{Strategy: consts.ExecutorStrategyConcurrent}
		cfg, err := NewConfig(WithExecutorConfig(exec))
		require.NoError(t, err)
		assert.Equal(t, exec, cfg.Executor)
	})

	t.Run("options failure", func(t *testing.T) {
		errOption := func(c *Config) error {
			return assert.AnError
		}
		cfg, err := NewConfig(errOption)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, cfg)
	})
}

func TestWithSzgenConfigFile(t *testing.T) {
	t.Run("load valid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		yamlData := `
metrics:
  tasks:
    - name: test.metric
      kind: counter
executor:
  strategy: serial
`
		err := os.WriteFile(configFile, []byte(yamlData), 0o644)
		require.NoError(t, err)

		cfg, err := NewConfig(WithSzgenConfigFile(configFile))
		require.NoError(t, err)

		require.NotNil(t, cfg.Metrics)
		require.Len(t, cfg.Metrics.Tasks, 1)
		assert.Equal(t, "test.metric", cfg.Metrics.Tasks[0].Name)
		assert.Equal(t, consts.ExecutorStrategySerial, cfg.Executor.Strategy)
	})

	t.Run("file not found", func(t *testing.T) {
		cfg, err := NewConfig(WithSzgenConfigFile("nonexistent.yaml"))
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(configFile, []byte("invalid: yaml: content: :"), 0o644)
		require.NoError(t, err)

		cfg, err := NewConfig(WithSzgenConfigFile(configFile))
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("empty path", func(t *testing.T) {
		cfg, err := NewConfig(WithDefaultConfig("test"), WithSzgenConfigFile(""))
		require.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			Metrics: &MetricsConfig{
				Tasks: []MetricTask{{
					Name:      "valid.metric",
					Kind:      consts.MetricTypeCounter,
					Type:      consts.ValueTypeInt64,
					Rate:      time.Second,
					Generator: consts.GeneratorConstant,
				}},
			},
			Executor: ExecutorConfig{Strategy: consts.ExecutorStrategySerial},
			OpenTelemetry: map[string]any{
				"file_format": "1.0",
			},
		}
		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("no tasks", func(t *testing.T) {
		cfg := &Config{
			Metrics: &MetricsConfig{Tasks: []MetricTask{}},
		}
		err := cfg.Validate()
		assert.ErrorContains(t, err, "no metrics defined")
	})

	t.Run("invalid executor", func(t *testing.T) {
		cfg := &Config{
			Metrics:  &MetricsConfig{Tasks: []MetricTask{{Name: "valid"}}},
			Executor: ExecutorConfig{Strategy: "invalid"},
		}
		err := cfg.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid metric task", func(t *testing.T) {
		cfg := &Config{
			Metrics: &MetricsConfig{
				Tasks: []MetricTask{{Name: ""}}, // Invalid
			},
			Executor: ExecutorConfig{Strategy: consts.ExecutorStrategySerial},
		}
		err := cfg.Validate()
		assert.ErrorContains(t, err, "metric[0]")
	})

	t.Run("no otel config", func(t *testing.T) {
		cfg := &Config{
			Metrics: &MetricsConfig{
				Tasks: []MetricTask{{
					Name:      "valid.metric",
					Kind:      consts.MetricTypeCounter,
					Type:      consts.ValueTypeInt64,
					Rate:      time.Second,
					Generator: consts.GeneratorConstant,
				}},
			},
			Executor:      ExecutorConfig{Strategy: consts.ExecutorStrategySerial},
			OpenTelemetry: nil,
		}
		err := cfg.Validate()
		assert.ErrorContains(t, err, "no opentelemetry configuration")
	})
}
