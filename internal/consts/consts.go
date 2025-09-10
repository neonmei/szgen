package consts

import "time"

const (
	ExportModeExecute        = "execute"
	ExportModeSave           = "save"
	ExportModeExecuteAndSave = "execute-and-save"
)

const (
	AggregationExplicitBucketHistogram = "explicit_bucket_histogram"
	AggregationExponentialHistogram    = "base2_exponential_histogram"
	ExecutorStrategySerial             = "serial"
	ExecutorStrategyConcurrent         = "concurrent"
	GeneratorConstant                  = "constant"
	GeneratorRandom                    = "random"
	GeneratorSequence                  = "sequence"
	GeneratorSine                      = "sine"
	GeneratorStep                      = "step"
	MetricTypeCounter                  = "counter"
	MetricTypeGauge                    = "gauge"
	MetricTypeHistogram                = "histogram"
	MetricTypeUpDownCounter            = "updowncounter"
	TemporalityCumulative              = "cumulative"
	TemporalityDelta                   = "delta"
	ValueTypeFloat64                   = "float64"
	ValueTypeInt64                     = "int64"
)

const (
	InstrumentKindUndefined         = ""
	InstrumentKindCounter           = "counter"
	InstrumentKindGauge             = "gauge"
	InstrumentKindHistogram         = "histogram"
	InstrumentKindObservableCounter = "observablecounter"
	InstrumentKindObservableGauge   = "observablegauge"
	InstrumentKindObservableUpDown  = "observableupdowncounter"
	InstrumentKindUpDownCounter     = "updowncounter"
)

const (
	DefaultConfigFile        = "szgen.yaml"
	DefaultCount             = 1
	DefaultDelta             = 1.0
	DefaultDescription       = "Metric generated with szgen"
	DefaultExecutorStrategy  = ExecutorStrategySerial
	DefaultExportMode        = ExportModeExecute
	DefaultExportTemporality = TemporalityDelta
	DefaultGenerator         = GeneratorConstant
	DefaultMetricKind        = MetricTypeCounter
	DefaultMetricName        = "szgen.metric"
	DefaultOTLPEndpoint      = "127.0.0.1:4317"
	DefaultOTLPInsecure      = true
	DefaultOTLPInterval      = time.Second
	DefaultRate              = time.Second
	DefaultServiceName       = "szgen"
	DefaultServiceVersion    = "0.1.0"
	DefaultSineGeneratorB    = 10
	DefaultValue             = "1"
	DefaultValueType         = ValueTypeFloat64
)
