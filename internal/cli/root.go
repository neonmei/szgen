// Package cli is the szgen command line: it parses flags and hands the pipeline over to package app.
package cli

import (
	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/cli/metricscmd"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewRootCommand returns the szgen command tree. Errors are returned, not printed.
func NewRootCommand(version string) *cobra.Command {
	opts := &app.Options{Version: version}

	rootCmd := &cobra.Command{
		Use:   "szgen",
		Short: "OpenTelemetry Data Generator CLI",
		Long: `szgen is a command-line tool for generating synthetic OpenTelemetry data
for testing, development, and demonstration purposes.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			logLevel, _ := cmd.Flags().GetString("log-level")
			logFormat, _ := cmd.Flags().GetString("log-format")
			telemetry.StartLogger(logLevel, logFormat)

			opts.Executor = executorOverrides(cmd.Flags())
		},
	}

	rootCmd.PersistentFlags().StringP("executor", "e", consts.DefaultExecutorStrategy, "Executor strategy (serial, concurrent)")
	rootCmd.PersistentFlags().IntP("max-concurrency", "j", 0, "Maximum concurrency for concurrent executor (0 = unlimited), only applies to concurrent executor")
	rootCmd.PersistentFlags().String("log-level", consts.DefaultLogLevel, "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json)")

	rootCmd.AddCommand(newRunCommand(opts), metricscmd.New(opts))
	rootCmd.AddCommand(platformCommands(opts)...)

	return rootCmd
}

// executorOverrides returns the executor flags explicitly set on the command line.
func executorOverrides(flags *pflag.FlagSet) app.ExecutorOverrides {
	var overrides app.ExecutorOverrides

	if flags.Changed("executor") {
		strategy, _ := flags.GetString("executor")
		overrides.Strategy = &strategy
	}

	if flags.Changed("max-concurrency") {
		maxConcurrency, _ := flags.GetInt("max-concurrency")
		overrides.MaxConcurrency = &maxConcurrency
	}

	return overrides
}
