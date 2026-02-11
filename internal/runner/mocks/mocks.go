package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/neonmei/szgen/internal/runner"
)

// MockTask is a manual implementation of runner.Task for testing.
type MockTask struct {
	NameFunc    func() string
	ExecuteFunc func(context.Context) error

	// Pre-configured behavior
	NameVal     string
	ExecuteTime time.Duration
	ExecuteErr  error

	// Call tracking
	ExecuteCalled int
	Mu            sync.Mutex
}

func (m *MockTask) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return m.NameVal
}

func (m *MockTask) Execute(ctx context.Context) error {
	m.Mu.Lock()
	m.ExecuteCalled++
	m.Mu.Unlock()

	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx)
	}

	if m.ExecuteTime > 0 {
		select {
		case <-time.After(m.ExecuteTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return m.ExecuteErr
}

// MockExecutor is a manual implementation of runner.Executor for testing.
type MockExecutor struct {
	ExecuteFunc func(context.Context, []runner.Task) error

	// Call tracking
	ExecuteCalled int
	ExecutedTasks []runner.Task
	Mu            sync.Mutex
}

func (m *MockExecutor) Execute(ctx context.Context, tasks []runner.Task) error {
	m.Mu.Lock()
	m.ExecuteCalled++
	m.ExecutedTasks = tasks
	m.Mu.Unlock()

	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, tasks)
	}

	return nil
}
