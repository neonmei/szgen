package main

import (
	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
)

var counterCmd = &cobra.Command{
	Use:     "counter",
	Aliases: []string{"c"},
	Short:   "Generate counter metrics",
	Long:    `Generate monotonically increasing counter metrics with configurable patterns.`,
	RunE:    runCounter,
}

func init() {
	metricsCmd.AddCommand(counterCmd)
}

func runCounter(cmd *cobra.Command, _ []string) error {
	return runMetricCommand(cmd, consts.MetricTypeCounter)
}
