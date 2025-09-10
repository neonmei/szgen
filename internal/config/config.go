package config

import (
	"fmt"
)

type Config struct {
	Metrics  *MetricsConfig `yaml:"metrics"`
	Export   ExportConfig   `yaml:"export"`
	Resource ResourceConfig `yaml:"resource,omitempty"`
	Executor ExecutorConfig `yaml:"executor,omitempty"`
}

type Option func(*Config)

func WithMetricsConfig(metrics *MetricsConfig) Option {
	return func(c *Config) {
		c.Metrics = metrics
	}
}

func WithExportConfig(export ExportConfig) Option {
	return func(c *Config) {
		c.Export = export
	}
}

func WithResourceConfig(resource ResourceConfig) Option {
	return func(c *Config) {
		c.Resource = resource
	}
}

func WithExecutorConfig(executor ExecutorConfig) Option {
	return func(c *Config) {
		c.Executor = executor
	}
}

func NewConfig(options ...Option) *Config {
	c := &Config{
		Metrics:  &MetricsConfig{Tasks: []MetricTask{}, Views: []MetricView{}},
		Export:   NewExportConfig(),
		Resource: NewResourceConfig(),
		Executor: NewExecutorConfig(),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Config) Validate() error {
	if len(c.Metrics.Tasks) == 0 {
		return fmt.Errorf("no metrics defined in configuration")
	}

	if err := c.Export.Validate(); err != nil {
		return err
	}

	if err := c.Resource.Validate(); err != nil {
		return err
	}

	if err := c.Executor.Validate(); err != nil {
		return err
	}

	for i, metric := range c.Metrics.Tasks {
		if err := metric.Validate(); err != nil {
			return fmt.Errorf("metric[%d]: %w", i, err)
		}
	}

	return nil
}
