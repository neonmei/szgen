package traces

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/neonmei/szgen/internal/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// newTestExecTask records spans through the global provider, like the SDK installed by app.Run.
// Tests using it must not run in parallel.
func newTestExecTask(t *testing.T, argv ...string) (runner.Task, *tracetest.SpanRecorder) {
	t.Helper()

	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(previous)
		_ = tp.Shutdown(context.Background())
	})

	task, err := NewExecTask(argv)
	require.NoError(t, err)

	return task, recorder
}

func executeTraced(t *testing.T, ctx context.Context, task runner.Task) error {
	t.Helper()

	err := task.Execute(ctx)
	if errors.Is(err, syscall.EPERM) {
		t.Skip("ptrace is not permitted in this environment")
	}
	return err
}

func TestExecTask(t *testing.T) {
	t.Run("process tree", func(t *testing.T) {
		task, recorder := newTestExecTask(t, "sh", "-c", `sleep 0; sh -c "exit 3"; exit 5`)

		err := executeTraced(t, context.Background(), task)
		var exitErr *ExitError
		require.ErrorAs(t, err, &exitErr)
		assert.Equal(t, 5, exitErr.ExitCode())

		spans := recorder.Ended()
		require.Len(t, spans, 3)

		byName := make(map[string][]sdktrace.ReadOnlySpan)
		for _, span := range spans {
			byName[span.Name()] = append(byName[span.Name()], span)
		}
		require.Len(t, byName["sleep"], 1)
		require.Len(t, byName["sh"], 2)

		var root, nested sdktrace.ReadOnlySpan
		for _, span := range byName["sh"] {
			if span.Parent().IsValid() {
				nested = span
			} else {
				root = span
			}
		}
		require.NotNil(t, root)
		require.NotNil(t, nested)

		assert.Equal(t, int64(5), attr(root, "process.exit.code").AsInt64())
		assert.Equal(t, int64(3), attr(nested, "process.exit.code").AsInt64())
		assert.Equal(t, int64(0), attr(byName["sleep"][0], "process.exit.code").AsInt64())

		for _, child := range []sdktrace.ReadOnlySpan{nested, byName["sleep"][0]} {
			assert.Equal(t, root.SpanContext().SpanID(), child.Parent().SpanID())
			assert.Equal(t, attr(root, "process.pid"), attr(child, "process.parent_pid"))
		}
	})

	t.Run("success", func(t *testing.T) {
		task, recorder := newTestExecTask(t, "true")

		require.NoError(t, executeTraced(t, context.Background(), task))
		require.Len(t, recorder.Ended(), 1)
		assert.Equal(t, codes.Unset, recorder.Ended()[0].Status().Code)
	})

	t.Run("multi-threaded process has a single span", func(t *testing.T) {
		// Any Go program starts several threads.
		if _, err := exec.LookPath("go"); err != nil {
			t.Skip("go is not in PATH")
		}
		task, recorder := newTestExecTask(t, "go", "version")

		require.NoError(t, executeTraced(t, context.Background(), task))
		for _, span := range recorder.Ended() {
			assert.Equal(t, "go", span.Name())
		}
		assert.Len(t, recorder.Ended(), 1)
	})

	t.Run("cancellation forwards a signal", func(t *testing.T) {
		task, recorder := newTestExecTask(t, "sleep", "30")

		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(200*time.Millisecond, cancel)

		start := time.Now()
		err := executeTraced(t, ctx, task)
		assert.Less(t, time.Since(start), 10*time.Second)

		var exitErr *ExitError
		require.ErrorAs(t, err, &exitErr)
		assert.Equal(t, 128+int(syscall.SIGTERM), exitErr.ExitCode())
		require.Len(t, recorder.Ended(), 1)
		assert.Equal(t, "terminated by signal SIGTERM", recorder.Ended()[0].Status().Description)
	})

	t.Run("untraced children do not block", func(t *testing.T) {
		// Like a shell that forks a helper and then execs szgen in the same process.
		unrelated := exec.Command("sleep", "30")
		require.NoError(t, unrelated.Start())
		t.Cleanup(func() {
			_ = unrelated.Process.Kill()
			_ = unrelated.Wait()
		})

		task, recorder := newTestExecTask(t, "true")

		done := make(chan error, 1)
		go func() { done <- task.Execute(context.Background()) }()

		select {
		case err := <-done:
			if errors.Is(err, syscall.EPERM) {
				t.Skip("ptrace is not permitted in this environment")
			}
			require.NoError(t, err)
			assert.Len(t, recorder.Ended(), 1)
		case <-time.After(10 * time.Second):
			t.Fatal("tracing waited for a child it does not trace")
		}
	})

	t.Run("unknown command", func(t *testing.T) {
		_, err := NewExecTask([]string{"szgen-no-such-command"})
		assert.ErrorContains(t, err, `exec: "szgen-no-such-command": executable file not found in $PATH`)

		_, err = NewExecTask(nil)
		assert.EqualError(t, err, "no command provided")
	})
}
