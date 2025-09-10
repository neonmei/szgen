# Default recipe
default:
    @just --list

# Build the binary
build:
    go build -o bin/szgen ./cmd/szgen

# Install locally
install:
    go install ./cmd/szgen

# Run tests
test:
    go test -race ./...

# Run tests without race detection
test-norace:
    go test ./...

# Run tests with coverage
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
    go test -v -bench=. ./...

# Run linting with golangci-lint
lint:
    golangci-lint run

# Format code
fmt:
    go fmt ./...
    gofumpt -w ./

# Vet code
vet:
    go vet ./...

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f coverage.out coverage.html

# Install dependencies
deps:
    go mod download
    go mod tidy

# Update dependencies
update-deps:
    go get -u ./...
    go mod tidy

# Run all checks (format, vet, lint, test)
check: fmt vet lint test

# Run example with counter
example-counter:
    just build
    ./bin/szgen metrics counter --name test_counter --value 1 --count 5 --rate 1s

# Run example with config file
example-config-run:
    just build
    ./bin/szgen run --config examples/example-config.yaml