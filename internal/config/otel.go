package config

func NewOTelConfig() map[string]any {
	return map[string]any{
		"file_format": "1.0-rc.2",
		"disabled":    false,
		"log_level":   "info",
		"meter_provider": map[string]any{
			"readers": []map[string]any{
				{
					"periodic": map[string]any{
						"interval": 1000,
						"timeout":  1000,
						"exporter": map[string]any{
							"otlp_grpc": map[string]any{
								"endpoint":                      "http://127.0.0.1:4317",
								"encoding":                      "protobuf",
								"compression":                   "gzip",
								"insecure":                      true,
								"timeout":                       1000,
								"temporality_preference":        "delta",
								"default_histogram_aggregation": "base2_exponential_bucket_histogram",
							},
						},
					},
				},
			},
			"views": []map[string]any{
				{
					"selector": map[string]any{
						"instrument_type": "histogram",
					},
					"stream": map[string]any{
						"aggregation": map[string]any{
							"base2_exponential_bucket_histogram": map[string]any{
								"max_size":  100,
								"max_scale": 10,
							},
						},
					},
				},
			},
		},
		"resource": map[string]any{
			"attributes": []map[string]any{
				{
					"name":  "service.name",
					"value": "szgen",
					"type":  "string",
				},
				{
					"name":  "service.version",
					"value": "0.1.0",
					"type":  "string",
				},
			},
		},
	}
}
