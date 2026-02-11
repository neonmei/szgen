package config

import (
	"testing"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutorConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := NewExecutorConfig()
		assert.Equal(t, consts.DefaultExecutorStrategy, cfg.Strategy)
		assert.Empty(t, cfg.Params)
	})

	t.Run("with options", func(t *testing.T) {
		params := map[string]any{"max_concurrency": 5}
		cfg := NewExecutorConfig(
			WithExecutorStrategy(consts.ExecutorStrategyConcurrent),
			WithExecutorParams(params),
		)
		assert.Equal(t, consts.ExecutorStrategyConcurrent, cfg.Strategy)
		assert.Equal(t, params, cfg.Params)
	})

	t.Run("empty options ignored", func(t *testing.T) {
		cfg := NewExecutorConfig(
			WithExecutorStrategy(""),
			WithExecutorParams(nil),
		)
		assert.Equal(t, consts.DefaultExecutorStrategy, cfg.Strategy)
		assert.Empty(t, cfg.Params)
	})
}

func TestExecutorConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ExecutorConfig
		wantErr bool
	}{
		{
			name:    "valid serial",
			cfg:     ExecutorConfig{Strategy: consts.ExecutorStrategySerial},
			wantErr: false,
		},
		{
			name:    "valid concurrent",
			cfg:     ExecutorConfig{Strategy: consts.ExecutorStrategyConcurrent},
			wantErr: false,
		},
		{
			name:    "invalid strategy",
			cfg:     ExecutorConfig{Strategy: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
