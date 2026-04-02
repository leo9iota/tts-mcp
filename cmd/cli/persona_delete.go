package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"

	"tts-mcp/internal/personas"
)

func deletePersona(mng *personas.Manager, p personas.Persona) {
	var confirm bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete '%s'?", p.Name)).
				Value(&confirm),
		),
	).WithTheme(OneDarkTheme()).Run()

	if err != nil || !confirm {
		return
	}

	fileName := strings.ToLower(strings.ReplaceAll(p.Name, " ", "_")) + ".toml"
	path := filepath.Join(mng.PersonasDir, fileName)

	if err := os.Remove(path); err != nil {
		PrintError(fmt.Sprintf("Error deleting file: %v", err))
	} else {
		delete(mng.Personas, p.Name)
		PrintSuccess(fmt.Sprintf("Successfully deleted %s", path))
	}
}
