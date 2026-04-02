set dotenv-load

# The default recipe when just is invoked with no arguments
default:
    @just --list

# ==========================================
# Build & Execution
# ==========================================

# Build the generic tts-mcp executable for the current host OS
build: deps
    go build -o bin/mcp.exe ./cmd/mcp
    go build -o bin/cli.exe ./cmd/cli

# Build cross-platform binaries (Windows, Linux, macOS)
build-all: deps
    GOOS=windows GOARCH=amd64 go build -o bin/mcp-windows-amd64.exe ./cmd/mcp
    GOOS=linux GOARCH=amd64 go build -o bin/mcp-linux-amd64 ./cmd/mcp
    GOOS=darwin GOARCH=arm64 go build -o bin/mcp-darwin-arm64 ./cmd/mcp

# Run the MCP server locally over stdio
mcp: build
    ./bin/mcp.exe

# Start the MCP inspector to test the server visually
inspect: build
    bunx @modelcontextprotocol/inspector ./bin/mcp.exe

# ==========================================
# Testing & Linting
# ==========================================

# Run the pure Go test suite
test:
    go test -v ./...

# Run the test suite with race detection enabled
test-race:
    go test -race -v ./...

# Run go vet to catch suspiciously constructed code
vet:
    go vet ./...

# Run standard formatting to clean up syntax
format:
    gofumpt -l -w .

# Run standard linting and formatting pipeline
lint: format vet
    @echo "Linting complete."

# ==========================================
# Maintenance & Utilities
# ==========================================

# Download modules and clean the go.mod and go.sum files
deps:
    go mod tidy
    go mod download

# Clean up binaries and temporary audio cache files
wipe:
    rm -rf bin/
    rm -f temp.wav temp.mp3

# Launch the interactive configuration wizard (CLI)
cli: deps
    go run ./cmd/cli

# Setup the environment using the configurator (alias mapping)
setup-env: cli

# Initialize a completely fresh developer environment
init: setup-env deps build
    @echo "Development environment initialized successfully!"
