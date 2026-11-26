// Package app is the szgen pipeline, independent of the command line: it loads the configuration,
// builds tasks and runs them with the OpenTelemetry SDK.
package app

import (
	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
)

// Options are settings that apply to every command.
type Options struct {
	// Version is reported as the service.version resource attribute.
	Version  string
	Executor ExecutorOverrides
}

// ExecutorOverrides replace executor settings of the configuration. A nil field keeps the configured value.
type ExecutorOverrides struct {
	Strategy       *string
	MaxConcurrency *int
}

func (o ExecutorOverrides) apply(ec *config.ExecutorConfig) {
	if o.Strategy != nil {
		ec.Strategy = *o.Strategy
	}

	if o.MaxConcurrency != nil {
		if ec.Params == nil {
			ec.Params = make(map[string]any)
		}
		ec.Params[consts.ParamMaxConcurrency] = *o.MaxConcurrency
	}
}
