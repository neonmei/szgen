//go:build !linux

package traces

import "context"

const execSupported = false

func (t *execTask) Execute(context.Context) error {
	return errExecUnsupported
}
