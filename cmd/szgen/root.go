package main

import (
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "szgen",
	Short: "OpenTelemetry Data Generator CLI",
	Long: `szgen is a command-line tool for generating synthetic OpenTelemetry data
for testing, development, and demonstration purposes.`,
	SilenceUsage:  true,
	SilenceErrors: true, // main prints the returned error
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")

		logging.StartLogger(logLevel, logFormat)
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

// applyExecutorFlags overrides the executor configuration with flags explicitly set on the command line.
func applyExecutorFlags(cmd *cobra.Command, ec *config.ExecutorConfig) {
	if cmd.Flags().Changed("executor") {
		ec.Strategy, _ = cmd.Flags().GetString("executor")
	}

	if cmd.Flags().Changed("max-concurrency") {
		maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
		if ec.Params == nil {
			ec.Params = make(map[string]any)
		}
		ec.Params[consts.ParamMaxConcurrency] = maxConcurrency
	}
}
