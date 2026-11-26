package traces

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// ExitError reports a traced command that exited with a non-zero code (128+N when killed by signal N).
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("command exited with code %d", e.Code)
}

// ExitCode is the exit code szgen should exit with.
func (e *ExitError) ExitCode() int {
	return e.Code
}

type execTask struct {
	argv []string
}

// New creates a task that runs argv and traces it and all its descendant processes.
func NewExecTask(argv []string) (runner.Task, error) {
	if len(argv) == 0 {
		return nil, errors.New("no command provided")
	}

	if _, err := exec.LookPath(argv[0]); err != nil {
		return nil, err
	}

	return &execTask{argv: argv}, nil
}

func (t *execTask) Name() string {
	return t.argv[0]
}

// tracer is resolved when the task runs, so the global provider is the one installed by the SDK.
func (t *execTask) tracer() trace.Tracer {
	return otel.Tracer(consts.DefaultTracerName)
}
