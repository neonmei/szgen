package config

import (
	"time"

	"github.com/neonmei/szgen/internal/consts"
)

type ExportConfig struct {
	Mode        string        `yaml:"mode,omitempty"`
	Endpoint    string        `yaml:"endpoint,omitempty"`
	File        string        `yaml:"file,omitempty"`
	Temporality string        `yaml:"temporality,omitempty"`
	Insecure    bool          `yaml:"insecure,omitempty"`
	Interval    time.Duration `yaml:"interval,omitempty"`
	CACert      string        `yaml:"ca_cert,omitempty"`
	ClientCert  string        `yaml:"client_cert,omitempty"`
	ClientKey   string        `yaml:"client_key,omitempty"`
}

type ExportOption func(*ExportConfig)

func WithExportMode(mode string) ExportOption {
	return func(ec *ExportConfig) {
		if mode != "" {
			ec.Mode = mode
		}
	}
}

func WithEndpoint(endpoint string) ExportOption {
	return func(ec *ExportConfig) {
		if endpoint != "" {
			ec.Endpoint = endpoint
		}
	}
}

func WithExportFile(file string) ExportOption {
	return func(ec *ExportConfig) {
		if file != "" {
			ec.File = file
		}
	}
}

func WithTemporality(temporality string) ExportOption {
	return func(ec *ExportConfig) {
		if temporality != "" {
			ec.Temporality = temporality
		}
	}
}

func WithInsecure(insecure bool) ExportOption {
	return func(ec *ExportConfig) {
		ec.Insecure = insecure
	}
}

func WithInterval(interval time.Duration) ExportOption {
	return func(ec *ExportConfig) {
		if interval != 0 {
			ec.Interval = interval
		}
	}
}

func WithCACert(caCert string) ExportOption {
	return func(ec *ExportConfig) {
		if caCert != "" {
			ec.CACert = caCert
		}
	}
}

func WithClientCert(clientCert string) ExportOption {
	return func(ec *ExportConfig) {
		if clientCert != "" {
			ec.ClientCert = clientCert
		}
	}
}

func WithClientKey(clientKey string) ExportOption {
	return func(ec *ExportConfig) {
		if clientKey != "" {
			ec.ClientKey = clientKey
		}
	}
}

func NewExportConfig(options ...ExportOption) ExportConfig {
	ec := ExportConfig{
		Mode:        consts.DefaultExportMode,
		Endpoint:    consts.DefaultOTLPEndpoint,
		Temporality: consts.DefaultExportTemporality,
		Interval:    consts.DefaultOTLPInterval,
	}

	for _, option := range options {
		option(&ec)
	}

	return ec
}

func (ec *ExportConfig) Validate() error {
	if err := ValidateMode(ec.Mode); err != nil {
		return err
	}

	if err := ValidateTemporality(ec.Temporality); err != nil {
		return err
	}
	return nil
}
