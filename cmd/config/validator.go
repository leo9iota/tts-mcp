package main

import (
	"fmt"
	"net/http"
	"time"
)

func ValidateConfig(cfg *SetupConfig) error {
	client := &http.Client{Timeout: 5 * time.Second}

	switch cfg.SelectedProvider {
	case "OPENAI":
		if cfg.OpenAIKey == "" {
			return fmt.Errorf("OpenAI API Key is empty")
		}
		req, _ := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
		req.Header.Add("Authorization", "Bearer "+cfg.OpenAIKey)
		if err := performPing(client, req, "OpenAI"); err != nil {
			return err
		}

	case "ELEVENLABS":
		if cfg.ElevenLabsKey == "" {
			return fmt.Errorf("ElevenLabs API Key is empty")
		}
		req, _ := http.NewRequest("GET", "https://api.elevenlabs.io/v1/voices", nil)
		req.Header.Add("xi-api-key", cfg.ElevenLabsKey)
		if err := performPing(client, req, "ElevenLabs"); err != nil {
			return err
		}

	case "FISH_AUDIO":
		if cfg.FishAudioKey == "" {
			return fmt.Errorf("Fish Audio API Key is empty")
		}
		req, _ := http.NewRequest("GET", "https://api.fish.audio/v1/models", nil)
		req.Header.Add("Authorization", "Bearer "+cfg.FishAudioKey)
		if err := performPing(client, req, "Fish Audio"); err != nil {
			return err
		}

	case "CARTESIA":
		if cfg.CartesiaKey == "" {
			return fmt.Errorf("Cartesia API Key is empty")
		}
		req, _ := http.NewRequest("GET", "https://api.cartesia.ai/voices", nil)
		req.Header.Add("X-API-Key", cfg.CartesiaKey)
		req.Header.Add("Cartesia-Version", "2024-06-10")
		if err := performPing(client, req, "Cartesia"); err != nil {
			return err
		}

	case "PLAYHT":
		if cfg.PlayHTKey == "" || cfg.PlayHTUser == "" {
			return fmt.Errorf("PlayHT Key or User ID is empty")
		}
		req, _ := http.NewRequest("GET", "https://api.play.ht/api/v2/voices", nil)
		req.Header.Add("Authorization", "Bearer "+cfg.PlayHTKey)
		req.Header.Add("X-User-Id", cfg.PlayHTUser)
		if err := performPing(client, req, "PlayHT"); err != nil {
			return err
		}

	case "AZURE":
		if cfg.AzureKey == "" || cfg.AzureRegion == "" {
			return fmt.Errorf("Azure Key or Region is empty")
		}
		url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list", cfg.AzureRegion)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Ocp-Apim-Subscription-Key", cfg.AzureKey)
		if err := performPing(client, req, "Azure"); err != nil {
			return err
		}

	case "NEETS":
		if cfg.NeetsKey == "" {
			return fmt.Errorf("Neets API Key is empty")
		}

	case "LOCAL":
		if cfg.LocalEndpoint == "" {
			return fmt.Errorf("Local Endpoint is empty")
		}
	}
	return nil
}

func performPing(client *http.Client, req *http.Request, name string) error {
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s network error: %v", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("%s API Key is invalid (HTTP %d)", name, resp.StatusCode)
	}
	return nil
}
