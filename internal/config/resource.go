package config

import (
	"fmt"

	"github.com/neonmei/szgen/internal/consts"
)

type ResourceConfig struct {
	ServiceName    string            `yaml:"service_name,omitempty"`
	ServiceVersion string            `yaml:"service_version,omitempty"`
	Attributes     map[string]string `yaml:"attributes,omitempty"`
}

type ResourceOption func(*ResourceConfig)

func WithServiceName(name string) ResourceOption {
	return func(rc *ResourceConfig) {
		if name != "" {
			rc.ServiceName = name
		}
	}
}

func WithServiceVersion(version string) ResourceOption {
	return func(rc *ResourceConfig) {
		if version != "" {
			rc.ServiceVersion = version
		}
	}
}

func WithAttributes(attrs map[string]string) ResourceOption {
	return func(rc *ResourceConfig) {
		if len(attrs) > 0 {
			rc.Attributes = attrs
		}
	}
}

func NewResourceConfig(options ...ResourceOption) ResourceConfig {
	rc := ResourceConfig{
		ServiceName:    consts.DefaultServiceName,
		ServiceVersion: consts.DefaultServiceVersion,
		Attributes:     make(map[string]string),
	}

	for _, option := range options {
		option(&rc)
	}

	return rc
}

func (rc *ResourceConfig) Validate() error {
	if rc.ServiceName == "" {
		return fmt.Errorf("empty service name")
	}

	return nil
}
