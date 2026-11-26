package metricscmd

import (
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMetricTestCmd(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{}
	cmd.Flags().StringToStringP("attributes", "a", nil, "")
	cmd.Flags().IntP("count", "c", consts.DefaultCount, "")
	cmd.Flags().DurationP("rate", "r", consts.DefaultRate, "")
	cmd.Flags().StringP("description", "d", "", "")
	cmd.Flags().StringP("name", "n", consts.DefaultMetricName, "")
	cmd.Flags().StringP("unit", "u", "", "")
	cmd.Flags().StringP("generator", "g", consts.DefaultGenerator, "")
	cmd.Flags().StringP("value", "v", consts.DefaultValue, "")
	cmd.Flags().StringP("type", "t", consts.DefaultValueType, "")

	require.NoError(t, cmd.ParseFlags(args))
	return cmd
}

func TestBuildMetricTask(t *testing.T) {
	t.Run("attributes", func(t *testing.T) {
		cmd := newMetricTestCmd(t, "-a", "env=prod,method=GET", "--attributes", "region=eu")

		mc, err := buildMetricTask(cmd, consts.MetricTypeCounter)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"env": "prod", "method": "GET", "region": "eu"}, mc.Attributes)
	})

	t.Run("defaults", func(t *testing.T) {
		cmd := newMetricTestCmd(t)

		mc, err := buildMetricTask(cmd, consts.MetricTypeGauge)
		require.NoError(t, err)
		assert.Equal(t, consts.MetricTypeGauge, mc.Kind)
		assert.Equal(t, consts.ValueTypeFloat64, mc.Type)
		assert.Equal(t, time.Second, mc.Rate)
		assert.Empty(t, mc.Description)
		assert.Empty(t, mc.Attributes)
	})

	t.Run("negative rate", func(t *testing.T) {
		cmd := newMetricTestCmd(t, "--rate", "-1s")

		_, err := buildMetricTask(cmd, consts.MetricTypeCounter)
		assert.ErrorContains(t, err, "rate must be positive")
	})
}

func TestNew(t *testing.T) {
	// Keep the user's global otelconf file out of the loaded configuration.
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	opts := &app.Options{Version: "test"}

	// load runs the flag handling of a kind subcommand, without starting the pipeline.
	load := func(t *testing.T, args ...string) (*config.Config, error) {
		t.Helper()

		kindCmd, flags, err := New(opts).Find(args)
		require.NoError(t, err)
		require.NoError(t, kindCmd.ParseFlags(flags))

		return loadConfig(kindCmd, *opts, kindCmd.Name())
	}

	t.Run("kind subcommand adds its task", func(t *testing.T) {
		for alias, kind := range map[string]string{
			"c":   consts.MetricTypeCounter,
			"g":   consts.MetricTypeGauge,
			"h":   consts.MetricTypeHistogram,
			"udc": consts.MetricTypeUpDownCounter,
		} {
			cfg, err := load(t, alias, "-n", "test.metric")
			require.NoError(t, err)
			require.Len(t, cfg.Metrics.Tasks, 1)
			assert.Equal(t, kind, cfg.Metrics.Tasks[0].Kind)
			assert.Equal(t, "test.metric", cfg.Metrics.Tasks[0].Name)
		}
	})

	t.Run("invalid flags", func(t *testing.T) {
		cfg, err := load(t, "counter", "--name", "a")
		assert.EqualError(t, err, "failed to build metric config: invalid metric name a, please see OTel guidelines")
		assert.Nil(t, cfg)
	})

	t.Run("deprecated histogram flags", func(t *testing.T) {
		cmd := New(opts)
		cmd.SetArgs([]string{"histogram", "--histogram-buckets", "1,2"})
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		err := cmd.Execute()
		assert.EqualError(t, err, "metric views configuration via CLI is deprecated, please use --config with otelconf configuration")
	})
}
