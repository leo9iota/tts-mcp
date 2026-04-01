package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tts-mcp/internal/config"
)

// GenerateOutputFile dynamically constructs a guaranteed unique artifact file inside the outputs/ directory.
// It anchors relative to the global OS user cache directory to ensure stable production/binary deployments.
func GenerateOutputFile(personaName, providerName string) (*os.File, string, error) {
	outputsDir := config.GetCacheDir()

	// Formatting logic (FEAT-002)
	prefix := personaName
	if prefix == "" {
		prefix = providerName
	}
	if prefix == "" {
		prefix = "audio"
	}

	// Purge OS-invalid filename strings
	prefix = strings.ReplaceAll(prefix, " ", "_")
	prefix = strings.ReplaceAll(prefix, "/", "")
	prefix = strings.ReplaceAll(prefix, "\\", "")

	// Include milliseconds to physically guarantee race-condition uniqueness natively
	timestamp := time.Now().Format("2006-01-02_15-04-05.000")
	// Clean the colons and dots since Windows hates them in explicitly constructed file extensions
	timestamp = strings.ReplaceAll(timestamp, ".", "-")

	fileName := fmt.Sprintf("%s_%s.mp3", prefix, timestamp)
	absPath, _ := filepath.Abs(filepath.Join(outputsDir, fileName))

	file, err := os.Create(absPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to securely lock output file path %s: %w", absPath, err)
	}

	return file, absPath, nil
}
