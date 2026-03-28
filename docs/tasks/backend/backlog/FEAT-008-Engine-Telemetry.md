# FEAT-008: Audio Engine Telemetry & Concurrency

> **Status:** Draft
> **Priority:** P1 High
> **Package:** `internal/audio/`
> **Stack:** Go (sync.Mutex, mcp-go NotifyProgress)
> **Domain:** Backend API / Telemetry

---

## 1. Overview

While `FEAT-002` provided raw streaming capability via `beep`, we must still extract the higher-level operational telemetry from `blacktop/mcp-tts`. Specifically: providing real-time playback percentage feedback to the IDE client, and protecting the audio pipeline from overlapping concurrent TTS requests.

### Why Now?

- **Anti-Timeout & UX:** VS Code and AI IDEs can time out if an MCP tool executes for 30 seconds without returning. Pinging the client with `mcp.NotifyProgress` keeps the connection alive and displays a sleek "Speaking: 2.1s / 5.0s" progress bar.
- **Overlapping Voices:** If the LLM produces two responses simultaneously, calling `generate_speech` twice will cross-contaminate the `beep` hardware buffer, causing the voices to talk over each other. We must implement `blacktop`'s global TTS locking.

---

## 2. Architecture & Strategy

**Approach Evaluated:**

- **Concurrency Mutexing:** Inside `internal/audio/player.go`, implement a `var TTSMutex sync.Mutex`. The `WaitAndPlay` command will actively `Lock()` the global context, forcing any secondary speech commands to queue sequentially rather than playing on top.
- **Progress Telemetry:** Import the architectural design of `progressReporter` from `blacktop`. We will use a `time.Ticker` every 250ms during playback to query `streamer.Position()` and report the fractional second output back through `session.NotifyProgress`.

---

## 3. Implementation Phases

### Phase 1: The Telemetry Reporter

- [ ] Create `internal/api/telemetry.go` or bind inside `audio/player.go`.
- [ ] Implement `StartProgressReporter(ctx context.Context, session *mcp.ServerSession, token any, total int, sampleRate int)`.
- [ ] Spin up a `time.Ticker` polling the `beep.Streamer` buffer fraction.

### Phase 2: The Global Mutex

- [ ] Wrap `WaitAndPlay` execution under the protection of `TTSMutex.Lock()`.
- [ ] Hook `defer TTSMutex.Unlock()` upon standard string completion and context cancellation alike.

---

## 4. Acceptance Criteria

- [ ] Rapid-firing the MCP tool inside Antigravity consecutively queues the audio flawlessly without overlapping.
- [ ] The IDE actively receives progress tokens and does not time out during long generation speeches.
