package executors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/runner"
	"github.com/neonmei/szgen/internal/runner/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestConcurrentExecutor_Execute(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("execute tasks concurrently", func(t *testing.T) {
		// Create long running tasks
		taskCount := 10
		tasks := make([]runner.Task, taskCount)
		for i := 0; i < taskCount; i++ {
			tasks[i] = &mocks.MockTask{
				NameVal:     "task",
				ExecuteTime: 10 * time.Millisecond,
			}
		}

		// With concurrency 5, it should take approx 2 * 10ms = 20ms (+ overhead)
		// Sequential would take 100ms
		exec := NewConcurrent(map[string]any{"max_concurrency": 5})

		start := time.Now()
		err := exec.Execute(context.Background(), tasks)
		elapsed := time.Since(start)

		require.NoError(t, err)
		// Allow some overhead, but it should be faster than sequential
		assert.Less(t, elapsed, 80*time.Millisecond, "execution should be concurrent")
	})

	t.Run("abort on error", func(t *testing.T) {
		tasks := []runner.Task{
			&mocks.MockTask{NameVal: "task1", ExecuteTime: 50 * time.Millisecond},
			&mocks.MockTask{NameVal: "error-task", ExecuteErr: errors.New("failed")},
			&mocks.MockTask{NameVal: "task3", ExecuteTime: 50 * time.Millisecond},
		}

		exec := NewConcurrent(map[string]any{"max_concurrency": 2})
		err := exec.Execute(context.Background(), tasks)
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		tasks := []runner.Task{
			&mocks.MockTask{NameVal: "task1", ExecuteTime: 100 * time.Millisecond},
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel shortly after start
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		exec := NewConcurrent(map[string]any{})
		err := exec.Execute(ctx, tasks)

		// Should return nil (successful completion of *what was possible* or error?)
		// The current implementation returns nil if ALL tasks succeed OR if context is canceled
		// Let's check implementation behavior: returns nil if context canceled?
		// Wait, looking at implementation:
		// if err := task.Execute(ctx); err != nil && err != context.Canceled { return error }
		// So if task.Execute returns context.Canceled (which MockTask does), it returns nil.
		require.NoError(t, err)
	})
}
