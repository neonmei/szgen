package metrictask

import (
	"go.opentelemetry.io/otel/attribute"
)

func parseAttributes(attrs map[string]string) []attribute.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	attributes := make([]attribute.KeyValue, 0, len(attrs))
	for key, value := range attrs {
		attributes = append(attributes, attribute.String(key, value))
	}

	return attributes
}
