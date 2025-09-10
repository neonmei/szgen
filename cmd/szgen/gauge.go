package main

import (
	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
)

var gaugeCmd = &cobra.Command{
	Use:     "gauge",
	Aliases: []string{"g"},
	Short:   "Generate gauge metrics",
	Long:    `Generate gauge metrics that can go up or down with configurable patterns.`,
	RunE:    runGauge,
}

func init() {
	metricsCmd.AddCommand(gaugeCmd)
}

func runGauge(cmd *cobra.Command, _ []string) error {
	return runMetricCommand(cmd, consts.MetricTypeGauge)
}
