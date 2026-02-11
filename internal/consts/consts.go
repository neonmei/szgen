package consts

import (
	"math"
	"time"
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
	InstrumentKindUndefined     = ""
	InstrumentKindCounter       = "counter"
	InstrumentKindGauge         = "gauge"
	InstrumentKindHistogram     = "histogram"
	InstrumentKindUpDownCounter = "updowncounter"
)

const (
	DefaultConfigFile        = "szgen.yaml"
	DefaultCount             = 1
	DefaultDelta             = 1.0
	DefaultDescription       = "Metric generated with szgen"
	DefaultExecutorStrategy  = ExecutorStrategySerial
	DefaultExportTemporality = TemporalityDelta
	DefaultGenerator         = GeneratorConstant
	DefaultMeterName         = "szgen"
	DefaultMetricKind        = MetricTypeCounter
	DefaultMetricName        = "szgen.metric"
	DefaultOTLPEndpoint      = "http://127.0.0.1:4317"
	DefaultOTLPInsecure      = true
	DefaultOTLPInterval      = time.Second
	DefaultRate              = time.Second
	DefaultSineGeneratorB    = 10
	DefaultValue             = "1"
	DefaultValueType         = ValueTypeFloat64

	DefaultFlushTimeout = 5 * time.Second

	DefaultOTelIntervalMillis = 1000
	DefaultOTelTimeoutMillis  = 1000
	DefaultOTelMaxSize        = 100
	DefaultOTelMaxScale       = 10

	DefaultFilePerm = 0o644

	ParamMaxConcurrency = "max_concurrency"

	SineParamIndexB      = 1
	SineParamIndexVShift = 2
	SineParamIndexHShift = 3
	SineFullCircle       = 2 * math.Pi
)
