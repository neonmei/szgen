package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/neonmei/szgen/internal/consts"
)

var (
	metricNameRegex = regexp.MustCompile(`^[a-zA-Z]{1}[\w\.]+$`)

	validMetricTypes   = []string{consts.MetricTypeCounter, consts.MetricTypeGauge, consts.MetricTypeHistogram, consts.MetricTypeUpDownCounter}
	ValidTemporalities = []string{consts.TemporalityCumulative, consts.TemporalityDelta}
	validModes         = []string{consts.ExportModeExecute, consts.ExportModeExecuteAndSave, consts.ExportModeSave}
	validGenerators    = []string{
		consts.GeneratorConstant,
		consts.GeneratorRandom,
		consts.GeneratorStep,
		consts.GeneratorSine,
		consts.GeneratorSequence,
	}
	validValueTypes      = []string{consts.ValueTypeInt64, consts.ValueTypeFloat64}
	validInstrumentKinds = []string{
		consts.InstrumentKindUndefined,
		consts.InstrumentKindCounter,
		consts.InstrumentKindGauge,
		consts.InstrumentKindHistogram,
		consts.InstrumentKindObservableCounter,
		consts.InstrumentKindObservableGauge,
		consts.InstrumentKindObservableUpDown,
		consts.InstrumentKindUpDownCounter,
	}
)

func ValidateMetricKind(metricKind string) error {
	if !slices.Contains(validMetricTypes, metricKind) {
		return fmt.Errorf("invalid type '%s', must be one of: %s", metricKind, strings.Join(validMetricTypes, ", "))
	}

	return nil
}

func ValidateInstrumentKind(instrumentKind string) error {
	if !slices.Contains(validInstrumentKinds, instrumentKind) {
		return fmt.Errorf("invalid instrument kind '%s', must be one of: %s", instrumentKind, strings.Join(validInstrumentKinds, ", "))
	}

	return nil
}

func ValidateGenerator(name string) error {
	if name == "" {
		return nil
	}

	if !slices.Contains(validGenerators, name) {
		return fmt.Errorf("invalid generator '%s', must be one of: %s", name, strings.Join(validGenerators, ", "))
	}

	return nil
}

func ValidateMode(mode string) error {
	if !slices.Contains(validModes, mode) {
		return fmt.Errorf("invalid mode '%s', must be one of: %s", mode, strings.Join(validModes, ", "))
	}

	return nil
}

func ValidateTemporality(temporality string) error {
	if !slices.Contains(ValidTemporalities, temporality) {
		return fmt.Errorf("export: invalid temporality '%s', must be one of: %s", temporality, strings.Join(ValidTemporalities, ", "))
	}

	return nil
}

func ValidateValueType(valueType string) error {
	if !slices.Contains(validValueTypes, valueType) {
		return fmt.Errorf("invalid value type '%s', must be one of: %s", valueType, strings.Join(validValueTypes, ", "))
	}

	return nil
}

// ValidateMetricName validates with OpenTelemetry naming conventions.
// This means name must start with a letter, and can contain letters, numbers, underscores, and dots.
//
// REF: https://opentelemetry.io/docs/specs/semconv/general/naming/
func ValidateMetricName(metricName string) error {
	if l := len(metricName); l == 0 || l > 255 {
		return fmt.Errorf("invalid metric length %d, should be between 1-255 characters", l)
	}

	if !metricNameRegex.MatchString(metricName) {
		return fmt.Errorf("invalid metric name %s, please see OTel guidelines", metricName)
	}
	return nil
}
