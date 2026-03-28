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
	"tts-mcp/internal/tts"
)

// Start initializes the toolsets and serves the MCP stdio handler
func Start() {
	s := server.NewMCPServer("tts-mcp", "1.0.0")

	tool := mcp.NewTool("generate_speech",
		mcp.WithDescription("Takes conversational text and a specific character voice ID, generates the audio, and plays it out loud on the host machine natively."),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The exact phrase the AI wants to say."),
		),
		mcp.WithString("voice_id",
			mcp.Required(),
			mcp.Description("The ID of the character model to use (e.g., Fish Audio model ID)."),
		),
	)

	s.AddTool(tool, generateSpeechHandler)

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

	// 1. Connect the read closer explicitly capturing the HTTP payload as it downloads
	respBody, err := tts.StreamSpeech(text, voiceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("TTS API streaming failed: %v", err)), nil
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

	audioComplete := make(chan error, 1)
	go func() {
		audioComplete <- audio.WaitAndPlay(streamer, format.SampleRate)
	}()

	// 4. Thread-lock the active response on the active OS block waiting for ctx.Done internally!
	select {
	case err := <-audioComplete:
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Audio engine execution failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully generated natively and played speech aloud to the user!\nSaved localized version at: %s", absPath)), nil

	case <-ctx.Done():
		audio.Stop()
		return mcp.NewToolResultError("Audio generation forcefully cancelled by user context!"), nil
	}
}
