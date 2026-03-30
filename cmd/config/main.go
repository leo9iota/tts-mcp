package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var noIcons bool

// icon strictly returns the Nerd Font icon if not disabled via flag, otherwise yields the ASCII fallback.
func icon(nerd, fallback string) string {
	if noIcons {
		return fallback
	}
	return nerd
}

type SetupConfig struct {
	SelectedProvider string

	OpenAIKey     string
	ElevenLabsKey string
	FishAudioKey  string
	CartesiaKey   string
	NeetsKey      string
	PlayHTKey     string
	PlayHTUser    string
	AzureKey      string
	AzureRegion   string
	LocalEndpoint string
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

func main() {
	flag.BoolVar(&noIcons, "no-icons", false, "Disable Nerd Font icons in the terminal UI")
	flag.Parse()

	var headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#61afef")).
		MarginBottom(1)

	fmt.Println(headerStyle.Render(icon("\uf085 ", "") + "Welcome to the tts-mcp Configuration Wizard!"))

	cfg := &SetupConfig{}

	// 1. Ask which providers
	providerOptions := []huh.Option[string]{
		huh.NewOption("OpenAI", "OPENAI"),
		huh.NewOption("ElevenLabs", "ELEVENLABS"),
		huh.NewOption("Fish Audio", "FISH_AUDIO"),
		huh.NewOption("Cartesia", "CARTESIA"),
		huh.NewOption("Neets", "NEETS"),
		huh.NewOption("PlayHT", "PLAYHT"),
		huh.NewOption("Azure Speech", "AZURE"),
		huh.NewOption("Local REST API", "LOCAL"),
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which TTS Provider do you wish to configure?").
				Description("Use ↑/↓ navigation and press Enter to select").
				Options(providerOptions...).
				Value(&cfg.SelectedProvider),
		),
	).WithTheme(OneDarkTheme()).Run()

	if err != nil {
		fmt.Printf("Setup cancelled: %v\n", err)
		os.Exit(1)
	}

	if cfg.SelectedProvider == "" {
		fmt.Println("No provider selected. Exiting.")
		os.Exit(0)
	}

	// 2. Dynamically build the input form based on selections
	var groups []*huh.Group

	switch cfg.SelectedProvider {
	case "OPENAI":
		groups = append(groups, huh.NewGroup(
			huh.NewInput().Title("OpenAI API Key (sk-...)").EchoMode(huh.EchoModePassword).Value(&cfg.OpenAIKey),
		))
		case "ELEVENLABS":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("ElevenLabs API Key").EchoMode(huh.EchoModePassword).Value(&cfg.ElevenLabsKey),
			))
		case "FISH_AUDIO":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("Fish Audio API Key").EchoMode(huh.EchoModePassword).Value(&cfg.FishAudioKey),
			))
		case "CARTESIA":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("Cartesia API Key").EchoMode(huh.EchoModePassword).Value(&cfg.CartesiaKey),
			))
		case "NEETS":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("Neets API Key").EchoMode(huh.EchoModePassword).Value(&cfg.NeetsKey),
			))
		case "PLAYHT":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("PlayHT API Key").EchoMode(huh.EchoModePassword).Value(&cfg.PlayHTKey),
				huh.NewInput().Title("PlayHT User ID").EchoMode(huh.EchoModePassword).Value(&cfg.PlayHTUser),
			))
		case "AZURE":
			groups = append(groups, huh.NewGroup(
				huh.NewInput().Title("Azure Speech Key").EchoMode(huh.EchoModePassword).Value(&cfg.AzureKey),
				huh.NewInput().Title("Azure Region (e.g. eastus)").Value(&cfg.AzureRegion),
			))
		case "LOCAL":
		groups = append(groups, huh.NewGroup(
			huh.NewInput().Title("Local TTS Endpoint URL").Value(&cfg.LocalEndpoint),
		))
	}

	err = huh.NewForm(groups...).WithTheme(OneDarkTheme()).Run()
	if err != nil {
		fmt.Printf("Input cancelled: %v\n", err)
		os.Exit(1)
	}

	// 3. Validation
	fmt.Printf("\n%sValidating keys...\n", icon("\uf002 ", ""))
	if err := ValidateConfig(cfg); err != nil {
		fmt.Printf("%sValidation Failed: %v\n", icon("\uf00d ", "Error: "), err)
		os.Exit(1)
	}

	// 4. Output Generation
	err = WriteEnvFile(cfg)
	if err != nil {
		fmt.Printf("%sFailed to write .env: %v\n", icon("\uf00d ", "Error: "), err)
		os.Exit(1)
	}

	var finishStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#98c379")).
		MarginTop(1)
		
	fmt.Println(finishStyle.Render(icon("\uf00c ", "") + "Configuration Saved Successfully!"))
	
	exePath, _ := os.Executable()
	dir := filepath.Dir(exePath)
	fmt.Printf("\nYour .env has been generated at %s/.env\n", filepath.Dir(dir))
	fmt.Println("You can now connect the server. Example MCP command array:")
	fmt.Println("[\"go\", \"run\", \".\"]")
}
