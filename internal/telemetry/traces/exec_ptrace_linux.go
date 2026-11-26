package traces

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/neonmei/szgen/internal/runner"
	"golang.org/x/sys/unix"
)

const execSupported = true

// ptraceOptions follow every new task and stop on exec. EXITKILL kills the tracees if szgen dies.
const ptraceOptions = unix.PTRACE_O_EXITKILL |
	unix.PTRACE_O_TRACEFORK |
	unix.PTRACE_O_TRACEVFORK |
	unix.PTRACE_O_TRACECLONE |
	unix.PTRACE_O_TRACEEXEC |
	unix.PTRACE_O_TRACEEXIT

// Execute runs the command under ptrace until it and all its descendants have exited.
// Cancelling ctx forwards the signal that cancelled it (SIGTERM otherwise) to the command.
func (t *execTask) Execute(ctx context.Context) error {
	// ptrace requests are only accepted from the thread that started the tracee.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	slog.Info("Tracing command", "command", t.argv[0], "args", t.argv[1:])

	cmd := exec.Command(t.argv[0], t.argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		if errors.Is(err, syscall.EPERM) {
			return fmt.Errorf("start %s: %w (ptrace may be restricted by kernel.yama.ptrace_scope or a seccomp policy)", t.argv[0], err)
		}
		return fmt.Errorf("start %s: %w", t.argv[0], err)
	}
	defer func() { _ = cmd.Process.Release() }()

	rootPid := cmd.Process.Pid
	abort := func(err error) error {
		_ = unix.Kill(rootPid, unix.SIGKILL)
		_, _ = unix.Wait4(rootPid, nil, unix.WALL, nil)
		return err
	}

	// The tracee stops with SIGTRAP right after execve, before running any of its own code.
	var ws unix.WaitStatus
	if _, err := unix.Wait4(rootPid, &ws, unix.WALL, nil); err != nil {
		return abort(fmt.Errorf("wait for %s to start: %w", t.argv[0], err))
	}
	if !ws.Stopped() {
		return fmt.Errorf("%s exited before it could be traced", t.argv[0])
	}

	tr := newProcessTree(context.WithoutCancel(ctx), t.tracer())
	tr.startProcess(rootPid, os.Getpid(), start, readProcInfo(rootPid))

	if err := unix.PtraceSetOptions(rootPid, ptraceOptions); err != nil {
		tr.endAll(time.Now())
		return abort(fmt.Errorf("set ptrace options: %w", err))
	}
	if err := unix.PtraceCont(rootPid, 0); err != nil {
		tr.endAll(time.Now())
		return abort(fmt.Errorf("resume %s: %w", t.argv[0], err))
	}

	stopForwarding := forwardSignal(ctx, cmd.Process)
	rootCode, err := waitTracees(tr, rootPid)
	stopForwarding()
	tr.endAll(time.Now())
	if err != nil {
		return err
	}

	slog.Info("Command finished", "pid", rootPid, "exit_code", rootCode, "duration", time.Since(start))
	if rootCode != 0 {
		return &ExitError{Code: rootCode}
	}
	return nil
}

// waitTracees handles ptrace stops of every tracee until none is left, and returns the root process exit code.
func waitTracees(tr *processTree, rootPid int) (int, error) {
	rootCode := 0
	// szgen may have children it did not start (e.g. inherited from a shell that forked before exec'ing it),
	// so the loop ends when every tracee is gone rather than when wait4 runs out of children.
	alive := map[int]bool{rootPid: true}
	// New tracees start with a SIGSTOP that must be suppressed. It may be reported before the
	// fork event of its parent, so it is tracked on its own.
	started := map[int]bool{rootPid: true}

	for len(alive) > 0 {
		var ws unix.WaitStatus
		wpid, err := unix.Wait4(-1, &ws, unix.WALL, nil)
		switch {
		case errors.Is(err, unix.ECHILD):
			return rootCode, nil
		case errors.Is(err, unix.EINTR):
			continue
		case err != nil:
			return rootCode, fmt.Errorf("wait for traced processes: %w", err)
		}
		now := time.Now()

		if ws.Exited() || ws.Signaled() {
			code, signal := ws.ExitStatus(), ""
			if ws.Signaled() {
				code, signal = consts.SignalExitCodeBase+int(ws.Signal()), unix.SignalName(ws.Signal())
			}
			if wpid == rootPid {
				rootCode = code
			}
			if tr.exit(wpid, now, code, signal) {
				slog.Debug("Process exited", "pid", wpid, "exit_code", code, "signal", signal)
			}
			delete(started, wpid)
			delete(alive, wpid)
			continue
		}

		if !ws.Stopped() {
			continue
		}
		alive[wpid] = true

		deliver := 0
		switch cause := ws.TrapCause(); cause {
		case unix.PTRACE_EVENT_FORK, unix.PTRACE_EVENT_VFORK, unix.PTRACE_EVENT_CLONE:
			if id, ok := newTracee(tr, wpid, cause, now); ok {
				alive[id] = true
			}
		case unix.PTRACE_EVENT_EXEC:
			// A thread calling exec takes over the process id. Its former thread id never reports an exit.
			if former, err := unix.PtraceGetEventMsg(wpid); err == nil && int(former) != wpid {
				tr.exit(int(former), now, 0, "")
				delete(alive, int(former))
				delete(started, int(former))
			}
			info := readProcInfo(wpid)
			tr.exec(wpid, now, info)
			slog.Debug("Process exec", "pid", wpid, "executable", info.path)
		case unix.PTRACE_EVENT_EXIT:
			// The exit status is recorded once the process is gone.
		default:
			if stopSig := ws.StopSignal(); stopSig == unix.SIGSTOP && !started[wpid] {
				started[wpid] = true
			} else {
				deliver = int(stopSig)
			}
		}

		if err := unix.PtraceCont(wpid, deliver); err != nil {
			slog.Debug("Failed to resume traced process", "pid", wpid, "error", err)
		}
	}

	return rootCode, nil
}

// newTracee records the task created by the tracee parentID, reported by a fork, vfork or clone event,
// and returns its id.
func newTracee(tr *processTree, parentID, cause int, at time.Time) (int, bool) {
	msg, err := unix.PtraceGetEventMsg(parentID)
	if err != nil {
		slog.Debug("Failed to get new task id", "parent", parentID, "error", err)
		return 0, false
	}
	id := int(msg)

	if cause == unix.PTRACE_EVENT_CLONE {
		tgid, err := threadGroupID(id)
		if err != nil {
			slog.Debug("Failed to read thread group", "tid", id, "error", err)
			return id, true
		}
		if tgid != id {
			tr.addThread(id, tgid)
			return id, true
		}
	}

	tr.startProcess(id, parentID, at, procInfo{})
	slog.Debug("Process started", "pid", id, "parent", parentID)
	return id, true
}

// forwardSignal sends the signal that cancels ctx to the process, until stop is called.
func forwardSignal(ctx context.Context, process *os.Process) (stop func()) {
	done := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			sig, ok := runner.SignalCause(ctx)
			if !ok {
				sig = syscall.SIGTERM
			}
			slog.Info("Forwarding signal", "signal", sig, "pid", process.Pid)
			if err := process.Signal(sig); err != nil {
				slog.Debug("Failed to forward signal", "pid", process.Pid, "error", err)
			}
		case <-done:
		}
	}()

	return func() { close(done) }
}
