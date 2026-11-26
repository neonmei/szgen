package traces

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// procInfo describes the program a process is running.
type procInfo struct {
	path string   // executable path
	comm string   // kernel command name, truncated to 15 bytes
	args []string // command line
}

// executableName follows semantic conventions: the base name of /proc/<pid>/exe.
func (i procInfo) executableName() string {
	if i.path != "" {
		return filepath.Base(i.path)
	}
	return i.comm
}

// spanName is the command as invoked, which differs from the executable name for symlinks
// and multi-call binaries (busybox, coreutils).
func (i procInfo) spanName() string {
	if len(i.args) > 0 && i.args[0] != "" {
		return filepath.Base(i.args[0])
	}
	if i.comm != "" {
		return i.comm
	}
	return i.executableName()
}

func (i procInfo) attributes() []attribute.KeyValue {
	var attrs []attribute.KeyValue

	if name := i.executableName(); name != "" {
		attrs = append(attrs, semconv.ProcessExecutableName(name))
	}
	if i.path != "" {
		attrs = append(attrs, semconv.ProcessExecutablePath(i.path))
	}
	if len(i.args) > 0 {
		attrs = append(attrs, semconv.ProcessCommand(i.args[0]), semconv.ProcessCommandArgs(i.args...))
	}

	return attrs
}

type proc struct {
	ctx  context.Context
	span trace.Span
	info procInfo
}

// processTree maps the traced processes to their spans. Threads get no span of their own,
// they are only tracked so processes they fork are parented to the process that owns them.
type processTree struct {
	tracer trace.Tracer
	root   context.Context
	procs  map[int]*proc
	owner  map[int]int // thread id → pid of its process
}

func newProcessTree(ctx context.Context, tracer trace.Tracer) *processTree {
	return &processTree{
		tracer: tracer,
		root:   ctx,
		procs:  make(map[int]*proc),
		owner:  make(map[int]int),
	}
}

// pid resolves a thread id to the pid of its process.
func (t *processTree) pid(id int) int {
	if pid, ok := t.owner[id]; ok {
		return pid
	}
	return id
}

// startProcess starts the span of a new process, as a child of its parent's span if the parent is traced.
// A forked process runs its parent's program until it execs, so an empty info is inherited from the parent.
//
//nolint:spancheck // the span lives in the tree until the process exits, it is ended by exit or endAll.
func (t *processTree) startProcess(pid, parentID int, at time.Time, info procInfo) {
	ppid := t.pid(parentID)
	parentCtx := t.root
	if parent, ok := t.procs[ppid]; ok {
		parentCtx = parent.ctx
		if info.spanName() == "" {
			info = parent.info
		}
	}

	name := info.spanName()
	if name == "" {
		name = "process"
	}

	attrs := append([]attribute.KeyValue{semconv.ProcessPID(pid), semconv.ProcessParentPID(ppid)}, info.attributes()...)
	ctx, span := t.tracer.Start(parentCtx, name,
		trace.WithTimestamp(at),
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(attrs...),
	)

	t.procs[pid] = &proc{ctx: ctx, span: span, info: info}
}

// addThread records a thread of the process pid.
func (t *processTree) addThread(tid, pid int) {
	t.owner[tid] = t.pid(pid)
}

// exec records a process replacing its program.
func (t *processTree) exec(id int, at time.Time, info procInfo) {
	p, ok := t.procs[t.pid(id)]
	if !ok {
		return
	}

	p.info = info
	if name := info.spanName(); name != "" {
		p.span.SetName(name)
	}
	p.span.SetAttributes(info.attributes()...)
	p.span.AddEvent("exec", trace.WithTimestamp(at), trace.WithAttributes(info.attributes()...))
}

// exit ends the span of a process that exited with code, or was killed by signal when it is not empty.
// It reports whether id was a traced process (as opposed to a thread or an unknown task).
func (t *processTree) exit(id int, at time.Time, code int, signal string) bool {
	if _, ok := t.owner[id]; ok {
		delete(t.owner, id)
		return false
	}

	p, ok := t.procs[id]
	if !ok {
		return false
	}

	p.span.SetAttributes(semconv.ProcessExitCode(code))
	switch {
	case signal != "":
		p.span.SetStatus(codes.Error, "terminated by signal "+signal)
	case code != 0:
		p.span.SetStatus(codes.Error, fmt.Sprintf("exit code %d", code))
	}
	p.span.End(trace.WithTimestamp(at))

	delete(t.procs, id)
	for tid, pid := range t.owner {
		if pid == id {
			delete(t.owner, tid)
		}
	}

	return true
}

// endAll ends the spans of processes whose exit was never observed.
func (t *processTree) endAll(at time.Time) {
	for pid, p := range t.procs {
		p.span.End(trace.WithTimestamp(at))
		delete(t.procs, pid)
	}
	clear(t.owner)
}
