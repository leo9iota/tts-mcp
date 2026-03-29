package api

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gopxl/beep/v2/mp3"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"tts-mcp/internal/audio"
	"tts-mcp/internal/personas"
	"tts-mcp/internal/providers"
)

// Start initializes the toolsets and serves the MCP stdio handler
func Start() {
	s := server.NewMCPServer("tts-mcp", "1.0.0")

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
		s.AddTool(tool, createHandler(s, p))
	}

	// Register unified Persona Tool if any exist
	if len(personaManager.GetOptions()) > 0 && personaManager.GetOptions()[0] != "" {
		personaTool := mcp.NewTool("speak_as_persona",
			mcp.WithDescription("Summon a specific character persona from the locally configured data/ directories. This abstracts the TTS backend and dynamically binds voices for seamless testing."),
			mcp.WithString("persona", mcp.Required(), mcp.Description("The loaded character persona to invoke."), mcp.Enum(personaManager.GetOptions()...)),
			mcp.WithString("text", mcp.Required(), mcp.Description("The text for the character to say.")),
		)
		s.AddTool(personaTool, createPersonaHandler(s, personaManager, providerList))
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}

func createHandler(s *server.MCPServer, provider providers.Provider) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid format"), nil
		}

		progressToken := request.Params.Meta.ProgressToken

		text, ok := args["text"].(string)
		if !ok {
			return mcp.NewToolResultError("Arguments missing or invalid: 'text' must be a string"), nil
		}

		var voiceID string
		if vid, ok := args["voice_id"].(string); ok {
			voiceID = vid
		}

		// 1. Connect the read closer explicitly capturing the HTTP payload as it downloads
		respBody, err := provider.StreamSpeech(ctx, text, voiceID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%s streaming failed: %v", provider.ToolName(), err)), nil
		}

		// 2. Clone the stream: Pass one to local file, pass one to hardware speaker pipe
		file, err := os.Create("temp.mp3")
		if err != nil {
			respBody.Close()
			return mcp.NewToolResultError(fmt.Sprintf("Failed to construct audio disk file: %v", err)), nil
		}
		absPath, _ := filepath.Abs("temp.mp3")

		pipeReader, pipeWriter := io.Pipe()

		go func() {
			defer pipeWriter.Close()
			defer respBody.Close()
			defer file.Close()

			tee := io.TeeReader(respBody, file)
			_, copyErr := io.Copy(pipeWriter, tee)
			if copyErr != nil {
				pipeWriter.CloseWithError(copyErr)
			}
		}()

		// 3. Mount pipe locally within beep audio decoder
		streamer, format, err := mp3.Decode(pipeReader)
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
			audioComplete <- audio.WaitAndPlay(streamer, format.SampleRate, reporter)
		}()

		// 4. Thread-lock the active response on the active OS block waiting for ctx.Done internally!
		select {
		case err := <-audioComplete:
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Audio engine execution failed: %v", err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Successfully generated natively and played speech aloud to the user using %s!\nSaved localized version at: %s", provider.ToolName(), absPath)), nil

		case <-ctx.Done():
			audio.Stop()
			return mcp.NewToolResultError("Audio generation forcefully cancelled by user context!"), nil
		}
	}
}

func createPersonaHandler(s *server.MCPServer, mng *personas.Manager, providerList []providers.Provider) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		request.Params.Arguments = map[string]interface{}{
			"text":     args["text"],
			"voice_id": persona.VoiceID,
		}

		// Delegate directly to the standard handler!
		return createHandler(s, targetProvider)(ctx, request)
	}
}
