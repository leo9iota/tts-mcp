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

var _ Provider = (*LocalProvider)(nil)

type LocalRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed,omitempty"`
}

type LocalProvider struct{}

func NewLocalProvider() *LocalProvider {
	return &LocalProvider{}
}

func (p *LocalProvider) ToolName() string {
	return "local_tts"
}

func (p *LocalProvider) Description() string {
	return "Generate speech using a locally hosted or open-source OpenAI-compatible TTS engine."
}

func (p *LocalProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to synthesize locally.")),
		mcp.WithString("voice_id", mcp.Description("The local voice ID / profile to use.")),
	}
}

func (p *LocalProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	endpoint := os.Getenv("LOCAL_TTS_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("LOCAL_TTS_ENDPOINT environment variable is not set")
	}

	apiKey := os.Getenv("LOCAL_TTS_API_KEY") // Optional for most local setups

	if voiceID == "" {
		voiceID = "default"
	}

	modelID := os.Getenv("LOCAL_TTS_MODEL_ID")
	if modelID == "" {
		modelID = "tts-1"
	}

	reqBody := LocalRequest{
		Model:          modelID,
		Input:          text,
		Voice:          voiceID,
		ResponseFormat: "mp3",
		Speed:          1.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

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
