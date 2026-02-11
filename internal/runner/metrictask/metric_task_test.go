package metrictask

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricTask_Execute(t *testing.T) {
	t.Run("execute and record values", func(t *testing.T) {
		// Create a generator that yields a few values
		genFunc := func(yield func(int64) bool) {
			for i := 0; i < 5; i++ {
				if !yield(int64(i)) {
					return
				}
			}
		}

		var recordedValues []int64
		var mu sync.Mutex

		recorder := func(ctx context.Context, val int64) {
			mu.Lock()
			defer mu.Unlock()
			recordedValues = append(recordedValues, val)
		}

		task := &metricTask[int64]{
			taskName:    "test-task",
			genInterval: 1 * time.Millisecond,
			genIter:     generator.ValueGenerator[int64](genFunc),
			recorder:    recorder,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := task.Execute(ctx)
		require.NoError(t, err)

		mu.Lock()
		defer mu.Unlock()
		assert.Len(t, recordedValues, 5)
		assert.Equal(t, []int64{0, 1, 2, 3, 4}, recordedValues)
	})

	t.Run("stop on context cancellation", func(t *testing.T) {
		// Infinite generator
		genFunc := func(yield func(int64) bool) {
			i := 0
			for {
				if !yield(int64(i)) {
					return
				}
				i++
			}
		}

		recorder := func(_ context.Context, _ int64) {}

		task := &metricTask[int64]{
			taskName:    "infinite-task",
			genInterval: 10 * time.Millisecond,
			genIter:     generator.ValueGenerator[int64](genFunc),
			recorder:    recorder,
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Run in goroutine
		errChan := make(chan error)
		go func() {
			errChan <- task.Execute(ctx)
		}()

		// Let it run for a bit
		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-errChan:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(1 * time.Second):
			t.Fatal("task did not stop after context cancellation")
		}
	})
}
