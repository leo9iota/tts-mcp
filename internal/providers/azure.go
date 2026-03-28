package providers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

var _ Provider = (*AzureProvider)(nil)

type AzureProvider struct{}

func NewAzureProvider() *AzureProvider {
	return &AzureProvider{}
}

func (p *AzureProvider) ToolName() string {
	return "azure_tts"
}

func (p *AzureProvider) Description() string {
	return "Generates highly reliable, enterprise-grade neural speech via Microsoft Azure Cognitive Services."
}

func (p *AzureProvider) ToolArguments() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("text", mcp.Required(), mcp.Description("The text to convert to speech.")),
		mcp.WithString("voice_id", mcp.Description("The specific Azure Neural voice ID (e.g. en-US-ChristopherNeural).")),
	}
}

func (p *AzureProvider) StreamSpeech(ctx context.Context, text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("AZURE_SPEECH_KEY")
	region := os.Getenv("AZURE_SPEECH_REGION")

	if apiKey == "" || region == "" {
		return nil, fmt.Errorf("AZURE_SPEECH_KEY and AZURE_SPEECH_REGION must be set")
	}

	if voiceID == "" {
		voiceID = "en-US-ChristopherNeural" // Safe standard male voice
	}

	// Azure requires Speech Synthesis Markup Language (SSML)
	ssml := fmt.Sprintf(`
		<speak version='1.0' xml:lang='en-US'>
			<voice name='%s'>%s</voice>
		</speak>
	`, voiceID, text)

	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", region)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(ssml))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-24khz-48kbitrate-mono-mp3")
	req.Header.Set("User-Agent", "tts-mcp")

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
