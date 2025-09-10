package main

import (
	"fmt"
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
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
	views, err := parseHistogramViews(cmd)
	if err != nil {
		return err
	}

	return runMetricCommandWithViews(cmd, consts.MetricTypeHistogram, views)
}

func parseHistogramViews(cmd *cobra.Command) ([]config.MetricView, error) {
	explicitBuckets, _ := cmd.Flags().GetFloat64Slice(histogramBuckets)
	expoMaxScale, _ := cmd.Flags().GetInt(expoMaxScale)
	expoMaxSize, _ := cmd.Flags().GetInt(expoMaxSize)
	expoNoMinMax, _ := cmd.Flags().GetBool(expoNoMinMax)

	hasExplicitBuckets := len(explicitBuckets) > 0
	hasExpoParams := expoMaxScale > 0 || expoMaxSize > 0

	if !hasExplicitBuckets && !hasExpoParams {
		return nil, nil
	}

	if hasExplicitBuckets && hasExpoParams {
		return nil, fmt.Errorf("for oneshot commands you can only configure either explicit buckets OR exponential histogram")
	}

	if hasExplicitBuckets {
		bucketView := config.NewMetricView(
			config.WithInstrumentKind(consts.InstrumentKindHistogram),
			config.WithExplicitBuckets(explicitBuckets),
		)

		if err := bucketView.Validate(); err != nil {
			return nil, fmt.Errorf("invalid explicit bucket histogram view: %w", err)
		}

		slog.Info("Configured explicit buckets view from CLI", "boundaries", explicitBuckets)
		return []config.MetricView{*bucketView}, nil
	}

	// Set defaults for exponential histogram if any parameter is missing
	if expoMaxScale == 0 {
		expoMaxScale = 20
	}
	if expoMaxSize == 0 {
		expoMaxSize = 160
	}

	expoView := config.NewMetricView(
		config.WithInstrumentKind(consts.InstrumentKindHistogram),
		config.WithExponentialHistogram(expoMaxScale, expoMaxSize, expoNoMinMax),
	)

	if err := expoView.Validate(); err != nil {
		return nil, fmt.Errorf("invalid exponential histogram view: %w", err)
	}

	slog.Info("Configured exponential histogram view from CLI",
		"maxSize", expoView.Stream.Aggregation.MaxSize,
		"maxScale", expoView.Stream.Aggregation.MaxScale,
		"noMinMax", expoView.Stream.Aggregation.NoMinMax,
	)

	return []config.MetricView{*expoView}, nil
}
