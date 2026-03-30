package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

var _ Provider = (*NeetsProvider)(nil)

type NeetsRequest struct {
	Text    string  `json:"text"`
	VoiceID string  `json:"voice_id"`
	Fmt     string  `json:"fmt"`
	Speed   float64 `json:"speed,omitempty"`
}

type NeetsProvider struct{}

func NewNeetsProvider() *NeetsProvider {
	return &NeetsProvider{}
}

func (p *NeetsProvider) ToolName() string {
	return "neets_tts"
}

func (p *NeetsProvider) Description() string {
	return "Generates ultra-cheap, highly realistic TTS using the Neets.ai REST API."
}

func (p *NeetsProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to synthesize.")),
		mcp.WithString("voice_id", mcp.Description("The explicit Neets.ai reference_id (e.g., 'us-male-1'). Default is 'us-male-1' if blank.")),
	}
}

func (p *NeetsProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("NEETS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("NEETS_API_KEY environment variable is not set")
	}

	if voiceID == "" {
		voiceID = "us-male-1" // Fallback fallback to a generic voice if empty
	}

	reqBody := NeetsRequest{
		Text:    text,
		VoiceID: voiceID,
		Fmt:     "mp3",
		Speed:   1.0,
	}

	if opts, ok := ctx.Value(OptionsKey).(map[string]interface{}); ok {
		if speed, ok := opts["speed"].(float64); ok {
			reqBody.Speed = speed
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.neets.ai/v1/tts", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp.Body, nil
}
