package metrictask

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

func parseAttributes(attrs map[string]any) []attribute.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	attributes := make([]attribute.KeyValue, 0, len(attrs))
	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			attributes = append(attributes, attribute.String(key, v))
		case int:
			attributes = append(attributes, attribute.Int(key, v))
		case float64:
			attributes = append(attributes, attribute.Float64(key, v))
		case bool:
			attributes = append(attributes, attribute.Bool(key, v))
		default:
			attributes = append(attributes, attribute.String(key, fmt.Sprint(v)))
		}
	}

	return attributes
}
