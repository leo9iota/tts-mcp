package personas

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Persona struct {
	Name     string `json:"name"`
	Trope    string `json:"trope"`
	Provider string `json:"provider"`
	VoiceID  string `json:"voice_id"`
}

type Manager struct {
	Personas    map[string]Persona
	PersonasDir string
}

func NewManager() (*Manager, error) {
	m := &Manager{
		Personas: make(map[string]Persona),
	}

	searchDirs := []string{"."}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Typically run via 'go run .' or compiled to 'bin/', so test both project root fallbacks
		searchDirs = append([]string{filepath.Join(exeDir, ".."), exeDir}, searchDirs...)
	}

	var personasDir string
	for _, sDir := range searchDirs {
		candidate := filepath.Join(sDir, "data")
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			personasDir = candidate
			break
		}
	}

	if personasDir == "" {
		// Fallback: create it in the CWD if none exists in typical places
		personasDir = filepath.Join(".", "data")
		if err := os.MkdirAll(personasDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	m.PersonasDir = personasDir

	files, err := os.ReadDir(personasDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read personas directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join(personasDir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading persona %s: %v\n", path, err)
			continue
		}

		var p Persona
		if err := json.Unmarshal(data, &p); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing persona %s: %v\n", path, err)
			continue
		}

		if p.Name != "" {
			m.Personas[p.Name] = p
		}
	}

	return m, nil
}

func (m *Manager) GetOptions() []string {
	var opts []string
	for name := range m.Personas {
		opts = append(opts, name)
	}
	// Return a default empty string if no personas configured yet to prevent MCP schema validation errors
	if len(opts) == 0 {
		return []string{""}
	}
	return opts
}

func (m *Manager) GetPersona(name string) (Persona, bool) {
	p, exists := m.Personas[name]
	return p, exists
}

// SavePersona maps a structural persona payload directly into the data folder as a perfectly structured JSON artifact.
func (m *Manager) SavePersona(p Persona) error {
	if p.Name == "" {
		return fmt.Errorf("persona must have a valid name")
	}

	fileName := strings.ToLower(strings.ReplaceAll(p.Name, " ", "_")) + ".json"
	path := filepath.Join(m.PersonasDir, fileName)

	bytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal persona: %w", err)
	}

	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write persona file: %w", err)
	}

	// Hot reload into the active execution pointer immediately
	m.Personas[p.Name] = p
	return nil
}
