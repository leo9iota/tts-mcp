package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Request struct {
	Text        string `json:"text"`
	ReferenceID string `json:"reference_id"`
	Format      string `json:"format"`
}

// StreamSpeech constructs the Fish Audio conversational TTS JSON and returns the explicit MP3 raw response body.
func StreamSpeech(text string, voiceID string) (io.ReadCloser, error) {
	apiKey := os.Getenv("FISH_AUDIO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FISH_AUDIO_API_KEY environment variable is not set")
	}

	reqBody := Request{
		Text:        text,
		ReferenceID: voiceID,
		Format:      "mp3", // Use robust mp3 format so we can stream it dynamically over beep!
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.fish.audio/v1/tts", bytes.NewBuffer(jsonData))
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
