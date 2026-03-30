package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var noIcons bool

func icon(nerd, fallback string) string {
	if noIcons {
		return fallback
	}
	return nerd
}

func PrintError(msg string) {
	fmt.Println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#e06c75")).Padding(0, 2).Render(icon("\uf00d ", "Error: ") + msg))
}

func PrintSuccess(msg string) {
	fmt.Println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#98c379")).Padding(0, 2).Render(icon("\uf00c ", "Success: ") + msg))
}

func PrintInfo(msg string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#61afef")).Render(icon("\uf05a ", "Info: ") + msg))
}

func PrintWarning(msg string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#e5c07b")).Render(icon("\uf071 ", "Warning: ") + msg))
}

// OneDarkTheme constructs a custom huh.Theme bound natively to the Atom One Dark palette.
func OneDarkTheme() *huh.Theme {
	t := huh.ThemeBase()

	var (
		fg     = lipgloss.AdaptiveColor{Light: "#383a42", Dark: "#abb2bf"}
		blue   = lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"}
		purple = lipgloss.AdaptiveColor{Light: "#a626a4", Dark: "#c678dd"}
		green  = lipgloss.AdaptiveColor{Light: "#50a14f", Dark: "#98c379"}
		red    = lipgloss.AdaptiveColor{Light: "#e45649", Dark: "#e06c75"}
		gray   = lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"}
	)

	t.Focused.Base = t.Focused.Base.Foreground(fg)
	t.Focused.Title = t.Focused.Title.Foreground(blue).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(blue).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(gray)
	t.Focused.SelectSelector = lipgloss.NewStyle().Foreground(purple).SetString("> ")
	t.Focused.Option = lipgloss.NewStyle().Foreground(fg)
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(green)
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(purple).SetString("> ")
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(green).SetString("[x] ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(gray).SetString("[ ] ")
	t.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(blue)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(purple)

	t.Focused.ErrorIndicator = lipgloss.NewStyle().Foreground(red).Bold(true)
	t.Focused.ErrorMessage = lipgloss.NewStyle().Foreground(red)

	t.Blurred.Title = t.Blurred.Title.Foreground(gray).Bold(true)
	t.Blurred.Description = t.Blurred.Description.Foreground(gray)

	return t
}
