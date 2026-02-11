package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Validation(t *testing.T) {
	// Setup temporary directory for config files
	tmpDir := t.TempDir()

	szgenDir := filepath.Join(tmpDir, "szgen")
	err := os.MkdirAll(szgenDir, 0o755)
	require.NoError(t, err)

	otelConfigPath := filepath.Join(szgenDir, "opentelemetry.yaml")

	t.Run("Wrapped Config (Legacy)", func(t *testing.T) {
		// Create a legacy wrapped config
		configContent := `
opentelemetry:
  file_format: "0.1"
  foo: bar
`
		err := os.WriteFile(otelConfigPath, []byte(configContent), 0o644)
		require.NoError(t, err)

		cmd := &cobra.Command{}
		// Initialize flags that loadConfig reads
		cmd.Flags().String("config", otelConfigPath, "")
		cmd.Flags().String("service-version", "", "")
		cmd.Flags().StringToString("resource-attributes", nil, "")
		cmd.Flags().String("executor", "", "")
		cmd.Flags().Int("max-concurrency", 0, "")

		cfg, err := loadConfig(cmd)
		require.NoError(t, err)
		require.NotNil(t, cfg.OpenTelemetry)

		assert.Equal(t, "0.1", cfg.OpenTelemetry["file_format"])
		assert.Equal(t, "bar", cfg.OpenTelemetry["foo"])
	})
}
