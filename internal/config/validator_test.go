package config

import (
	"strings"
	"testing"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestValidateMetricKind(t *testing.T) {
	tests := []struct {
		kind    string
		wantErr bool
	}{
		{consts.MetricTypeCounter, false},
		{consts.MetricTypeGauge, false},
		{consts.MetricTypeHistogram, false},
		{consts.MetricTypeUpDownCounter, false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			err := ValidateMetricKind(tt.kind)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateInstrumentKind(t *testing.T) {
	tests := []struct {
		kind    string
		wantErr bool
	}{
		{consts.InstrumentKindCounter, false},
		{consts.InstrumentKindGauge, false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			err := ValidateInstrumentKind(tt.kind)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGenerator(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{consts.GeneratorConstant, false},
		{consts.GeneratorRandom, false},
		{"", false}, // Allowed empty
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenerator(tt.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTemporality(t *testing.T) {
	tests := []struct {
		temporality string
		wantErr     bool
	}{
		{consts.TemporalityDelta, false},
		{consts.TemporalityCumulative, false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.temporality, func(t *testing.T) {
			err := ValidateTemporality(tt.temporality)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateValueType(t *testing.T) {
	tests := []struct {
		valType string
		wantErr bool
	}{
		{consts.ValueTypeInt64, false},
		{consts.ValueTypeFloat64, false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.valType, func(t *testing.T) {
			err := ValidateValueType(tt.valType)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMetricName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid.metric.name", false},
		{"valid_metric_name", false},
		{"ValidMetricName", false},
		{"", true},
		{strings.Repeat("a", 256), true},
		{"invalid metric name", true}, // spaces
		{"1invalid", true},            // starts with number
		{".invalid", true},            // starts with dot
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricName(tt.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
