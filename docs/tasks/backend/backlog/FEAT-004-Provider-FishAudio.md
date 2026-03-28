# FEAT-004: FishAudio Provider Implementation

> **Status:** Draft
> **Priority:** P1 High
> **Package:** `internal/tts/`
> **Stack:** Go (REST API)
> **Domain:** Backend API / TTS Providers

---

## 1. Overview

Migrating our existing, hard-coded Fish Audio driver to respect the new `tts.Provider` architecture explicitly mapped in `FEAT-003`. 

### Why Now?

- **Backwards Compatibility:** We already proved Fish Audio works impeccably with the `beep` streaming engine in `FEAT-002`. We just need to shift it into the interface block instead of holding it as a loose package-level command.

---

## 2. Architecture & Strategy

**Approach Evaluated:**
Convert `internal/tts/client.go` into `internal/tts/fish.go`. We will build a struct `FishAudioProvider{}` that completely satisfies the interface.

---

## 3. Implementation Phases

### Phase 1: Struct Implementation

- [ ] Delete `internal/tts/client.go`.
- [ ] Create `internal/tts/fish.go`.
- [ ] Define `Type FishAudioProvider struct{}`.
- [ ] Implement `StreamSpeech(ctx context.Context, text string, voiceID string)` pulling exactly the logic we just proved works for FishAudio!

### Phase 2: MCP Tool Definition

- [ ] Map `ToolName()` to return `"fishaudio_tts"`.
- [ ] Map `Description()` to announce: `"Generates anime-style expressive TTS via Fish Audio REST API"`.

---

## 4. Acceptance Criteria

- [ ] The `FishAudioProvider` completely satisfies the `tts.Provider` interface at compile-time.
- [ ] Re-testing via MCP correctly fires `fishaudio_tts`.
