package otel

import (
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func NewResource(resourceConfig *config.ResourceConfig) *resource.Resource {
	if resourceConfig == nil {
		return resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(consts.DefaultServiceName),
			semconv.ServiceVersion(consts.DefaultServiceVersion),
		)
	}

	attrs := []attribute.KeyValue{}
	if resourceConfig.ServiceName != "" {
		attrs = append(attrs, semconv.ServiceName(resourceConfig.ServiceName))
	}

	if resourceConfig.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(resourceConfig.ServiceVersion))
	}

	for key, value := range resourceConfig.Attributes {
		attrs = append(attrs, attribute.String(key, value))
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	)
}
