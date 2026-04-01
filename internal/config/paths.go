package config

import (
	"os"
	"path/filepath"
)

// GetAppConfigDir resolves the base XDG configuration directory for tts-mcp.
func GetAppConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "." // Fallback to PWD if OS completely fails (rare)
	}
	appDir := filepath.Join(configDir, "tts-mcp")
	_ = os.MkdirAll(appDir, 0o755)
	return appDir
}

// GetEnvPath resolves the `.env` file location explicitly into the core configuration space.
func GetEnvPath() string {
	return filepath.Join(GetAppConfigDir(), ".env")
}

// GetPersonasDir resolves the personas directory, ensuring it exists before mapping operations.
func GetPersonasDir() string {
	dir := filepath.Join(GetAppConfigDir(), "personas")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// GetCacheDir returns the appropriate XDG cache base directory mapping for output logs and artifacts.
func GetCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "." // Fallback
	}
	dir := filepath.Join(cacheDir, "tts-mcp", "output")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}
