package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"tts-mcp/internal/config"
	"tts-mcp/internal/mcp"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	// Automatically parse environment variables from the global path
	_ = godotenv.Load(config.GetEnvPath())

	// Begin the STDIO processing loop
	mcp.Start()
}
