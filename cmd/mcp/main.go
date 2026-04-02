package main

import (
	"github.com/joho/godotenv"

	"tts-mcp/internal/config"
	"tts-mcp/internal/mcp"
)

func main() {
	// Automatically parse environment variables from the global path
	_ = godotenv.Load(config.GetEnvPath())

	// Begin the STDIO processing loop
	mcp.Start()
}
