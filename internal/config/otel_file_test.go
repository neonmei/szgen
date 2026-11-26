package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithOtelConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	szgenDir := filepath.Join(tmpDir, "szgen")
	require.NoError(t, os.MkdirAll(szgenDir, 0o755))
	otelConfigPath := filepath.Join(szgenDir, "opentelemetry.yaml")

	t.Run("merges into defaults", func(t *testing.T) {
		configContent := `
meter_provider:
  readers:
    - periodic:
        exporter:
          console: {}
`
		require.NoError(t, os.WriteFile(otelConfigPath, []byte(configContent), 0o644))

		cfg, err := NewConfig(WithDefaultConfig("test"), WithOtelConfigFile())
		require.NoError(t, err)

		mp := cfg.OpenTelemetry["meter_provider"].(map[string]any)
		assert.Len(t, mp["readers"], 1)
		assert.NotEmpty(t, mp["views"], "default views are kept")
		assert.Equal(t, "1.0", cfg.OpenTelemetry["file_format"])
	})

	t.Run("invalid yaml", func(t *testing.T) {
		require.NoError(t, os.WriteFile(otelConfigPath, []byte("- not a map"), 0o644))

		_, err := NewConfig(WithDefaultConfig("test"), WithOtelConfigFile())
		assert.ErrorContains(t, err, "failed to unmarshal otelconf file "+otelConfigPath)
	})

	t.Run("missing file", func(t *testing.T) {
		require.NoError(t, os.Remove(otelConfigPath))

		cfg, err := NewConfig(WithDefaultConfig("test"), WithOtelConfigFile())
		require.NoError(t, err)
		assert.Equal(t, NewOTelConfig("test"), cfg.OpenTelemetry)
	})
}

func TestWithMetricTasks(t *testing.T) {
	t.Run("appends to existing tasks", func(t *testing.T) {
		cfg, err := NewConfig(
			WithMetricsConfig(&MetricsConfig{Tasks: []MetricTask{{Name: "first"}}}),
			WithMetricTasks(MetricTask{Name: "second"}),
		)
		require.NoError(t, err)
		assert.Equal(t, []MetricTask{{Name: "first"}, {Name: "second"}}, cfg.Metrics.Tasks)
	})

	t.Run("nil metrics", func(t *testing.T) {
		cfg, err := NewConfig(WithMetricTasks(MetricTask{Name: "only"}))
		require.NoError(t, err)
		assert.Equal(t, []MetricTask{{Name: "only"}}, cfg.Metrics.Tasks)
	})
}
