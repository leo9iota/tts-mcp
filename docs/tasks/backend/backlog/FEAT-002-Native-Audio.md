# FEAT-002: Native Go Audio Playback & Client Delegation

> **Status:** Draft
> **Priority:** P1 High
> **Package:** `internal/audio/`
> **Stack:** Go (gopxl/beep)
> **Domain:** Backend API / Audio Engine

---

## 1. Overview

This specification addresses the dependency vulnerability established in FEAT-001 by replacing the external `ffplay` shell execution with a pure, cross-platform Go audio library. It also expands the MCP tool response to gracefully return the absolute path of the generated audio to the client, bridging the gap between host-side playback and client-side webview rendering.

### Why Now?

- **Dependency Free:** Forcing developers to install FFmpeg globally to use a Go binary defeats the purpose of Go's cross-platform portability.
- **VS Code Ideations:** AI IDEs (Antigravity, Cursor, Cline) typically do not contain complex, auto-playing audio renderers for tool call payloads in their textual chat screens. If we blindly delegate playback to the IDE, the audio simply won't play. We _must_ retain host-level execution while allowing future IDE clients the option to parse the payload themselves.

---

## 2. Architecture & Strategy

**Approach Evaluated:**

1. **Full Client Delegation (Rejected):** Returning base64 or a URI to Antigravity. Rejected because standard VS Code webviews don't auto-resolve returned MCP strings into audio `<audio>` players. The silent failure ruins the AI experience.
2. **Pure Go Audio Libraries (Accepted):** Moving to `github.com/gopxl/beep`. This library buffers standard `wav` streams directly into the Windows/macOS/Linux hardware APIs entirely without CGO (using `PureGo` via its `oto` sub-dependency).

**Execution Path:** We will implement a dual-pronged strategy.

1. `internal/audio/player.go` will be refactored to read the `.wav` headers dynamically and flush the PCM raw data to the OS via the `beep/speaker` API safely.
2. `internal/api/api.go` will be updated to return an MCP response containing the absolute local path to the generated file, allowing advanced clients to render it if they choose to.

---

## 3. Implementation Phases

### Phase 1: Pure Go Audio Dependencies

- [ ] Execute `go get github.com/gopxl/beep`.
- [ ] Add `encoding/binary` and standard IO routines to process the `wav` chunk headers.

### Phase 2: Refactor `audio` Engine

- [ ] Update `internal/audio/player.go` to remove the `os/exec` ffplay routine.
- [ ] Open `temp.wav` safely from disk.
- [ ] Initialize the `beep/speaker` context dynamically based on the decoded sample rate of the specific Fish Audio `wav` file (typically 44100Hz or 48000Hz).
- [ ] Buffer and play the audio byte chunks explicitly until EOF.

### Phase 3: Update Client Delegation Context

- [ ] Modify `internal/api/api.go` to grab the absolute file path of `temp.wav` using `filepath.Abs("temp.wav")`.
- [ ] Update the final `mcp.NewToolResultText` wrapper to return a detailed JSON or clear text message: e.g., `"Successfully generated speech. File saved locally at: C:\...\temp.wav"`.

### Phase 4: Build & Verify

- [ ] Run `just deps` to secure the new modules.
- [ ] Execute `just build` to confirm the binary compiles completely independently of CGO on Windows.

---

## 4. Acceptance Criteria

- [ ] The `tts-mcp.exe` binary runs on a Windows machine that does completely lacks FFmpeg or `ffplay`.
- [ ] The tool successfully reads a valid `.wav` file and plays it aloud without crashing or panicking.
- [ ] The audio buffer empties synchronously, meaning the server correctly waits for the audio to finish speaking before returning control to the AI.
- [ ] The output sent to the AI explicitly lists the absolute file path to the audio payload.
