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

var _ Provider = (*ElevenLabsProvider)(nil)

type ElevenLabsSynthesisOptions struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

type ElevenLabsRequest struct {
	Text          string                     `json:"text"`
	ModelID       string                     `json:"model_id"`
	VoiceSettings ElevenLabsSynthesisOptions `json:"voice_settings"`
}

type ElevenLabsProvider struct{}

func NewElevenLabsProvider() *ElevenLabsProvider {
	return &ElevenLabsProvider{}
}

func (p *ElevenLabsProvider) ToolName() string {
	return "elevenlabs_tts"
}

func (p *ElevenLabsProvider) Description() string {
	return "Generates incredibly human, cinematic-grade speech via ElevenLabs."
}

func (p *ElevenLabsProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to convert to speech.")),
		mcp.WithString("voice_id", mcp.Description("The ElevenLabs voice ID to use. Falls back to a default if empty.")),
	}
}

func (p *ElevenLabsProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ELEVENLABS_API_KEY is not set")
	}

	if voiceID == "" {
		voiceID = "1SM7GgM6IMuvQlz2BwM3" // Inherited default from Blacktop
	}

	modelID := os.Getenv("ELEVENLABS_MODEL_ID")
	if modelID == "" {
		modelID = "eleven_v3"
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream", voiceID)
	reqBody := ElevenLabsRequest{
		Text:    text,
		ModelID: modelID,
		VoiceSettings: ElevenLabsSynthesisOptions{
			Stability:       0.5,
			SimilarityBoost: 0.75,
			Style:           0.5,
			UseSpeakerBoost: false,
		},
	}

	if opts, ok := ctx.Value(OptionsKey).(map[string]interface{}); ok {
		if stability, ok := opts["stability"].(float64); ok {
			reqBody.VoiceSettings.Stability = stability
		}
		if similarity, ok := opts["similarity_boost"].(float64); ok {
			reqBody.VoiceSettings.SimilarityBoost = similarity
		}
		if style, ok := opts["style"].(float64); ok {
			reqBody.VoiceSettings.Style = style
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

	req.Header.Set("xi-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "audio/mpeg")

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
