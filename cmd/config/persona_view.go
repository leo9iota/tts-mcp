package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

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
					Description("Manage your dynamic character profiles in " + mng.PersonasDir).
					Options(
						huh.NewOption("List Personas", "LIST"),
						huh.NewOption("Create New Persona", "CREATE"),
						huh.NewOption("Edit Persona", "EDIT"),
						huh.NewOption("Delete Persona", "DELETE"),
						huh.NewOption("Return to Main Menu", "BACK"),
					).
					Value(&action),
			),
		).WithTheme(OneDarkTheme()).Run()

		if err != nil || action == "BACK" {
			return
		}

		switch action {
		case "LIST":
			listPersonas(mng)
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
		}
	}
}

func listPersonas(mng *personas.Manager) {
	opts := mng.GetOptions()
	if len(opts) == 1 && opts[0] == "" {
		PrintInfo("No personas found.")
		return
	}
	
	var listStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#61afef")).
		Padding(1, 4).MarginTop(1).MarginBottom(1)
		
	var content strings.Builder
	content.WriteString("Installed Personas:\n")
	for _, o := range opts {
		p, _ := mng.GetPersona(o)
		content.WriteString(fmt.Sprintf("\n- %s (%s) => %s / %s", p.Name, p.Trope, p.Provider, p.VoiceID))
	}
	fmt.Println(listStyle.Render(content.String()))
}

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

	fileName := strings.ToLower(strings.ReplaceAll(p.Name, " ", "_")) + ".json"
	path := filepath.Join(mng.PersonasDir, fileName)

	if err := os.Remove(path); err != nil {
		PrintError(fmt.Sprintf("Error deleting file: %v", err))
	} else {
		delete(mng.Personas, p.Name)
		PrintSuccess(fmt.Sprintf("Successfully deleted %s", path))
	}
}

func createPersona(mng *personas.Manager, existing *personas.Persona) {
	var p personas.Persona
	if existing != nil {
		p = *existing
	} else {
		p = personas.Persona{}
	}

	providerOptions := []huh.Option[string]{
		huh.NewOption("OpenAI", "openai_tts"),
		huh.NewOption("ElevenLabs", "elevenlabs_tts"),
		huh.NewOption("Fish Audio", "fishaudio_tts"),
		huh.NewOption("Cartesia", "cartesia_tts"),
		huh.NewOption("Neets", "neets_tts"),
		huh.NewOption("PlayHT", "playht_tts"),
		huh.NewOption("Azure Speech", "azure_tts"),
		huh.NewOption("Local REST API", "local_tts"),
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Persona Name (e.g., GLaDOS)").Value(&p.Name),
			huh.NewInput().Title("Vocal Trope Context (e.g., Sarcastic AI)").Value(&p.Trope),
			huh.NewSelect[string]().Title("Backing Provider").Options(providerOptions...).Value(&p.Provider),
			huh.NewInput().Title("Explicit Voice UUID/Hex").Value(&p.VoiceID),
		),
	).WithTheme(OneDarkTheme()).Run()

	if err != nil {
		PrintWarning("Persona configuration cancelled.")
		return
	}

	if p.Name == "" || p.Provider == "" || p.VoiceID == "" {
		PrintError("Name, Provider, and VoiceID cannot be empty.")
		return
	}

	err = mng.SavePersona(p)
	if err != nil {
		PrintError(fmt.Sprintf("Error saving persona: %v", err))
		return
	}

	var finishStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#98c379")).
		Padding(1, 4).
		MarginTop(1)

	var b strings.Builder
	b.WriteString(icon("\uf00c ", "") + fmt.Sprintf("Persona '%s' saved successfully!\n\n", p.Name))
	b.WriteString(fmt.Sprintf("{\n  \"name\": \"%s\",\n  \"trope\": \"%s\",\n  \"provider\": \"%s\",\n  \"voice_id\": \"%s\"\n}", p.Name, p.Trope, p.Provider, p.VoiceID))

	fmt.Println(finishStyle.Render(b.String()))
}
