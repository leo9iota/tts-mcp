package main

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"tts-mcp/internal/personas"
)

func selectPersona(mng *personas.Manager, action string) (personas.Persona, bool) {
	opts := mng.GetOptions()
	if len(opts) == 1 && opts[0] == "" {
		PrintInfo("No personas found.")
		return personas.Persona{}, false
	}

	var selected string
	options := make([]huh.Option[string], 0, len(opts)+1)
	for _, o := range opts {
		options = append(options, huh.NewOption(o, o))
	}
	options = append(options, huh.NewOption("Go Back", ""))

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("%s Persona", action)).
				Options(options...).
				Value(&selected),
		),
	).WithTheme(OneDarkTheme()).Run()

	if err != nil || selected == "" {
		return personas.Persona{}, false
	}

	p, _ := mng.GetPersona(selected)
	return p, true
}
