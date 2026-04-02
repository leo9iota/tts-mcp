package main

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"tts-mcp/internal/personas"
)

func RunPersonaWizard() {
	mng, err := personas.NewManager()
	if err != nil {
		PrintError(fmt.Sprintf("Error initializing persona manager: %v", err))
		return
	}

	for {
		var action string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Persona Management").
					Description("Manage your dynamic character profiles in "+mng.PersonasDir).
					Options(
						huh.NewOption("Create New Persona", "CREATE"),
						huh.NewOption("Edit Persona", "EDIT"),
						huh.NewOption("Delete Persona", "DELETE"),
						huh.NewOption("List Personas", "LIST"),
						huh.NewOption("Return to Main Menu", "BACK"),
					).
					Value(&action),
			),
		).WithTheme(OneDarkTheme()).Run()

		if err != nil || action == "BACK" {
			return
		}

		switch action {
		case "CREATE":
			createPersona(mng, nil)
		case "EDIT":
			p, ok := selectPersona(mng, "Edit")
			if ok {
				createPersona(mng, &p)
			}
		case "DELETE":
			p, ok := selectPersona(mng, "Delete")
			if ok {
				deletePersona(mng, p)
			}
		case "LIST":
			listPersonas(mng)
		}
	}
}
