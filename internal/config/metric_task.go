package config

import (
	"fmt"
	"time"

	"github.com/neonmei/szgen/internal/consts"
	"gopkg.in/yaml.v3"
)

type (
	MetricTaskOption func(*MetricTask)
	MetricTask       struct {
		Name        string            `yaml:"name"`
		Kind        string            `yaml:"kind"`
		Type        string            `yaml:"type,omitempty"`
		Rate        time.Duration     `yaml:"rate,omitempty"`
		Count       int               `yaml:"count,omitempty"`
		Value       string            `yaml:"value,omitempty"`
		Attributes  map[string]string `yaml:"attributes,omitempty"`
		Generator   string            `yaml:"generator,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Unit        string            `yaml:"unit,omitempty"`
	}
)

func NewMetricTask(options ...MetricTaskOption) *MetricTask {
	mt := &MetricTask{
		Name:      consts.DefaultMetricName,
		Kind:      consts.DefaultMetricKind,
		Type:      consts.DefaultValueType,
		Rate:      consts.DefaultRate,
		Count:     consts.DefaultCount,
		Value:     consts.DefaultValue,
		Generator: consts.DefaultGenerator,
	}

	for _, option := range options {
		option(mt)
	}

	return mt
}

func (mc *MetricTask) Validate() error {
	if err := ValidateMetricName(mc.Name); err != nil {
		return err
	}

	if err := ValidateMetricKind(mc.Kind); err != nil {
		return err
	}

	if err := ValidateGenerator(mc.Generator); err != nil {
		return err
	}

	if err := ValidateValueType(mc.Type); err != nil {
		return err
	}

	if mc.Rate == 0 {
		return fmt.Errorf("empty rate")
	}

	return nil
}

func (mc *MetricTask) UnmarshalYAML(node *yaml.Node) error {
	// Start with defaults
	defaultTask := NewMetricTask()

	// Unmarshal YAML over the defaults
	// Use type alias to avoid infinite recursion
	type rawMetricTask MetricTask
	if err := node.Decode((*rawMetricTask)(defaultTask)); err != nil {
		return err
	}

	// Copy result to receiver
	*mc = *defaultTask

	return nil
}

func WithName(name string) MetricTaskOption {
	return func(mt *MetricTask) {
		if name != "" {
			mt.Name = name
		}
	}
}

func WithKind(kind string) MetricTaskOption {
	return func(mt *MetricTask) {
		if kind != "" {
			mt.Kind = kind
		}
	}
}

func WithType(valueType string) MetricTaskOption {
	return func(mt *MetricTask) {
		if valueType != "" {
			mt.Type = valueType
		}
	}
}

func WithRate(rate time.Duration) MetricTaskOption {
	return func(mt *MetricTask) {
		if rate != 0 {
			mt.Rate = rate
		}
	}
}

func WithCount(count int) MetricTaskOption {
	return func(mt *MetricTask) {
		if count != 0 {
			mt.Count = count
		}
	}
}

func WithValue(value string) MetricTaskOption {
	return func(mt *MetricTask) {
		if value != "" {
			mt.Value = value
		}
	}
}

func WithMetricAttributes(attrs map[string]string) MetricTaskOption {
	return func(mt *MetricTask) {
		if attrs != nil {
			mt.Attributes = attrs
		}
	}
}

func WithGenerator(generator string) MetricTaskOption {
	return func(mt *MetricTask) {
		if generator != "" {
			mt.Generator = generator
		}
	}
}

func WithDescription(description string) MetricTaskOption {
	return func(mt *MetricTask) {
		if description != "" {
			mt.Description = description
		}
	}
}

func WithUnit(unit string) MetricTaskOption {
	return func(mt *MetricTask) {
		if unit != "" {
			mt.Unit = unit
		}
	}
}
