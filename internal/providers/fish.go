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

// ensure interface is implemented at compile time
var _ Provider = (*FishAudioProvider)(nil)

type Request struct {
	Text        string `json:"text"`
	ReferenceID string `json:"reference_id"`
	Format      string `json:"format"`
	Latency     string `json:"latency,omitempty"`
}

type FishAudioProvider struct{}

func NewFishAudioProvider() *FishAudioProvider {
	return &FishAudioProvider{}
}

func (p *FishAudioProvider) ToolName() string {
	return "fishaudio_tts"
}

func (p *FishAudioProvider) Description() string {
	return "Generates anime-style expressive TTS via Fish Audio REST API. Extremely fast, low latency, stylized WAIFU/anime voices."
}

func (p *FishAudioProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to synthesize.")),
		mcp.WithString("voice_id", mcp.Description("The explicit FishAudio reference_id for the model character.")),
	}
}

// StreamSpeech constructs the Fish Audio conversational TTS JSON and exactly streams the mp3 format body output directly.
func (p *FishAudioProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("FISH_AUDIO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FISH_AUDIO_API_KEY environment variable is not set")
	}

	reqBody := Request{
		Text:        text,
		ReferenceID: voiceID,
		Format:      "mp3", // Use robust mp3 format so we can stream it dynamically
		Latency:     "normal",
	}

	if opts, ok := ctx.Value(OptionsKey).(map[string]interface{}); ok {
		// Example: map exact structural latency overrides for the anime API
		if latency, ok := opts["latency"].(string); ok {
			reqBody.Latency = latency
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Make request with context to hook natively into graceful cancellation!
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.fish.audio/v1/tts", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("model", "s2-pro")

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
