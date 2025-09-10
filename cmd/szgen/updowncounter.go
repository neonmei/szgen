package main

import (
	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/cobra"
)

var updowncounterCmd = &cobra.Command{
	Use:     "updowncounter",
	Aliases: []string{"udc"},
	Short:   "Generate updowncounter metrics",
	Long:    `Generate updowncounter metrics that can increase or decrease with configurable patterns.`,
	RunE:    runUpDownCounter,
}

func init() {
	metricsCmd.AddCommand(updowncounterCmd)
}

func runUpDownCounter(cmd *cobra.Command, _ []string) error {
	return runMetricCommand(cmd, consts.MetricTypeUpDownCounter)
}
