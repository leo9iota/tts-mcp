package main

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"tts-mcp/internal/api"
)

func main() {
	// Dynamically load .env from the executable's directory or the parent directory.
	// This ensures the API keys are loaded regardless of the client's working directory 
	// (e.g., when the MCP Inspector or Claude runs it from an arbitrary temp path).
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		_ = godotenv.Load(filepath.Join(exeDir, ".env"))       // Check alongside binary
		_ = godotenv.Load(filepath.Join(exeDir, "..", ".env")) // Check project root
	}
	// Fallback to Current Working Directory
	_ = godotenv.Load()

	// Begin the STDIO processing loop
	api.Start()
}
