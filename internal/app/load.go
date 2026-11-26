package app

import (
	"fmt"

	"github.com/neonmei/szgen/internal/config"
)

// Load builds the configuration from the defaults, the global otelconf file and the given layers,
// then applies the executor overrides.
func Load(opts Options, layers ...config.Option) (*config.Config, error) {
	options := append([]config.Option{config.WithDefaultConfig(opts.Version), config.WithOtelConfigFile()}, layers...)

	cfg, err := config.NewConfig(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	opts.Executor.apply(&cfg.Executor)
	return cfg, nil
}
