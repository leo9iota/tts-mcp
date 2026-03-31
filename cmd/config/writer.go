package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteEnvFile(cfg *SetupConfig) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dataDir := filepath.Join(cwd, "data")
	_ = os.MkdirAll(dataDir, 0o755)

	envPath := filepath.Join(dataDir, ".env")

	// 1. Parse existing .env explicitly retaining mapping
	envMap := make(map[string]string)
	if data, err := os.ReadFile(envPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}

		_ = os.WriteFile(envPath+".bak", data, 0644)
	}

	switch cfg.SelectedProvider {
	case "OPENAI":
		envMap["OPENAI_API_KEY"] = cfg.OpenAIKey
	case "ELEVENLABS":
		envMap["ELEVENLABS_API_KEY"] = cfg.ElevenLabsKey
	case "FISH_AUDIO":
		envMap["FISH_AUDIO_API_KEY"] = cfg.FishAudioKey
	case "CARTESIA":
		envMap["CARTESIA_API_KEY"] = cfg.CartesiaKey
	case "NEETS":
		envMap["NEETS_API_KEY"] = cfg.NeetsKey
	case "PLAYHT":
		envMap["PLAYHT_API_KEY"] = cfg.PlayHTKey
		envMap["PLAYHT_USER_ID"] = cfg.PlayHTUser
	case "AZURE":
		envMap["AZURE_SPEECH_KEY"] = cfg.AzureKey
		envMap["AZURE_SPEECH_REGION"] = cfg.AzureRegion
	case "LOCAL":
		envMap["LOCAL_TTS_ENDPOINT"] = cfg.LocalEndpoint
	}

	f, err := os.Create(envPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range envMap {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}

	return nil
}
