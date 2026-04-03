# Contributing to TTS-MCP

First off, thank you for considering contributing to TTS-MCP! It's people like you that make TTS-MCP an amazing tool for everyone.

## Development Workflow

The project utilizes `just` (a command runner) to manage environments, builds, and dependencies. You should have [`just`](https://github.com/casey/just) installed on your system.

### 1. Initial Setup

Run the following to pull down dependencies and configure your environment:

```bash
just init
```

### 2. Building

Compile the `tts-mcp` code into binaries for both the core MCP server and the configuration CLI:

```bash
just build
```

### 3. Testing and Linting

Before submitting any Pull Request, ensure that your code is formatted and passes tests:

```bash
just test
# Run with race detection
just test-race
# Format and lint
just lint
```

## Pull Request Process

1. Ensure any new features include unit tests (if applicable) under the `internal/` or `cmd/` packages.
2. Update the `README.md` if your PR introduces new user-facing functionality or config parameters.
3. Keep PRs scope-limited to a single meaningful change (bugfix, feature, or refactor).
