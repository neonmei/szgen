package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/neonmei/szgen/internal/cli"
)

// version is set at build time via ldflags:
//
//	go build -ldflags "-X main.version=1.2.3"
var version = "dev"

func main() {
	if err := cli.NewRootCommand(version).Execute(); err != nil {
		// A traced command reports its own failure, szgen only passes its exit code on.
		var exitErr interface{ ExitCode() int }
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
