package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"tts-mcp/internal/personas"
)

func createPersona(mng *personas.Manager, existing *personas.Persona) {
	var p personas.Persona
	if existing != nil {
		p = *existing
	} else {
		p = personas.Persona{}
	}

	var volumeStr string = "1.0"
	if p.Options != nil {
		if vol, ok := p.Options["volume"]; ok {
			if volFloat, ok := vol.(float64); ok {
				volumeStr = strconv.FormatFloat(volFloat, 'f', -1, 64)
			} else if volFloat, ok := vol.(float32); ok {
				volumeStr = strconv.FormatFloat(float64(volFloat), 'f', -1, 64)
			}
		}
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
			huh.NewInput().Title("Persona Name (example: GLaDOS)").Value(&p.Name),
			huh.NewInput().Title("Vocal Trope Context (example: Sarcastic AI)").Value(&p.Trope),
			huh.NewSelect[string]().Title("Backing Provider").Options(providerOptions...).Value(&p.Provider),
			huh.NewInput().Title("Explicit Voice UUID/Hex").Value(&p.VoiceID),
			huh.NewInput().Title("Volume Multiplier (1.0 = Default, 0.5 = 50%)").Value(&volumeStr),
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

	if volumeStr != "" {
		if vol, err := strconv.ParseFloat(volumeStr, 64); err == nil {
			if p.Options == nil {
				p.Options = make(map[string]interface{})
			}
			p.Options["volume"] = vol
		}
	}

	err = mng.SavePersona(p)
	if err != nil {
		PrintError(fmt.Sprintf("Error saving persona: %v", err))
		return
	}

	finishStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#98c379")).
		Padding(1, 4).
		MarginTop(1)

	var b strings.Builder
	b.WriteString(icon("\uf00c ", "") + fmt.Sprintf("Persona '%s' saved successfully!\n\n", p.Name))
	b.WriteString(fmt.Sprintf("name = \"%s\"\ntrope = \"%s\"\nprovider = \"%s\"\nvoice_id = \"%s\"", p.Name, p.Trope, p.Provider, p.VoiceID))

	if volumeStr != "" && volumeStr != "1.0" {
		b.WriteString(fmt.Sprintf("\nvolume = %s", volumeStr))
	}

	fmt.Println(finishStyle.Render(b.String()))
}
