# TTS MCP

A high-performance, polyglot Text-to-Speech server bridging dynamic character personas and real-time audio playback natively into Google Antigravity and other MCP-enabled IDEs.

## Architecture & Data Flow

```mermaid
%%{init: {'theme': 'dark'}}%%
sequenceDiagram
    participant User as Google Antigravity IDE
    participant MCP as TTS MCP Server
    participant Persona as data/
    participant Provider as TTS Engine (FishAudio/Neets)
    participant Speaker as Native OS Audio Decoder
    participant Output as output/

    User->>MCP: Call `speak_as_persona` (text, persona: "Megumin")
    MCP->>Persona: Parse `data/megumin.json`
    Persona-->>MCP: Extract `voice_id` & `provider`
    MCP->>Provider: Stream text & voice mapping via HTTP payload

    par Real-time Bridging
        Provider-->>MCP: Raw `.mp3` byte stream
        MCP->>Speaker: Mount `go-beep` pipe directly to local speakers
        MCP->>Output: Clone sequence to `output/Megumin_2026-XX...mp3`
    end

    Speaker-->>User: Plays audio in real-time
    Output-->>User: Permanently saved artifact history
```
