package cli

import (
	"errors"

	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/config"
	"github.com/spf13/cobra"
)

func newRunCommand(opts *app.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Execute configuration file",
		Long:    `Execute a YAML configuration file with multiple generation tasks.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				return errors.New("failed to load config: no config file path provided")
			}

			cfg, err := app.Load(*opts, config.WithSzgenConfigFile(configPath))
			if err != nil {
				return err
			}

			return app.Execute(cfg)
		},
	}

	cmd.Flags().String("config", "", "Configuration file path")

	return cmd
}
