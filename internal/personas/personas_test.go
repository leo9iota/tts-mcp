package personas

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestManager_SavePersona(t *testing.T) {
	// 1. Arrange: Create a strictly localized isolated mock environment
	tempDir := t.TempDir()

	m := &Manager{
		Personas:    make(map[string]Persona),
		PersonasDir: tempDir,
	}

	testPersona := Persona{
		Name:     "Test Waifu",
		Trope:    "Sassy tester",
		Provider: "fishaudio_tts",
		VoiceID:  "hex123abc",
		Options: map[string]interface{}{
			"latency": "fast",
		},
	}

	// 2. Act: Execute the storage boundary loop
	err := m.SavePersona(testPersona)
	if err != nil {
		t.Fatalf("SavePersona returned unexpected error: %v", err)
	}

	// 3. Assert: Memory State Maps properly bind
	if p, exists := m.GetPersona("Test Waifu"); !exists {
		t.Errorf("Expected 'Test Waifu' natively inside mapping, got missing")
	} else if p.Provider != "fishaudio_tts" {
		t.Errorf("Expected Provider 'fishaudio_tts', got %s", p.Provider)
	}

	// 4. Assert: Physical File System accurately structures artifact
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to explicitly read isolated TempDir: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected exactly 1 metadata file bounded by SavePersona, got %d", len(files))
	}

	// fileName is dynamically slugged by `strings.ToLower(strings.ReplaceAll(...)`
	expectedFileName := "test_waifu.json"
	if files[0].Name() != expectedFileName {
		t.Errorf("Expected filename '%s', got '%s'", expectedFileName, files[0].Name())
	}

	// 5. Assert: Verify the written byte blocks didn't drop nested maps like 'options' due to structural marshaling drift
	fileBytes, _ := os.ReadFile(filepath.Join(tempDir, expectedFileName))
	var unmarshaled Persona
	if err := json.Unmarshal(fileBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unpack the newly minted generic artifact structurally: %v", err)
	}

	if unmarshaled.Options["latency"] != "fast" {
		t.Errorf("Expected explicit internal option latency='fast' dynamically, got %v", unmarshaled.Options["latency"])
	}
}

func TestManager_GetOptions(t *testing.T) {
	m := &Manager{
		Personas: make(map[string]Persona),
	}

	// Default fallback prevents MCP protocol schema from crushing under empty constraints
	opts := m.GetOptions()
	if len(opts) != 1 || opts[0] != "" {
		t.Errorf("Expected default enum hook [''] specifically preventing schema panics, got %v", opts)
	}

	// Adding structural dependencies manually
	m.Personas["Megumin"] = Persona{Name: "Megumin"}
	m.Personas["Geralt"] = Persona{Name: "Geralt"}

	opts = m.GetOptions()
	if len(opts) != 2 {
		t.Errorf("Expected exactly 2 generated schema objects mapped tightly from dictionary, got %v", opts)
	}

	// Dictionaries are implicitly unordered in Go, just ensure presence explicitly
	foundMegumin, foundGeralt := false, false
	for _, opt := range opts {
		if opt == "Megumin" {
			foundMegumin = true
		}
		if opt == "Geralt" {
			foundGeralt = true
		}
	}

	if !foundMegumin || !foundGeralt {
		t.Errorf("Missing exactly mapped keys natively: Megumin=%v, Geralt=%v", foundMegumin, foundGeralt)
	}
}
