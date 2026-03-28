package main

import (
	"github.com/joho/godotenv"

	"tts-mcp/internal/api"
)

func main() {
	// Try to load .env file if it exists, purely to set process os.Env dynamically if not provided by host
	_ = godotenv.Load()

	// Begin the STDIO processing loop
	api.Start()
}
