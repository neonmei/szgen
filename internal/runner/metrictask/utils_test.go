package metrictask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestParseAttributes(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected []attribute.KeyValue
	}{
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: nil,
		},
		{
			name:     "single string attribute",
			input:    map[string]any{"key": "value"},
			expected: []attribute.KeyValue{attribute.String("key", "value")},
		},
		{
			name:     "typed attributes",
			input:    map[string]any{"s": "text", "i": 42, "f": 3.14, "b": true},
			expected: []attribute.KeyValue{attribute.String("s", "text"), attribute.Int("i", 42), attribute.Float64("f", 3.14), attribute.Bool("b", true)},
		},
		{
			name:     "fallback to string",
			input:    map[string]any{"arr": []int{1, 2, 3}},
			expected: []attribute.KeyValue{attribute.String("arr", "[1 2 3]")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAttributes(tt.input)

			if len(tt.expected) == 0 {
				assert.Empty(t, got)
			} else {
				assert.Len(t, got, len(tt.expected))
				for _, expectedKV := range tt.expected {
					found := false
					for _, gotKV := range got {
						if gotKV.Key == expectedKV.Key && gotKV.Value == expectedKV.Value {
							found = true
							break
						}
					}
					assert.True(t, found, "expected key-value pair %v not found", expectedKV)
				}
			}
		})
	}
}
