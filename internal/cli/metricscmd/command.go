// Package metricscmd implements the `szgen metrics` command tree.
package metricscmd

import (
	"fmt"

	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
)

var kinds = []struct {
	kind, alias, short, long string
}{
	{
		kind:  consts.MetricTypeCounter,
		alias: "c",
		short: "Generate counter metrics",
		long:  `Generate monotonically increasing counter metrics with configurable patterns.`,
	},
	{
		kind:  consts.MetricTypeGauge,
		alias: "g",
		short: "Generate gauge metrics",
		long:  `Generate gauge metrics that can go up or down with configurable patterns.`,
	},
	{
		kind:  consts.MetricTypeHistogram,
		alias: "h",
		short: "Generate histogram metrics",
		long:  `Generate histogram metrics with configurable buckets or exponential histogram parameters.`,
	},
	{
		kind:  consts.MetricTypeUpDownCounter,
		alias: "udc",
		short: "Generate updowncounter metrics",
		long:  `Generate updowncounter metrics that can increase or decrease with configurable patterns.`,
	},
}

// New returns the `metrics` command with one subcommand per metric kind.
// Each subcommand builds a task from its flags and runs it.
func New(opts *app.Options) *cobra.Command {
	metricsCmd := &cobra.Command{
		Use:     "metrics",
		Aliases: []string{"m"},
		Short:   "Generate OpenTelemetry metrics",
		Long:    `Generate various types of OpenTelemetry metrics including counters, gauges, histograms, and updowncounters.`,
	}

	metricsCmd.PersistentFlags().StringToStringP("attributes", "a", nil, "Comma-separated key=value pairs")
	metricsCmd.PersistentFlags().IntP("count", "c", consts.DefaultCount, "Number of data points to generate")
	metricsCmd.PersistentFlags().DurationP("rate", "r", consts.DefaultRate, "Time interval between each generated data point")
	metricsCmd.PersistentFlags().StringP("description", "d", "", "Metric description")
	metricsCmd.PersistentFlags().StringP("name", "n", consts.DefaultMetricName, "Metric name")
	metricsCmd.PersistentFlags().StringP("unit", "u", "", "Metric unit")
	metricsCmd.PersistentFlags().StringP("generator", "g", consts.DefaultGenerator, "Value generation pattern")
	metricsCmd.PersistentFlags().StringP("value", "v", consts.DefaultValue, "Static value or value range")
	metricsCmd.PersistentFlags().StringP("type", "t", consts.DefaultValueType, "Value type (int64, float64)")

	for _, k := range kinds {
		kindCmd := &cobra.Command{
			Use:     k.kind,
			Aliases: []string{k.alias},
			Short:   k.short,
			Long:    k.long,
			RunE: func(cmd *cobra.Command, _ []string) error {
				cfg, err := loadConfig(cmd, *opts, k.kind)
				if err != nil {
					return err
				}

				return app.Execute(cfg)
			},
		}

		if k.kind == consts.MetricTypeHistogram {
			withDeprecatedHistogramFlags(kindCmd)
		}

		metricsCmd.AddCommand(kindCmd)
	}

	return metricsCmd
}

// loadConfig builds a task of the given kind from the flags and adds it to the configuration.
func loadConfig(cmd *cobra.Command, opts app.Options, kind string) (*config.Config, error) {
	task, err := buildMetricTask(cmd, kind)
	if err != nil {
		return nil, fmt.Errorf("failed to build metric config: %w", err)
	}

	return app.Load(opts, config.WithMetricTasks(*task))
}
