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

var _ Provider = (*PlayHTProvider)(nil)

type PlayHTRequest struct {
	Text         string `json:"text"`
	Voice        string `json:"voice"`
	OutputFormat string `json:"output_format"`
	VoiceEngine  string `json:"voice_engine"`
	Speed        float64 `json:"speed,omitempty"`
}

type PlayHTProvider struct{}

func NewPlayHTProvider() *PlayHTProvider {
	return &PlayHTProvider{}
}

func (p *PlayHTProvider) ToolName() string {
	return "playht_tts"
}

func (p *PlayHTProvider) Description() string {
	return "Generates highly expressive and cloneable voice outputs using PlayHT."
}

func (p *PlayHTProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to convert to speech.")),
		mcp.WithString("voice_id", mcp.Description("The PlayHT voice manifest URL to use. Falls back to a default if empty.")),
	}
}

func (p *PlayHTProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("PLAYHT_API_KEY")
	userID := os.Getenv("PLAYHT_USER_ID")

	if apiKey == "" || userID == "" {
		return nil, fmt.Errorf("PLAYHT_API_KEY and PLAYHT_USER_ID must be set")
	}

	if voiceID == "" {
		// Default to a solid standard voice (e.g. female-cs or generic Susan)
		voiceID = "s3://voice-cloning-zero-shot/d9ff78ba-d016-47f6-b0ef-dd630f59414e/female-cs/manifest.json" 
	}

	url := "https://api.play.ht/api/v2/tts/stream"
	reqBody := PlayHTRequest{
		Text:         text,
		Voice:        voiceID,
		OutputFormat: "mp3",
		VoiceEngine:  "PlayHT2.0",
		Speed:        1.0,
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

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("X-User-Id", userID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/mpeg")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp.Body, nil
}
