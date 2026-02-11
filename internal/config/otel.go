package config

import "github.com/neonmei/szgen/internal/consts"

func NewOTelConfig(serviceVersion string) map[string]any {
	return map[string]any{
		"file_format": "1.0",
		"disabled":    false,
		"log_level":   "info",
		"meter_provider": map[string]any{
			"readers": []map[string]any{
				{
					"periodic": map[string]any{
						"interval": consts.DefaultOTelIntervalMillis,
						"timeout":  consts.DefaultOTelTimeoutMillis,
						"exporter": map[string]any{
							"otlp_grpc": map[string]any{
								"endpoint":                      consts.DefaultOTLPEndpoint,
								"encoding":                      "protobuf",
								"compression":                   "gzip",
								"insecure":                      consts.DefaultOTLPInsecure,
								"timeout":                       consts.DefaultOTelTimeoutMillis,
								"temporality_preference":        consts.DefaultExportTemporality,
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
								"max_size":   consts.DefaultOTelMaxSize,
								"max_scale":  consts.DefaultOTelMaxScale,
								"no_min_max": false,
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
					"value": serviceVersion,
					"type":  "string",
				},
			},
		},
	}
}
