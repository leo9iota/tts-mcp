set dotenv-load

# The default recipe when just is invoked with no arguments
default:
    @just --list

# ==========================================
# Build & Execution
# ==========================================

# Build the generic tts-mcp executable for the current host OS
build: deps
    go build -o bin/tts-mcp.exe ./cmd/tts-mcp

# Build cross-platform binaries (Windows, Linux, macOS)
build-all: deps
    GOOS=windows GOARCH=amd64 go build -o bin/tts-mcp-windows-amd64.exe ./cmd/tts-mcp
    GOOS=linux GOARCH=amd64 go build -o bin/tts-mcp-linux-amd64 ./cmd/tts-mcp
    GOOS=darwin GOARCH=arm64 go build -o bin/tts-mcp-darwin-arm64 ./cmd/tts-mcp

# Run the MCP server locally over stdio
run: build
    ./bin/tts-mcp.exe

# Start the MCP inspector to test the server visually
inspect: build
    bunx @modelcontextprotocol/inspector ./bin/tts-mcp.exe

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

# Configure the MCP server API keys interactively
config: deps
    go run ./cmd/config

# Setup the environment using the configurator (alias mapping)
setup-env: config

# Initialize a completely fresh developer environment
init: setup-env deps build
    @echo "Development environment initialized successfully!"
