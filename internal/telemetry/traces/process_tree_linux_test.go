package traces

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func newTestProcessTree(t *testing.T) (*processTree, *tracetest.SpanRecorder) {
	t.Helper()

	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	return newProcessTree(context.Background(), tp.Tracer("test")), recorder
}

func spansByName(recorder *tracetest.SpanRecorder) map[string]sdktrace.ReadOnlySpan {
	spans := make(map[string]sdktrace.ReadOnlySpan)
	for _, span := range recorder.Ended() {
		spans[span.Name()] = span
	}
	return spans
}

func attr(span sdktrace.ReadOnlySpan, key string) attribute.Value {
	for _, kv := range span.Attributes() {
		if string(kv.Key) == key {
			return kv.Value
		}
	}
	return attribute.Value{}
}

func TestProcessTree(t *testing.T) {
	shell := procInfo{path: "/bin/sh", comm: "sh", args: []string{"sh", "-c", "ls"}}
	ls := procInfo{path: "/usr/bin/coreutils", comm: "ls", args: []string{"/usr/bin/ls", "-l"}}
	now := time.Now()

	t.Run("fork, exec and exit", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		tr.startProcess(100, 1, now, shell)
		tr.startProcess(101, 100, now, procInfo{})
		assert.Equal(t, shell, tr.procs[101].info, "forked process inherits its parent's program")

		tr.exec(101, now, ls)
		assert.True(t, tr.exit(101, now, 0, ""))
		assert.True(t, tr.exit(100, now, 0, ""))

		spans := spansByName(recorder)
		require.Len(t, spans, 2)
		root, child := spans["sh"], spans["ls"]

		assert.False(t, root.Parent().IsValid())
		assert.Equal(t, root.SpanContext().SpanID(), child.Parent().SpanID())
		assert.Equal(t, int64(101), attr(child, "process.pid").AsInt64())
		assert.Equal(t, int64(100), attr(child, "process.parent_pid").AsInt64())
		assert.Equal(t, "/usr/bin/coreutils", attr(child, "process.executable.path").AsString())
		assert.Equal(t, "coreutils", attr(child, "process.executable.name").AsString())
		assert.Equal(t, "/usr/bin/ls", attr(child, "process.command").AsString())
		assert.Equal(t, []string{"/usr/bin/ls", "-l"}, attr(child, "process.command_args").AsStringSlice())
		assert.Equal(t, int64(0), attr(child, "process.exit.code").AsInt64())
		assert.Equal(t, codes.Unset, child.Status().Code)
		require.Len(t, child.Events(), 1)
		assert.Equal(t, "exec", child.Events()[0].Name)
	})

	t.Run("threads get no span", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		tr.startProcess(100, 1, now, shell)
		tr.addThread(102, 100)
		tr.startProcess(103, 102, now, procInfo{})

		assert.False(t, tr.exit(102, now, 0, ""), "thread exit is not a process exit")
		assert.True(t, tr.exit(103, now, 0, ""))
		assert.True(t, tr.exit(100, now, 0, ""))

		spans := recorder.Ended()
		require.Len(t, spans, 2)
		child := spans[0]
		assert.Equal(t, int64(100), attr(child, "process.parent_pid").AsInt64(), "parented to the thread's process")
		assert.Equal(t, spans[1].SpanContext().SpanID(), child.Parent().SpanID())
	})

	t.Run("failures", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		tr.startProcess(100, 1, now, shell)
		tr.startProcess(101, 100, now, procInfo{})
		tr.exit(101, now, 137, "SIGKILL")
		tr.exit(100, now, 3, "")

		spans := recorder.Ended()
		require.Len(t, spans, 2)
		assert.Equal(t, int64(137), attr(spans[0], "process.exit.code").AsInt64())
		assert.Equal(t, sdktrace.Status{Code: codes.Error, Description: "terminated by signal SIGKILL"}, spans[0].Status())
		assert.Equal(t, sdktrace.Status{Code: codes.Error, Description: "exit code 3"}, spans[1].Status())
	})

	t.Run("unknown task", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		assert.False(t, tr.exit(42, now, 0, ""))
		tr.exec(42, now, ls)
		assert.Empty(t, recorder.Ended())
	})

	t.Run("end all", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		tr.startProcess(100, 1, now, shell)
		tr.startProcess(101, 100, now, procInfo{})
		tr.addThread(102, 100)
		tr.endAll(now)

		assert.Len(t, recorder.Ended(), 2)
		assert.Empty(t, tr.procs)
		assert.Empty(t, tr.owner)
	})

	t.Run("span name fallbacks", func(t *testing.T) {
		assert.Equal(t, "ls", procInfo{comm: "ls", path: "/usr/bin/coreutils"}.spanName())
		assert.Equal(t, "coreutils", procInfo{path: "/usr/bin/coreutils"}.spanName())
		assert.Equal(t, "ls", procInfo{path: "/usr/bin/coreutils", args: []string{""}, comm: "ls"}.spanName())
	})

	t.Run("unnamed process", func(t *testing.T) {
		tr, recorder := newTestProcessTree(t)

		tr.startProcess(100, 1, now, procInfo{})
		tr.endAll(now)

		require.Len(t, recorder.Ended(), 1)
		assert.Equal(t, "process", recorder.Ended()[0].Name())
	})
}
