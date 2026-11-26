package cli

import (
	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/cli/traceexeccmd"
	"github.com/spf13/cobra"
)

// platformCommands returns the commands that only exist on this platform.
func platformCommands(opts *app.Options) []*cobra.Command {
	return []*cobra.Command{traceexeccmd.New(opts)}
}
