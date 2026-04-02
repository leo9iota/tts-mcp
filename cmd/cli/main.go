package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"os"

	"tts-mcp/internal/config"
)

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

func main() {
	flag.BoolVar(&noIcons, "no-icons", false, "Disable Nerd Font icons in the terminal UI")
	flag.Parse()

	var headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#61afef")).
		MarginBottom(1)

	for {
		fmt.Println(headerStyle.Render(icon("\uf130 ", "") + "Welcome to the tts-mcp Configuration Wizard!"))

		var action string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Main Menu").
					Description("Select a configuration module").
					Options(
						huh.NewOption("Configure API Keys", "ENV"),
						huh.NewOption("Manage Personas", "PERSONA"),
						huh.NewOption("Exit Wizard", "EXIT"),
					).
					Value(&action),
			),
		).WithTheme(OneDarkTheme()).Run()

		if err != nil || action == "EXIT" {
			PrintInfo("Exiting wizard. Goodbye!")
			os.Exit(0)
		}

		switch action {
		case "ENV":
			RunEnvWizard()
		case "PERSONA":
			RunPersonaWizard()
		}
	}
}

func RunEnvWizard() {

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
		PrintWarning(fmt.Sprintf("Setup cancelled: %v", err))
		os.Exit(1)
	}

	if cfg.SelectedProvider == "" {
		PrintWarning("No provider selected. Exiting.")
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
		PrintWarning(fmt.Sprintf("Input cancelled: %v", err))
		os.Exit(1)
	}

	// 3. Validation
	PrintInfo("Validating API keys with external network...")
	if err := ValidateConfig(cfg); err != nil {
		PrintError(fmt.Sprintf("Validation Failed: %v", err))
		os.Exit(1)
	}

	// 4. Output Generation
	err = WriteEnvFile(cfg)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to write .env storage: %v", err))
		os.Exit(1)
	}

	var infoStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#98c379")).
		Padding(1, 4).
		MarginTop(1).
		MarginBottom(1)

	content := icon("\uf00c ", "") + "Configuration Saved Successfully!\n\n"
	content += fmt.Sprintf("Your environment configuration was automatically localized inside your global OS configuration directory:\n%s\n\n", config.GetEnvPath())
	content += "You can connect this MCP server via standard stdio:\n"
	content += "[\"go\", \"run\", \".\"]"

	fmt.Println(infoStyle.Render(content))
}
