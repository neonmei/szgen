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
								"max_size":       consts.DefaultOTelMaxSize,
								"max_scale":      consts.DefaultOTelMaxScale,
								"record_min_max": true,
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

// MergeOTelConfig deep merges src into dst and returns dst. Nested maps are merged recursively,
// any other value from src (including lists) replaces the one in dst, and null values in src are ignored.
func MergeOTelConfig(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any, len(src))
	}

	for key, srcVal := range src {
		if srcVal == nil {
			continue
		}

		srcMap, srcIsMap := srcVal.(map[string]any)
		dstMap, dstIsMap := dst[key].(map[string]any)
		if srcIsMap && dstIsMap {
			dst[key] = MergeOTelConfig(dstMap, srcMap)
			continue
		}

		dst[key] = srcVal
	}

	return dst
}
