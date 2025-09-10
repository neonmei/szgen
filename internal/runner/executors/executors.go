package executors

import (
	"fmt"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
)

func New(cfg config.ExecutorConfig) (runner.Executor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid executor configuration: %w", err)
	}

	switch cfg.Strategy {
	case consts.ExecutorStrategySerial:
		return NewSerial(), nil
	case consts.ExecutorStrategyConcurrent:
		return NewConcurrent(cfg.Params), nil
	default:
		return nil, fmt.Errorf("unknown executor strategy: %s", cfg.Strategy)
	}
}
