# FEAT-005: Neets.ai Provider Implementation

> **Status:** Draft
> **Priority:** P1 High
> **Package:** `internal/tts/`
> **Stack:** Go (REST API)
> **Domain:** Backend API / TTS Providers

---

## 1. Overview

Integrating Neets.ai (`neets_tts`) into our polymorphic TTS provider registry (as established by `FEAT-003`). Neets represents an ultra-low-latency, massively subsidized alternative TTS engine, supporting multiple character accents.

### Why Now?

- **Cost Efficiency:** `blacktop/mcp-tts` supports OpenAI and ElevenLabs, which are exceptionally expensive per character. For developers constantly testing LLM agents, Neets acts as the definitive "cheap proxy".

---

## 2. Architecture & Strategy

**Approach Evaluated:**
Neets.ai functions virtually identically to OpenAI's TTS streaming POST models (REST endpoint returning raw audio buffers).

- Create `internal/tts/neets.go` implementing `Provider`.
- Inject `NEETS_API_KEY` mapping out from the OS environment.
- Use the target model id format to natively return robust `"mp3"` encoded binary responses that feed safely into the `FEAT-002` decoding architecture.

---

## 3. Implementation Phases

### Phase 1: Struct Implementation

- [ ] Create `internal/tts/neets.go`.
- [ ] Define `Type NeetsProvider struct{}`.
- [ ] Bind custom `NeetsRequest` structs satisfying the model parameters: `{"text", "voice_id", "fmt": "mp3"}`.
- [ ] Implement `StreamSpeech(ctx context.Context, text string, voiceID string)`.

### Phase 2: MCP Tool Definition

- [ ] Map `ToolName()` to return `"neets_tts"`.
- [ ] Map `Description()` to announce: `"Uses the Neets.ai REST API to generate conversational high-speed speech"`.

---

## 4. Acceptance Criteria

- [ ] Re-testing via MCP correctly executes `neets_tts`.
- [ ] Stream plays seamlessly over the `beep` hardware layer without crashing on chunking headers.
