package config

type MetricsConfig struct {
	Tasks []MetricTask `yaml:"tasks"`
	Views []MetricView `yaml:"views"`
}
