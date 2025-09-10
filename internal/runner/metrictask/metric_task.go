package metrictask

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neonmei/szgen/internal/config"
	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/generator"
	"github.com/neonmei/szgen/internal/runner"
)

type valueRecorder[T int64 | float64] func(context.Context, T)

type metricTask[T int64 | float64] struct {
	recorder    valueRecorder[T]
	genIter     generator.ValueGenerator[T]
	genInterval time.Duration
	taskName    string
}

func (im *metricTask[T]) Name() string {
	return im.taskName
}

func (im *metricTask[T]) Execute(ctx context.Context) error {
	slog.Info("Iterator task running", "name", im.taskName, "interval", im.genInterval)

	ticker := time.NewTicker(im.genInterval)
	defer ticker.Stop()

	for value := range im.genIter {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			im.recorder(ctx, value)
			slog.Debug("Recorded data point",
				"metric", im.taskName,
				"value", value,
			)
		}
	}

	slog.Info("Completed execution", "metric", im.taskName)
	return nil
}

// New creates a runnable task from model (file, cli, etc) configuration.
// The context here allows cancelling generation at the producer (i.e: value generator) level.
func New(ctx context.Context, mTask config.MetricTask) (runner.Task, error) {
	if err := mTask.Validate(); err != nil {
		return nil, err
	}

	switch mTask.Type {
	case consts.ValueTypeInt64:
		return newInstrument[int64](ctx, mTask)
	case consts.ValueTypeFloat64:
		return newInstrument[float64](ctx, mTask)
	default:
		return nil, fmt.Errorf("unsupported metric type: %s", mTask.Type)
	}
}
