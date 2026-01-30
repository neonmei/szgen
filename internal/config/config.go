package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"go.opentelemetry.io/contrib/otelconf"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Metrics       *MetricsConfig `yaml:"metrics"`
	OpenTelemetry map[string]any `yaml:"opentelemetry"`
	Executor      ExecutorConfig `yaml:"executor,omitempty"`
}

type Option func(*Config) error

// NewConfig creates a new Config with default values.
// It applies options in order. If an option fails, it returns the error immediately.
func NewConfig(options ...Option) (*Config, error) {
	c := &Config{}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithDefaultConfig() Option {
	return func(c *Config) error {
		c.Metrics = &MetricsConfig{Tasks: []MetricTask{}}
		c.OpenTelemetry = NewOTelConfig()
		c.Executor = NewExecutorConfig()
		return nil
	}
}

func WithMetricsConfig(metrics *MetricsConfig) Option {
	return func(c *Config) error {
		c.Metrics = metrics
		return nil
	}
}

func WithOpenTelemetryConfig(otel map[string]any) Option {
	return func(c *Config) error {
		c.OpenTelemetry = otel
		return nil
	}
}

func WithExecutorConfig(executor ExecutorConfig) Option {
	return func(c *Config) error {
		c.Executor = executor
		return nil
	}
}

// WithOtelConfigFile does a best effort load of the global otelconf from ~/.config/szgen/opentelemetry.yaml.
// this file contains OpenTelemetry SDK configuration.
func WithOtelConfigFile() Option {
	return func(c *Config) error {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil
			}
			xdgConfigHome = filepath.Join(home, ".config")
		}
		userOtelConfig := filepath.Join(xdgConfigHome, "szgen", "opentelemetry.yaml")

		if _, err := os.Stat(userOtelConfig); err == nil {
			slog.Debug("Loading user config", "path", userOtelConfig)

			data, err := os.ReadFile(userOtelConfig)
			if err != nil {
				return fmt.Errorf("failed to read otelconf file %s: %w", userOtelConfig, err)
			}

			// We unmarshal into c.OpenTelemetry, effectively merging/overwriting
			if err := yaml.Unmarshal(data, c.OpenTelemetry); err != nil {
				return fmt.Errorf("failed to unmarshal otelconf file %s: %w", userOtelConfig, err)
			}
		}
		return nil
	}
}

// WithSzgenConfigFile loads the main szgen configuration file with pre-recorded tasks in a yaml file.
func WithSzgenConfigFile(path string) Option {
	return func(c *Config) error {
		if path == "" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read szgen config file %s: %w", path, err)
		}

		if err := yaml.Unmarshal(data, c); err != nil {
			return fmt.Errorf("failed to unmarshal szgen config file %s: %w", path, err)
		}
		return nil
	}
}

func (c *Config) Validate() error {
	if len(c.Metrics.Tasks) == 0 {
		return fmt.Errorf("no metrics defined in configuration")
	}

	if err := c.Executor.Validate(); err != nil {
		return err
	}

	for i, metric := range c.Metrics.Tasks {
		if err := metric.Validate(); err != nil {
			return fmt.Errorf("metric[%d]: %w", i, err)
		}
	}

	if len(c.OpenTelemetry) == 0 {
		return fmt.Errorf("no opentelemetry configuration found")
	}

	data, err := yaml.Marshal(c.OpenTelemetry)
	if err != nil {
		return fmt.Errorf("failed to marshal otelconf configuration: %w", err)
	}

	_, err = otelconf.ParseYAML(data)
	if err != nil {
		return fmt.Errorf("failed to parse otelconf schema: %w", err)
	}

	return nil
}
