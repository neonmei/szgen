package main

import (
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/otel"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "szgen",
	Short: "OpenTelemetry Data Generator CLI",
	Long: `szgen is a command-line tool for generating synthetic OpenTelemetry data
for testing, development, and demonstration purposes.`,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")

		otel.StartLogger(logLevel, logFormat)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("executor", "e", consts.DefaultExecutorStrategy, "Executor strategy (serial, concurrent)")
	rootCmd.PersistentFlags().IntP("max-concurrency", "j", 0, "Maximum concurrency for concurrent executor (0 = unlimited), only applies to concurrent executor")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json)")
}

func parseExecutorConfigFromCli(cmd *cobra.Command) (*config.ExecutorConfig, error) {
	strategy, _ := cmd.Flags().GetString("executor")
	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")

	// Create a new config to ensure defaults
	ec := config.NewExecutorConfig()

	if strategy != "" {
		ec.Strategy = strategy
	}

	if maxConcurrency > 0 {
		if ec.Params == nil {
			ec.Params = make(map[string]any)
		}
		ec.Params["max_concurrency"] = maxConcurrency
	}

	if err := ec.Validate(); err != nil {
		return nil, err
	}

	return &ec, nil
}
