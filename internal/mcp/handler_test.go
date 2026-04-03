package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"tts-mcp/internal/audio"
	"tts-mcp/internal/personas"
	"tts-mcp/internal/providers"
)

func TestCreatePersonaHandler_MissingArgs(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	mng := &personas.Manager{Personas: map[string]personas.Persona{
		"test_persona": {Name: "test_persona", Provider: "dummy", VoiceID: "123"},
	}}
	var provs []providers.Provider
	engine := audio.NewEngine()

	handler := createPersonaHandler(s, mng, provs, engine)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = nil // invalid arguments

	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for invalid arguments")
	}
}

func TestCreatePersonaHandler_NotFound(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	mng := &personas.Manager{Personas: make(map[string]personas.Persona)}
	var provs []providers.Provider
	engine := audio.NewEngine()

	handler := createPersonaHandler(s, mng, provs, engine)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"persona": "unknown",
		"text":    "hello",
	}

	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for unknown persona")
	}

	if len(res.Content) > 0 {
		txt, ok := res.Content[0].(mcp.TextContent)
		if ok {
			if !strings.Contains(txt.Text, "not found") {
				t.Errorf("expected 'not found' message, got %v", txt.Text)
			}
		}
	}
}

func TestCreatePersonaGeneratorHandler_MissingParams(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	mng := &personas.Manager{Personas: make(map[string]personas.Persona)}

	handler := createPersonaGeneratorHandler(s, mng)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"name": "OnlyNameNoID",
	}

	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for missing voice_id/provider params")
	}
}
