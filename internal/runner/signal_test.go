package runner

import (
	"context"
	"fmt"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignalCause(t *testing.T) {
	t.Run("cancelled by signal", func(t *testing.T) {
		ctx, cancel := context.WithCancelCause(context.Background())
		cancel(fmt.Errorf("wrapped: %w", &SignalError{Signal: syscall.SIGINT}))

		sig, ok := SignalCause(ctx)
		assert.True(t, ok)
		assert.Equal(t, syscall.SIGINT, sig)
		assert.ErrorIs(t, ctx.Err(), context.Canceled)
	})

	t.Run("cancelled otherwise", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, ok := SignalCause(ctx)
		assert.False(t, ok)
	})

	t.Run("not cancelled", func(t *testing.T) {
		_, ok := SignalCause(context.Background())
		assert.False(t, ok)
	})
}
