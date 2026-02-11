package config

import (
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNewMetricTask(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		mt := NewMetricTask()
		assert.Equal(t, consts.DefaultMetricName, mt.Name)
		assert.Equal(t, consts.DefaultMetricKind, mt.Kind)
		assert.Equal(t, consts.DefaultValueType, mt.Type)
		assert.Equal(t, consts.DefaultRate, mt.Rate)
		assert.Equal(t, consts.DefaultCount, mt.Count)
		assert.Equal(t, consts.DefaultValue, mt.Value)
		assert.Equal(t, consts.DefaultGenerator, mt.Generator)
	})

	t.Run("with options", func(t *testing.T) {
		mt := NewMetricTask(
			WithName("custom.metric"),
			WithKind(consts.MetricTypeGauge),
			WithType(consts.ValueTypeInt64),
			WithRate(5*time.Second),
			WithCount(100),
			WithValue("50"),
			WithGenerator(consts.GeneratorRandom),
			WithDescription("test metric"),
			WithUnit("ms"),
			WithMetricAttributes(map[string]any{"key": "val"}),
		)

		assert.Equal(t, "custom.metric", mt.Name)
		assert.Equal(t, consts.MetricTypeGauge, mt.Kind)
		assert.Equal(t, consts.ValueTypeInt64, mt.Type)
		assert.Equal(t, 5*time.Second, mt.Rate)
		assert.Equal(t, 100, mt.Count)
		assert.Equal(t, "50", mt.Value)
		assert.Equal(t, consts.GeneratorRandom, mt.Generator)
		assert.Equal(t, "test metric", mt.Description)
		assert.Equal(t, "ms", mt.Unit)
		assert.Equal(t, map[string]any{"key": "val"}, mt.Attributes)
	})

	t.Run("options are unconditional setters", func(t *testing.T) {
		mt := NewMetricTask(
			WithName(""),
			WithKind(""),
			WithType(""),
			WithRate(0),
			WithCount(0),
			WithValue(""),
			WithGenerator(""),
			WithDescription(""),
			WithUnit(""),
			WithMetricAttributes(nil),
		)

		// Options set zero values; defaults are only from NewMetricTask() without options
		assert.Equal(t, "", mt.Name)
		assert.Equal(t, "", mt.Kind)
	})
}

func TestMetricTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    MetricTask
		wantErr bool
	}{
		{
			name: "valid task",
			task: MetricTask{
				Name:      "valid.metric",
				Kind:      consts.MetricTypeCounter,
				Type:      consts.ValueTypeFloat64,
				Generator: consts.GeneratorConstant,
				Rate:      1 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "invalid name",
			task:    MetricTask{Name: ""},
			wantErr: true,
		},
		{
			name: "invalid kind",
			task: MetricTask{
				Name: "valid.metric",
				Kind: "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid generator",
			task: MetricTask{
				Name:      "valid.metric",
				Kind:      consts.MetricTypeCounter,
				Generator: "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid value type",
			task: MetricTask{
				Name:      "valid.metric",
				Kind:      consts.MetricTypeCounter,
				Generator: consts.GeneratorConstant,
				Type:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty rate",
			task: MetricTask{
				Name:      "valid.metric",
				Kind:      consts.MetricTypeCounter,
				Type:      consts.ValueTypeFloat64,
				Generator: consts.GeneratorConstant,
				Rate:      0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricTask_UnmarshalYAML(t *testing.T) {
	t.Run("partial config uses defaults", func(t *testing.T) {
		yamlData := `
name: partial.config
kind: gauge
`
		var mt MetricTask
		err := yaml.Unmarshal([]byte(yamlData), &mt)
		assert.NoError(t, err)

		assert.Equal(t, "partial.config", mt.Name)
		assert.Equal(t, "gauge", mt.Kind)

		// Defaults
		assert.Equal(t, consts.DefaultValueType, mt.Type)
		assert.Equal(t, consts.DefaultRate, mt.Rate)
	})

	t.Run("full config overrides defaults", func(t *testing.T) {
		yamlData := `
name: full.config
kind: histogram
type: int64
rate: 5s
count: 50
value: "10"
generator: step
description: "desc"
unit: "s"
attributes:
  env: prod
`
		var mt MetricTask
		err := yaml.Unmarshal([]byte(yamlData), &mt)
		assert.NoError(t, err)

		assert.Equal(t, "full.config", mt.Name)
		assert.Equal(t, "histogram", mt.Kind)
		assert.Equal(t, "int64", mt.Type)
		assert.Equal(t, 5*time.Second, mt.Rate)
		assert.Equal(t, 50, mt.Count)
		assert.Equal(t, "10", mt.Value)
		assert.Equal(t, "step", mt.Generator)
		assert.Equal(t, "desc", mt.Description)
		assert.Equal(t, "s", mt.Unit)
		assert.Equal(t, map[string]any{"env": "prod"}, mt.Attributes)
	})
}
