package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"tts-mcp/internal/personas"
)

func listPersonas(mng *personas.Manager) {
	opts := mng.GetOptions()
	if len(opts) == 1 && opts[0] == "" {
		PrintInfo("No personas found.")
		return
	}

	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#61afef")).
		Padding(1, 4).MarginTop(1).MarginBottom(1)

	var content strings.Builder
	content.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e5c07b")).Render("Installed Personas:") + "\n\n")

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5c6370")).Width(10)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))

	for _, o := range opts {
		p, _ := mng.GetPersona(o)

		content.WriteString(icon("\uf2bd ", "• ") + nameStyle.Render(p.Name) + "\n")

		if p.Trope != "" {
			content.WriteString("  " + labelStyle.Render("Trope:") + valueStyle.Render(p.Trope) + "\n")
		}

		content.WriteString("  " + labelStyle.Render("Provider:") + valueStyle.Render(p.Provider) + "\n")
		content.WriteString("  " + labelStyle.Render("Voice ID:") + valueStyle.Render(p.VoiceID) + "\n\n")
	}

	fmt.Println(listStyle.Render(strings.TrimSuffix(content.String(), "\n\n")))
}
