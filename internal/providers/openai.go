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

var _ Provider = (*OpenAIProvider)(nil)

type OpenAIRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed,omitempty"`
}

type OpenAIProvider struct{}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{}
}

func (p *OpenAIProvider) ToolName() string {
	return "openai_tts"
}

func (p *OpenAIProvider) Description() string {
	return "Generate state-of-the-art conversational speech via the OpenAI TTS API."
}

func (p *OpenAIProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to synthesize.")),
		mcp.WithString("voice_id", mcp.Description("The OpenAI voice ID (e.g., alloy, echo, fable, onyx, nova, shimmer). Default is 'alloy'.")),
	}
}

func (p *OpenAIProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	if voiceID == "" {
		voiceID = "alloy"
	}

	reqBody := OpenAIRequest{
		Model:          "tts-1", // tts-1 is optimized for speed/streaming
		Input:          text,
		Voice:          voiceID,
		ResponseFormat: "mp3",
		Speed:          1.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
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
