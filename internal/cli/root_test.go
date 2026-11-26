package cli

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutorOverrides(t *testing.T) {
	parse := func(t *testing.T, args ...string) *pflag.FlagSet {
		t.Helper()

		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		flags.StringP("executor", "e", consts.DefaultExecutorStrategy, "")
		flags.IntP("max-concurrency", "j", 0, "")
		require.NoError(t, flags.Parse(args))
		return flags
	}

	t.Run("unset flags", func(t *testing.T) {
		overrides := executorOverrides(parse(t))

		assert.Nil(t, overrides.Strategy)
		assert.Nil(t, overrides.MaxConcurrency)
	})

	t.Run("set flags, including their defaults", func(t *testing.T) {
		overrides := executorOverrides(parse(t, "-e", "serial", "-j", "0"))

		require.NotNil(t, overrides.Strategy)
		require.NotNil(t, overrides.MaxConcurrency)
		assert.Equal(t, consts.ExecutorStrategySerial, *overrides.Strategy)
		assert.Equal(t, 0, *overrides.MaxConcurrency)
	})
}

func TestNewRootCommand(t *testing.T) {
	execute := func(t *testing.T, args ...string) (string, error) {
		t.Helper()

		var out bytes.Buffer
		cmd := NewRootCommand("1.2.3")
		cmd.SetOut(&out)
		cmd.SetArgs(args)

		err := cmd.Execute()
		return out.String(), err
	}

	t.Run("version", func(t *testing.T) {
		out, err := execute(t, "--version")
		require.NoError(t, err)
		assert.Equal(t, "szgen version 1.2.3\n", out)
	})

	t.Run("commands", func(t *testing.T) {
		cmd := NewRootCommand("test")
		for _, name := range []string{"run", "r", "metrics", "m"} {
			found, _, err := cmd.Find([]string{name})
			require.NoError(t, err, name)
			assert.NotEqual(t, cmd, found, name)
		}
	})

	t.Run("trace-exec only exists on linux", func(t *testing.T) {
		if runtime.GOOS == "linux" {
			found, _, err := NewRootCommand("test").Find([]string{"trace-exec"})
			require.NoError(t, err)
			assert.Equal(t, "trace-exec", found.Name())
			return
		}

		_, err := execute(t, "trace-exec", "true")
		assert.EqualError(t, err, `unknown command "trace-exec" for "szgen"`)
	})

	t.Run("run requires a config file", func(t *testing.T) {
		_, err := execute(t, "run")
		assert.EqualError(t, err, "failed to load config: no config file path provided")
	})

	t.Run("run with a missing config file", func(t *testing.T) {
		_, err := execute(t, "run", "--config", "/nonexist.yaml")
		assert.EqualError(t, err, "failed to load config: failed to read szgen config file /nonexist.yaml: open /nonexist.yaml: no such file or directory")
	})

	t.Run("config flag only exists on run", func(t *testing.T) {
		_, err := execute(t, "metrics", "counter", "--config", "x")
		assert.EqualError(t, err, "unknown flag: --config")
	})
}
