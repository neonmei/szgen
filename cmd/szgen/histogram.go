package main

import (
	"fmt"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
)

const (
	histogramBuckets = "histogram-buckets"
	expoMaxSize      = "expo-max-size"
	expoMaxScale     = "expo-max-scale"
	expoNoMinMax     = "expo-no-minmax"
)

var histogramCmd = &cobra.Command{
	Use:     "histogram",
	Aliases: []string{"h"},
	Short:   "Generate histogram metrics",
	Long:    `Generate histogram metrics with configurable buckets or exponential histogram parameters.`,
	RunE:    runHistogram,
}

func init() {
	metricsCmd.AddCommand(histogramCmd)
	histogramCmd.Flags().Float64Slice(histogramBuckets, nil, "Enable explicit buckets histogram through a view")
	histogramCmd.Flags().Int(expoMaxScale, 0, "Exponential histogram max scale")
	histogramCmd.Flags().Int(expoMaxSize, 0, "Exponential histogram max size")
	histogramCmd.Flags().Bool(expoNoMinMax, false, "Exponential histogram should include min and max")
}

func runHistogram(cmd *cobra.Command, _ []string) error {
	explicitBuckets, _ := cmd.Flags().GetFloat64Slice(histogramBuckets)
	expoMaxScale, _ := cmd.Flags().GetInt(expoMaxScale)
	expoMaxSize, _ := cmd.Flags().GetInt(expoMaxSize)

	if len(explicitBuckets) > 0 || expoMaxScale > 0 || expoMaxSize > 0 {
		return fmt.Errorf("metric views configuration via CLI is deprecated, please use --config with otelconf configuration")
	}

	return runMetricCommand(cmd, consts.MetricTypeHistogram)
}
