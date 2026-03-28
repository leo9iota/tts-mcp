package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Try to load .env file if it exists
	_ = godotenv.Load()

	// Initialize the MCP server
	s := server.NewMCPServer(
		"tts-mcp",
		"1.0.0",
		// We avoid verbose logging by default when using stdio to prevent polluting the transport
	)

	// Register the tool with the AI
	tool := mcp.NewTool("generate_speech",
		mcp.WithDescription("Takes conversational text and a specific character voice ID, generates the audio, and plays it out loud on the host machine."),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The exact phrase the AI wants to say."),
		),
		mcp.WithString("voice_id",
			mcp.Required(),
			mcp.Description("The ID of the character model to use (e.g., Fish Audio model ID)."),
		),
	)

	// Link the handler
	s.AddTool(tool, generateSpeechHandler)

	// Start serving over STDIO
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}

func generateSpeechHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Arguments missing or invalid format"), nil
	}

	text, ok := args["text"].(string)
	if !ok {
		return mcp.NewToolResultError("Arguments missing or invalid: 'text' must be a string"), nil
	}

	voiceID, ok := args["voice_id"].(string)
	if !ok {
		return mcp.NewToolResultError("Arguments missing or invalid: 'voice_id' must be a string"), nil
	}

	// Phase 1: Call the TTS Client
	err := GenerateSpeech(text, voiceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("TTS API generation failed: %v", err)), nil
	}

	// Phase 2: Play the localized audio file
	err = PlayAudio()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Local audio playback execution failed: %v", err)), nil
	}

	// Return success cleanly
	return mcp.NewToolResultText("Successfully generated and played speech aloud to the user."), nil
}
