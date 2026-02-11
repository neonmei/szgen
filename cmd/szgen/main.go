package main

import (
	"fmt"
	"os"
)

// version is set at build time via ldflags:
//
//	go build -ldflags "-X main.version=1.2.3"
var version = "dev"

func main() {
	rootCmd.Version = version

	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
