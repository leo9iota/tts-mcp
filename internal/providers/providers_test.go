package providers

import (
	"context"
	"testing"
)

func TestProviderInstantiation(t *testing.T) {
	provs := []Provider{
		NewFishAudioProvider(),
		NewNeetsProvider(),
		NewElevenLabsProvider(),
		NewOpenAIProvider(),
		NewCartesiaProvider(),
		NewPlayHTProvider(),
		NewAzureProvider(),
		NewLocalProvider(),
	}

	for _, p := range provs {
		if p.ToolName() == "" {
			t.Errorf("expected provider to have a name")
		}
		if p.Description() == "" {
			t.Errorf("expected provider %s to have a description", p.ToolName())
		}
		if len(p.ToolArguments()) == 0 {
			t.Errorf("expected provider %s to have tool arguments", p.ToolName())
		}
	}
}

func TestLocalProvider_StreamSpeech_NoURL(t *testing.T) {
	t.Setenv("LOCAL_TTS_ENDPOINT", "")
	p := NewLocalProvider()
	_, err := p.StreamSpeech(context.Background(), "hello", "1")
	if err == nil {
		t.Errorf("expected error when LOCAL_TTS_ENDPOINT is empty")
	}
}

func TestFishProvider_StreamSpeech_NoKey(t *testing.T) {
	t.Setenv("FISH_AUDIO_API_KEY", "")
	p := NewFishAudioProvider()
	_, err := p.StreamSpeech(context.Background(), "hello", "1")
	if err == nil {
		t.Errorf("expected error when API key is empty")
	}
}

func TestOpenAIProvider_StreamSpeech_NoKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	p := NewOpenAIProvider()
	_, err := p.StreamSpeech(context.Background(), "hello", "1")
	if err == nil {
		t.Errorf("expected error when API key is empty")
	}
}
