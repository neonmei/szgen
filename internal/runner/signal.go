package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// SignalError is the cancellation cause of a context cancelled by a process signal.
type SignalError struct {
	Signal os.Signal
}

func (e *SignalError) Error() string {
	return fmt.Sprintf("received signal %s", e.Signal)
}

// SignalCause returns the signal that cancelled ctx, if it was cancelled by one.
func SignalCause(ctx context.Context) (os.Signal, bool) {
	var sigErr *SignalError
	if errors.As(context.Cause(ctx), &sigErr) {
		return sigErr.Signal, true
	}
	return nil, false
}
