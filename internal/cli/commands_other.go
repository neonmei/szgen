//go:build !linux

package cli

import (
	"github.com/neonmei/szgen/internal/app"
	"github.com/spf13/cobra"
)

// platformCommands returns the commands that only exist on this platform.
func platformCommands(*app.Options) []*cobra.Command {
	return nil
}
