package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"tts-mcp/internal/audio"
	"tts-mcp/internal/output"
	"tts-mcp/internal/personas"
	"tts-mcp/internal/providers"
)

// Start initializes the toolsets and serves the MCP stdio handler
func Start() {
	s := server.NewMCPServer("tts-mcp", "1.0.0")

	// Phase 2: Scaffold decoupled engine instance per server session
	audioEngine := audio.NewEngine()

	// 1. Load available Personas
	personaManager, _ := personas.NewManager()

	var providerList []providers.Provider

	if os.Getenv("FISH_AUDIO_API_KEY") != "" {
		providerList = append(providerList, providers.NewFishAudioProvider())
	}
	if os.Getenv("NEETS_API_KEY") != "" {
		providerList = append(providerList, providers.NewNeetsProvider())
	}
	if os.Getenv("ELEVENLABS_API_KEY") != "" {
		providerList = append(providerList, providers.NewElevenLabsProvider())
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		providerList = append(providerList, providers.NewOpenAIProvider())
	}
	if os.Getenv("CARTESIA_API_KEY") != "" {
		providerList = append(providerList, providers.NewCartesiaProvider())
	}
	if os.Getenv("PLAYHT_API_KEY") != "" && os.Getenv("PLAYHT_USER_ID") != "" {
		providerList = append(providerList, providers.NewPlayHTProvider())
	}
	if os.Getenv("AZURE_SPEECH_KEY") != "" && os.Getenv("AZURE_SPEECH_REGION") != "" {
		providerList = append(providerList, providers.NewAzureProvider())
	}
	if os.Getenv("LOCAL_TTS_ENDPOINT") != "" {
		providerList = append(providerList, providers.NewLocalProvider())
	}

	for _, p := range providerList {
		opts := []mcp.ToolOption{
			mcp.WithDescription(p.Description()),
		}
		opts = append(opts, p.ToolArguments()...)

		tool := mcp.NewTool(p.ToolName(), opts...)
		s.AddTool(tool, WithRecovery(createHandler(s, p, audioEngine)))
	}

	// Register unified Persona Tool if any exist
	if len(personaManager.GetOptions()) > 0 && personaManager.GetOptions()[0] != "" {
		personaTool := mcp.NewTool("speak_as_persona",
			mcp.WithDescription("Invoke a specific character persona from the locally configured directories. Abstracts the TTS backend and dynamically binds voices for seamless testing."),
			mcp.WithString("persona", mcp.Required(), mcp.Description("The loaded character persona to invoke."), mcp.Enum(personaManager.GetOptions()...)),
			mcp.WithString("text", mcp.Required(), mcp.Description("The text for the character to say.")),
			mcp.WithNumber("volume", mcp.Description("Optional volume multiplier (e.g. 0.5 for 50%, 2.0 for 200%). Defaults to 1.0.")),
		)
		s.AddTool(personaTool, WithRecovery(createPersonaHandler(s, personaManager, providerList, audioEngine)))
	}

	// Register IDE Persona Generator Tool
	generatorTool := mcp.NewTool("create_persona",
		mcp.WithDescription("Create and load a new persona configuration bound to a specific TTS provider and voice ID. Registers the configuration into the file system for persistent availability."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The precise single-word formal name of the character (e.g., 'Megumin').")),
		mcp.WithString("trope", mcp.Required(), mcp.Description("A brief semantic description of their personality or vocal trope (e.g., 'gruff medieval narrator').")),
		mcp.WithString("provider", mcp.Required(), mcp.Description("The exact ToolName of the underlying TTS provider (e.g., 'fishaudio_tts', 'elevenlabs_tts').")),
		mcp.WithString("voice_id", mcp.Required(), mcp.Description("The exact hex or UUID string natively mapping to the provider's specific voice model.")),
	)
	s.AddTool(generatorTool, WithRecovery(createPersonaGeneratorHandler(s, personaManager)))

	log.Info("TTS-MCP Server Initialized \nWaiting for Antigravity IDE JSON-RPC connections via stdio...",
		"personas_loaded", len(personaManager.Personas),
		"providers_active", len(providerList),
	)

	if err := server.ServeStdio(s); err != nil {
		log.Error("Server error", "err", err)
	}
}

func createHandler(s *server.MCPServer, provider providers.Provider, audioEngine *audio.AudioEngine) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid format"), nil
		}

		var progressToken interface{}
		if request.Params.Meta != nil {
			progressToken = request.Params.Meta.ProgressToken
		}

		text, ok := args["text"].(string)
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid: 'text' must be a string"), nil
		}

		var voiceID string
		if vid, ok := args["voice_id"].(string); ok {
			voiceID = vid
		}

		var volume float64 = 1.0
		if vol, ok := args["volume"].(float64); ok {
			volume = vol
		} else if volInt, ok := args["volume"].(int); ok {
			volume = float64(volInt)
		} else if volInt, ok := args["volume"].(int32); ok {
			volume = float64(volInt)
		} else if volInt, ok := args["volume"].(int64); ok {
			volume = float64(volInt)
		}

		// Inject the entire Argument footprint into Context so polymorphic providers can extract custom mappings
		ctx = context.WithValue(ctx, providers.OptionsKey, args)

		// 1. Connect the read closer explicitly capturing the HTTP payload as it downloads
		respBody, err := provider.StreamSpeech(ctx, text, voiceID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%s streaming failed: %v", provider.ToolName(), err)), nil
		}

		defer respBody.Close()

		// 2. Read full network payload into memory buffer to guarantee jitter-free beep playback
		audioBytes, err := io.ReadAll(respBody)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to download audio stream: %v", err)), nil
		}

		// 3. Write memory slice directly to disk artifact independent of hardware stream
		personaName, _ := args["persona"].(string) // Safe to cast implicitly fallback to "" if missing
		file, absPath, err := output.GenerateOutputFile(personaName, provider.ToolName())
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to construct history audio file: %v", err)), nil
		}
		if _, err := file.Write(audioBytes); err != nil {
			file.Close()
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write to local history file: %v", err)), nil
		}
		file.Close()

		// 4. Mount fully loaded memory buffer locally into beep audio decoder
		audioReader := io.NopCloser(bytes.NewReader(audioBytes))
		streamer, format, err := mp3.Decode(audioReader)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to decode mp3: %v", err)), nil
		}
		defer streamer.Close()
		// 4. Construct lightweight telemetry closure
		var reporter func(pos int, total int, message string) = nil
		lastPercent := -1
		if progressToken != nil && progressToken != "" {
			reporter = func(pos int, total int, message string) {
				percent := 0
				if total > 0 {
					percent = int((float64(pos) / float64(total)) * 100)
				}

				if percent != lastPercent || total <= 0 {
					lastPercent = percent
					progFloat := float64(pos)
					if total <= 0 {
						progFloat = 0 // indeterminate
					}

					// Send dynamic RPC boundary over HTTP map structure
					s.SendNotificationToAllClients("notifications/progress", map[string]interface{}{
						"progressToken": progressToken,
						"progress":      progFloat,
						"total":         float64(total),
						"message":       message,
					})
				}
			}
		}

		// 5. Stream sequence locking
		audioComplete := make(chan error, 1)
		go func() {
			audioComplete <- audioEngine.WaitAndPlay(ctx, streamer, format.SampleRate, volume, reporter)
		}()

		// 4. Thread-lock the active response on the active OS block waiting for ctx.Done internally!
		select {
		case err := <-audioComplete:
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Audio engine execution failed: %v", err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Successfully generated natively and played speech aloud to the user using %s!\nSaved localized version at: %s", provider.ToolName(), absPath)), nil

		case <-ctx.Done():
			audioEngine.Stop(ctx)
			return mcp.NewToolResultError("Audio generation forcefully cancelled by user context!"), nil
		}
	}
}

func createPersonaHandler(s *server.MCPServer, mng *personas.Manager, providerList []providers.Provider, audioEngine *audio.AudioEngine) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid format"), nil
		}

		personaName, ok := args["persona"].(string)
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid: 'persona' must be a string"), nil
		}

		persona, exists := mng.GetPersona(personaName)
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("Persona '%s' not found in loaded configurations", personaName)), nil
		}

		// Find the matching provider
		var targetProvider providers.Provider
		for _, p := range providerList {
			if p.ToolName() == persona.Provider {
				targetProvider = p
				break
			}
		}

		if targetProvider == nil {
			return mcp.NewToolResultError(fmt.Sprintf("Provider '%s' required by Persona '%s' is not active. Check API Keys in .env", persona.Provider, persona.Name)), nil
		}

		// Map to standard arguments dynamically
		mappedArgs := map[string]interface{}{
			"text":     args["text"],
			"voice_id": persona.VoiceID,
			"persona":  personaName,
		}
		if vol, ok := args["volume"]; ok {
			mappedArgs["volume"] = vol
		}

		// Phase 2: Tool Argument Hydration (FEAT-003)
		// Hydrate arbitrary modulation options securely scaling provider integrations
		if persona.Options != nil {
			for k, v := range persona.Options {
				// Don't overwrite the core structural strings manually set above
				if _, exists := mappedArgs[k]; !exists {
					mappedArgs[k] = v
				}
			}
		}

		request.Params.Arguments = mappedArgs

		// Delegate directly to the standard handler!
		return createHandler(s, targetProvider, audioEngine)(ctx, request)
	}
}

func createPersonaGeneratorHandler(s *server.MCPServer, mng *personas.Manager) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid format"), nil
		}

		name, _ := args["name"].(string)
		trope, _ := args["trope"].(string)
		provider, _ := args["provider"].(string)
		voiceID, _ := args["voice_id"].(string)

		if name == "" || provider == "" || voiceID == "" {
			return mcp.NewToolResultError("Missing required parameters: 'name', 'provider', and 'voice_id' are strictly required."), nil
		}

		p := personas.Persona{
			Name:     name,
			Trope:    trope,
			Provider: provider,
			VoiceID:  voiceID,
		}

		if err := mng.SavePersona(p); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to construct and save persona file: %v", err)), nil
		}

		successMsg := fmt.Sprintf("Successfully generated and hot-loaded new persona: '%s' (%s) bound to %s via Voice ID: %s",
			p.Name, p.Trope, p.Provider, p.VoiceID)

		// Fire an autonomous ghost tool ping updating IDE instances dynamically
		s.SendNotificationToAllClients("notifications/tools/list_changed", nil)

		return mcp.NewToolResultText(successMsg), nil
	}
}
