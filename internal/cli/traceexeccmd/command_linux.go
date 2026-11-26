package traceexeccmd

import (
	"github.com/neonmei/szgen/internal/app"
	"github.com/neonmei/szgen/internal/telemetry/traces"
	"github.com/spf13/cobra"
)

// New returns the `trace-exec` command, which runs a command and traces its process tree.
func New(opts *app.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trace-exec [flags] [--] <command> [args...]",
		Short: "Execute a command and trace its process tree",
		Long: `Execute a command and trace it and all its descendant processes using ptrace.
Each process becomes a span, parented to the span of the process that forked it, from its start
to its exit, with its pid, executable, command line and exit code as attributes.
szgen exits with the exit code of the command.`,
		Example: `  szgen trace-exec -- make -j4
  szgen --log-level debug trace-exec sh -c 'ls | wc -l'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			task, err := traces.NewExecTask(args)
			if err != nil {
				return err
			}

			cfg, err := app.Load(*opts)
			if err != nil {
				return err
			}

			return app.ExecuteTasks(cfg, task)
		},
	}

	// Flags after the command name belong to the command: `trace-exec ls -la`.
	cmd.Flags().SetInterspersed(false)

	return cmd
}
