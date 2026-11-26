package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/neonmei/szgen/internal/runner"
)

// signalContext returns a context cancelled on the first SIGINT or SIGTERM, with a *runner.SignalError as cause.
// The handler is removed afterwards, so a second signal terminates the process.
func signalContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigCh:
			slog.Info("Received signal, shutting down gracefully", "signal", sig)
			cancel(&runner.SignalError{Signal: sig})
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
		close(sigCh)
	}()

	return ctx, func() { cancel(context.Canceled) }
}
