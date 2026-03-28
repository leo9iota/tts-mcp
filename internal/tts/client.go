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

// GenerateSpeech takes conversational text and a Voice ID and downloads the audio from Fish Audio API
func GenerateSpeech(text string, voiceID string) error {
	apiKey := os.Getenv("FISH_AUDIO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("FISH_AUDIO_API_KEY environment variable is not set")
	}

	reqBody := Request{
		Text:        text,
		ReferenceID: voiceID,
		Format:      "wav",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.fish.audio/v1/tts", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("model", "s2-pro") 

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	outFile, err := os.Create("temp.wav")
	if err != nil {
		return fmt.Errorf("failed to create output file temp.wav: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save audio stream: %v", err)
	}

	return nil
}
