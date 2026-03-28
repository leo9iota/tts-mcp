# FEAT-007: OpenAI TTS Provider Extraction

> **Status:** Draft
> **Priority:** P2 Medium
> **Package:** `internal/tts/`
> **Stack:** Go (OpenAI SDK)
> **Domain:** Backend API / TTS Providers

---

## 1. Overview

To conclude full feature parity with `blacktop/mcp-tts`, integrating `openai_tts` cleanly matches the standard multi-model TTS provider.

### Why Now?

- **Industry Standard Fallback:** OpenAI's TTS is almost universally accessible for enterprise testing. It acts as a stable reference implementation during QA scenarios.

---

## 2. Architecture & Strategy

**Approach Evaluated:**
- Rely dynamically on `github.com/openai/openai-go` or raw HTTP REST implementations inside `internal/tts/openai.go` to construct the voice sequence. By returning raw `.mp3` bytes under `io.ReadCloser`, we ensure the `FEAT-002: Blacktop Streaming Engine` remains completely oblivious to where the audio came from.

---

## 3. Implementation Phases

### Phase 1: Struct Implementation

- [ ] Write `internal/tts/openai.go`.
- [ ] Ensure `OPENAI_API_KEY` safety blocks exist.
- [ ] Use standard `"alloy"` fallback if `voiceID` resolves empty.

### Phase 2: MCP Tool Definition

- [ ] Map `ToolName()` to return `"openai_tts"`.

---

## 4. Acceptance Criteria

- [ ] `openai_tts` seamlessly streams its byte output into our local hardware player identical to Fish Audio.
