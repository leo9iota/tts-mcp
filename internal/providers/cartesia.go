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

var _ Provider = (*CartesiaProvider)(nil)

type CartesiaVoice struct {
	Mode string `json:"mode"`
	ID   string `json:"id"`
}

type CartesiaOutputFormat struct {
	Container  string `json:"container"`
	Encoding   string `json:"encoding"`
	SampleRate int    `json:"sample_rate"`
}

type CartesiaRequest struct {
	ModelID      string               `json:"model_id"`
	Transcript   string               `json:"transcript"`
	Voice        CartesiaVoice        `json:"voice"`
	OutputFormat CartesiaOutputFormat `json:"output_format"`
}

type CartesiaProvider struct{}

func NewCartesiaProvider() *CartesiaProvider {
	return &CartesiaProvider{}
}

func (p *CartesiaProvider) ToolName() string {
	return "cartesia_tts"
}

func (p *CartesiaProvider) Description() string {
	return "Generate ultra-low latency, expressive speech using Cartesia Sonic."
}

func (p *CartesiaProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to convert to speech.")),
		mcp.WithString("voice_id", mcp.Description("The Cartesia voice ID to use. Falls back to a default if empty.")),
	}
}

func (p *CartesiaProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("CARTESIA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CARTESIA_API_KEY is not set")
	}

	if voiceID == "" {
		voiceID = "a0e99841-438c-4a64-b6a9-62f1060dcb11" // Default English voice
	}

	modelID := os.Getenv("CARTESIA_MODEL_ID")
	if modelID == "" {
		modelID = "sonic-english"
	}

	url := "https://api.cartesia.ai/tts/bytes"
	reqBody := CartesiaRequest{
		ModelID:    modelID,
		Transcript: text,
		Voice: CartesiaVoice{
			Mode: "id",
			ID:   voiceID,
		},
		OutputFormat: CartesiaOutputFormat{
			Container:  "mp3",
			Encoding:   "mp3",
			SampleRate: 44100,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Cartesia-Version", "2024-06-10")
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
