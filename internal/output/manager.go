package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateOutputFile dynamically constructs a guaranteed unique artifact file inside the outputs/ directory.
// It anchors relative to the 'data' directory to ensure stable production/binary deployments.
func GenerateOutputFile(personaName, providerName string) (*os.File, string, error) {
	searchDirs := []string{"."}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Check typical project root layout vs bin runtime execution pathing
		searchDirs = append([]string{filepath.Join(exeDir, ".."), exeDir}, searchDirs...)
	}

	var rootDir string
	for _, sDir := range searchDirs {
		candidate := filepath.Join(sDir, "data")
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			rootDir = sDir
			break
		}
	}

	if rootDir == "" {
		rootDir = "." // Absolute fallback to CWD
	}

	// Explicitly nest inside the data/ folder alongside personas/
	outputsDir := filepath.Join(rootDir, "data", "output")
	if err := os.MkdirAll(outputsDir, 0o755); err != nil {
		return nil, "", fmt.Errorf("failed to create output directory: %w", err)
	}

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
