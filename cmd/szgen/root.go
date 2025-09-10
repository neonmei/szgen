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
	rootCmd.PersistentFlags().String("otlp-endpoint", consts.DefaultOTLPEndpoint, "OTLP gRPC endpoint URL")
	rootCmd.PersistentFlags().Bool("otlp-insecure", consts.DefaultOTLPInsecure, "Enable insecure gRPC connection")
	rootCmd.PersistentFlags().Duration("otlp-interval", consts.DefaultOTLPInterval, "Time between exports for periodic reader")
	rootCmd.PersistentFlags().String("ca-cert", "", "Trusted Certificate Authority file path")
	rootCmd.PersistentFlags().String("client-cert", "", "Client certificate file path")
	rootCmd.PersistentFlags().String("client-key", "", "Client private key file path")
	rootCmd.PersistentFlags().String("config", "", "Configuration file path")
	rootCmd.PersistentFlags().String("export-temporality", consts.DefaultExportTemporality, "Global metric temporality behavior (cumulative, delta)")
	rootCmd.PersistentFlags().String("export-file", "", "OTLP JSON export file path")
	rootCmd.PersistentFlags().String("export-mode", consts.DefaultExportMode, "Export mode (execute, save, execute-and-save)")
	rootCmd.PersistentFlags().String("service-name", consts.DefaultServiceName, "Service name for OpenTelemetry resource")
	rootCmd.PersistentFlags().String("service-version", consts.DefaultServiceVersion, "Service version for OpenTelemetry resource")
	rootCmd.PersistentFlags().StringToString("resource-attributes", map[string]string{}, "Additional resource attributes in key=value,key2=value2 format")
	rootCmd.PersistentFlags().StringP("executor", "e", consts.DefaultExecutorStrategy, "Executor strategy (serial, concurrent)")
	rootCmd.PersistentFlags().IntP("max-concurrency", "j", 0, "Maximum concurrency for concurrent executor (0 = unlimited), only applies to concurrent executor")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json)")
}

func parseExportConfigFromCli(cmd *cobra.Command) (*config.ExportConfig, error) {
	endpoint, _ := cmd.Flags().GetString("otlp-endpoint")
	insecure, _ := cmd.Flags().GetBool("otlp-insecure")
	interval, _ := cmd.Flags().GetDuration("otlp-interval")
	caCert, _ := cmd.Flags().GetString("ca-cert")
	clientCert, _ := cmd.Flags().GetString("client-cert")
	clientKey, _ := cmd.Flags().GetString("client-key")
	temporality, _ := cmd.Flags().GetString("export-temporality")
	exportFile, _ := cmd.Flags().GetString("export-file")
	exportMode, _ := cmd.Flags().GetString("export-mode")

	options := []config.ExportOption{
		config.WithEndpoint(endpoint),
		config.WithInsecure(insecure),
		config.WithInterval(interval),
		config.WithCACert(caCert),
		config.WithClientCert(clientCert),
		config.WithClientKey(clientKey),
		config.WithTemporality(temporality),
		config.WithExportFile(exportFile),
		config.WithExportMode(exportMode),
	}

	ec := config.NewExportConfig(options...)

	if err := ec.Validate(); err != nil {
		return nil, err
	}

	return &ec, nil
}

func parseExecutorConfigFromCli(cmd *cobra.Command) (*config.ExecutorConfig, error) {
	strategy, _ := cmd.Flags().GetString("executor")
	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")

	params := make(map[string]any)
	if maxConcurrency > 0 {
		params["max_concurrency"] = maxConcurrency
	}

	options := []config.ExecutorOption{
		config.WithExecutorStrategy(strategy),
		config.WithExecutorParams(params),
	}

	ec := config.NewExecutorConfig(options...)

	if err := ec.Validate(); err != nil {
		return nil, err
	}

	return &ec, nil
}

func parseResourceConfigFromCli(cmd *cobra.Command) (*config.ResourceConfig, error) {
	serviceName, _ := cmd.Flags().GetString("service-name")
	serviceVersion, _ := cmd.Flags().GetString("service-version")
	resourceAttributes, _ := cmd.Flags().GetStringToString("resource-attributes")

	options := []config.ResourceOption{
		config.WithServiceName(serviceName),
		config.WithServiceVersion(serviceVersion),
		config.WithAttributes(resourceAttributes),
	}

	rc := config.NewResourceConfig(options...)

	if err := rc.Validate(); err != nil {
		return nil, err
	}

	return &rc, nil
}
