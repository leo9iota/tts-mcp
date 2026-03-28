# FEAT-006: ElevenLabs Provider Extraction

> **Status:** Draft
> **Priority:** P2 Medium
> **Package:** `internal/tts/`
> **Stack:** Go (REST API)
> **Domain:** Backend API / TTS Providers

---

## 1. Overview

Importing the premium capability `elevenlabs_tts` directly from our `docs/code/` referenced `blacktop/mcp-tts` project into our lightweight Provider registry (`FEAT-003`).

### Why Now?

- **Feature Parity:** To fully consume the capabilities of the original inspirational project without its CLI bloat.

---

## 2. Architecture & Strategy

**Approach Evaluated:**
- Steal the ElevenLabs POST JSON schema natively written inside `blacktop/mcp-tts/cmd/root.go` and isolate it into a pure Go struct.
- Port their exact model configuration parameters (Stability, SimilarityBoost, etc.) into an `internal/tts/elevenlabs.go` implementation file.

---

## 3. Implementation Phases

### Phase 1: Struct Implementation

- [ ] Copy `SynthesisOptions` and `ElevenLabsParams` directly into `internal/tts/elevenlabs.go`.
- [ ] Define `Type ElevenLabsProvider struct{}` safely exposing the `ELEVENLABS_API_KEY` validation blocks.
- [ ] Tie down `StreamSpeech` dynamically building the `/v1/text-to-speech/{voiceID}/stream` POST.

### Phase 2: MCP Tool Definition

- [ ] Map `ToolName()` to return `"elevenlabs_tts"`.
- [ ] Match `Description()` block to identify ElevenLabs specific requirements.

---

## 4. Acceptance Criteria

- [ ] `elevenlabs_tts` triggers correctly if environmental values present.
