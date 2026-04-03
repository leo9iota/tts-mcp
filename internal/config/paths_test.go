package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAppConfigDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("APPDATA", tempHome)         // Windows
	t.Setenv("XDG_CONFIG_HOME", tempHome) // Linux
	t.Setenv("HOME", tempHome)            // macOS

	dir := GetAppConfigDir()
	if !strings.HasPrefix(dir, tempHome) {
		t.Errorf("expected dir to start with %s, got %s", tempHome, dir)
	}
	if !strings.HasSuffix(dir, "tts-mcp") {
		t.Errorf("expected dir to end with tts-mcp, got %s", dir)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected directory to be created, got error: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected returned value to be a directory")
	}
}

func TestGetCacheDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("LOCALAPPDATA", tempHome)   // Windows
	t.Setenv("XDG_CACHE_HOME", tempHome) // Linux
	t.Setenv("HOME", tempHome)           // macOS

	dir := GetCacheDir()
	expectedSuffix := filepath.Join("tts-mcp", "output")
	if !strings.HasSuffix(dir, expectedSuffix) {
		t.Errorf("expected dir to end with %s, got %s", expectedSuffix, dir)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected directory to be created, got error: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected returned path to be a directory")
	}
}

func TestGetPersonasDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("APPDATA", tempHome)
	t.Setenv("XDG_CONFIG_HOME", tempHome)

	dir := GetPersonasDir()
	if !strings.HasSuffix(dir, filepath.Join("tts-mcp", "personas")) {
		t.Errorf("expected personas dir suffix, got: %s", dir)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected personas directory to be created, got err: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected personas path to be a directory")
	}
}

func TestGetEnvPath(t *testing.T) {
	p := GetEnvPath()
	if filepath.Base(p) != ".env" {
		t.Errorf("expected env file name to be .env, got: %s", filepath.Base(p))
	}
	if !strings.Contains(p, "tts-mcp") {
		t.Errorf("expected env file to be inside tts-mcp app dir, got: %s", p)
	}
}
