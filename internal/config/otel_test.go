package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeOTelConfig(t *testing.T) {
	t.Run("nested maps are merged", func(t *testing.T) {
		dst := map[string]any{
			"meter_provider": map[string]any{"readers": []any{"otlp"}, "views": []any{"expo"}},
			"file_format":    "1.0",
		}
		src := map[string]any{
			"meter_provider": map[string]any{"views": []any{"mine"}},
		}

		got := MergeOTelConfig(dst, src)
		assert.Equal(t, map[string]any{
			"meter_provider": map[string]any{"readers": []any{"otlp"}, "views": []any{"mine"}},
			"file_format":    "1.0",
		}, got)
	})

	t.Run("lists and scalars replace", func(t *testing.T) {
		dst := map[string]any{"resource": map[string]any{"attributes": []any{"a", "b"}}, "disabled": false}
		src := map[string]any{"resource": map[string]any{"attributes": []any{}}, "disabled": true}

		got := MergeOTelConfig(dst, src)
		assert.Equal(t, map[string]any{"resource": map[string]any{"attributes": []any{}}, "disabled": true}, got)
	})

	t.Run("map replaces non-map and null is ignored", func(t *testing.T) {
		dst := map[string]any{"a": "scalar", "b": map[string]any{"x": 1}}
		src := map[string]any{"a": map[string]any{"y": 2}, "b": nil}

		got := MergeOTelConfig(dst, src)
		assert.Equal(t, map[string]any{"a": map[string]any{"y": 2}, "b": map[string]any{"x": 1}}, got)
	})

	t.Run("nil dst", func(t *testing.T) {
		got := MergeOTelConfig(nil, map[string]any{"a": 1})
		assert.Equal(t, map[string]any{"a": 1}, got)
	})
}

func TestWithSzgenConfigFile_MergesOTelConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	yamlData := `
opentelemetry:
  meter_provider:
    views:
      - selector:
          instrument_name: "custom.*"
        stream:
          aggregation:
            sum: {}
metrics:
  tasks:
    - name: test.metric
`
	require.NoError(t, os.WriteFile(configFile, []byte(yamlData), 0o644))

	cfg, err := NewConfig(WithDefaultConfig("test"), WithSzgenConfigFile(configFile))
	require.NoError(t, err)

	meterProvider, ok := cfg.OpenTelemetry["meter_provider"].(map[string]any)
	require.True(t, ok)
	assert.NotEmpty(t, meterProvider["readers"], "default readers should be kept")
	assert.Len(t, meterProvider["views"], 1, "views should be replaced")
	assert.Equal(t, "1.0", cfg.OpenTelemetry["file_format"])
	require.NoError(t, cfg.Validate())
}
