// Package traces builds tasks that emit traces.
//
// The exec task (Linux only) runs a command and records its process tree as an OpenTelemetry trace,
// with one span per process from fork (or start) to exit.
package traces
