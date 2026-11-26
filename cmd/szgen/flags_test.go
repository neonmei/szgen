package main

import (
	"testing"
	"time"

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
	cmd.Flags().StringP("executor", "e", consts.DefaultExecutorStrategy, "")
	cmd.Flags().IntP("max-concurrency", "j", 0, "")

	require.NoError(t, cmd.ParseFlags(args))
	return cmd
}

func TestBuildMetricConfig(t *testing.T) {
	t.Run("attributes", func(t *testing.T) {
		cmd := newMetricTestCmd(t, "-a", "env=prod,method=GET", "--attributes", "region=eu")

		mc, err := buildMetricConfig(cmd, consts.MetricTypeCounter)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"env": "prod", "method": "GET", "region": "eu"}, mc.Attributes)
	})

	t.Run("defaults", func(t *testing.T) {
		cmd := newMetricTestCmd(t)

		mc, err := buildMetricConfig(cmd, consts.MetricTypeGauge)
		require.NoError(t, err)
		assert.Equal(t, consts.MetricTypeGauge, mc.Kind)
		assert.Equal(t, consts.ValueTypeFloat64, mc.Type)
		assert.Equal(t, time.Second, mc.Rate)
		assert.Empty(t, mc.Description)
		assert.Empty(t, mc.Attributes)
	})

	t.Run("negative rate", func(t *testing.T) {
		cmd := newMetricTestCmd(t, "--rate", "-1s")

		_, err := buildMetricConfig(cmd, consts.MetricTypeCounter)
		assert.ErrorContains(t, err, "rate must be positive")
	})
}

func TestApplyExecutorFlags(t *testing.T) {
	fileConfig := func() config.ExecutorConfig {
		return config.NewExecutorConfig(
			config.WithExecutorStrategy(consts.ExecutorStrategyConcurrent),
			config.WithExecutorParams(map[string]any{consts.ParamMaxConcurrency: 4}),
		)
	}

	t.Run("unset flags keep config", func(t *testing.T) {
		ec := fileConfig()
		applyExecutorFlags(newMetricTestCmd(t), &ec)

		assert.Equal(t, fileConfig(), ec)
	})

	t.Run("set flags override config", func(t *testing.T) {
		ec := fileConfig()
		applyExecutorFlags(newMetricTestCmd(t, "-e", "serial", "-j", "0"), &ec)

		assert.Equal(t, consts.ExecutorStrategySerial, ec.Strategy)
		assert.Equal(t, 0, ec.Params[consts.ParamMaxConcurrency])
	})

	t.Run("nil params", func(t *testing.T) {
		ec := config.ExecutorConfig{Strategy: consts.ExecutorStrategySerial}
		applyExecutorFlags(newMetricTestCmd(t, "-j", "2"), &ec)

		assert.Equal(t, 2, ec.Params[consts.ParamMaxConcurrency])
	})
}
