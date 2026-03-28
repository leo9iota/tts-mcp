package providers

import (
	"context"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
)

// Provider defines the polymorphic interface for any AI TTS engine.
// All providers must return an io.ReadCloser of exact `.mp3` bytes or standard chunks
// capable of being directly decoded natively by go-beep.
type Provider interface {
	// ToolName returns the unique MCP tool identifier (e.g. "fishaudio_tts").
	ToolName() string

	// Description provides the instructions for the AI on when and how to use this specific engine.
	Description() string

	// ToolArguments returns the exact JSON-Schema parameters expected by this particular model.
	ToolArguments() []mcp.ToolOption

	// StreamSpeech connects to the upstream AI engine and begins piping raw audio bytes.
	// It is critical that this yields the ReadCloser as fast as possible to prevent stream lag.
	StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error)
}
