package metricscmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	histogramBuckets = "histogram-buckets"
	expoMaxSize      = "expo-max-size"
	expoMaxScale     = "expo-max-scale"
	expoNoMinMax     = "expo-no-minmax"
)

// withDeprecatedHistogramFlags registers the former view flags and rejects them before anything else runs.
func withDeprecatedHistogramFlags(cmd *cobra.Command) {
	cmd.Flags().Float64Slice(histogramBuckets, nil, "Enable explicit buckets histogram through a view")
	cmd.Flags().Int(expoMaxScale, 0, "Exponential histogram max scale")
	cmd.Flags().Int(expoMaxSize, 0, "Exponential histogram max size")
	cmd.Flags().Bool(expoNoMinMax, false, "Exponential histogram should include min and max")

	runE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		explicitBuckets, _ := cmd.Flags().GetFloat64Slice(histogramBuckets)
		expoMaxScale, _ := cmd.Flags().GetInt(expoMaxScale)
		expoMaxSize, _ := cmd.Flags().GetInt(expoMaxSize)

		if len(explicitBuckets) > 0 || expoMaxScale > 0 || expoMaxSize > 0 {
			return errors.New("metric views configuration via CLI is deprecated, please use --config with otelconf configuration")
		}

		return runE(cmd, args)
	}
}
