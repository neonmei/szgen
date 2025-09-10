package config

import (
	"fmt"

	"github.com/neonmei/szgen/internal/consts"
)

type ExecutorConfig struct {
	Strategy string         `yaml:"strategy,omitempty"`
	Params   map[string]any `yaml:"params,omitempty"`
}

type ExecutorOption func(*ExecutorConfig)

func WithExecutorStrategy(strategy string) ExecutorOption {
	return func(ec *ExecutorConfig) {
		if strategy != "" {
			ec.Strategy = strategy
		}
	}
}

func WithExecutorParams(params map[string]any) ExecutorOption {
	return func(ec *ExecutorConfig) {
		if params != nil {
			ec.Params = params
		}
	}
}

func NewExecutorConfig(options ...ExecutorOption) ExecutorConfig {
	ec := ExecutorConfig{
		Strategy: consts.DefaultExecutorStrategy,
		Params:   make(map[string]any),
	}

	for _, option := range options {
		option(&ec)
	}

	return ec
}

func (ec *ExecutorConfig) Validate() error {
	if err := ValidateExecutorStrategy(ec.Strategy); err != nil {
		return err
	}

	return nil
}

func ValidateExecutorStrategy(strategy string) error {
	switch strategy {
	case consts.ExecutorStrategySerial, consts.ExecutorStrategyConcurrent:
		return nil
	default:
		return fmt.Errorf("invalid executor strategy '%s', must be one of: %s, %s",
			strategy, consts.ExecutorStrategySerial, consts.ExecutorStrategyConcurrent)
	}
}
