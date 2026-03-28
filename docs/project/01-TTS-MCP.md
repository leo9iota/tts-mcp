# TTS MCP

A lightweight, custom Model Context Protocol (MCP) server written in Go. This server acts as a cheaper, anime-focused alternative to ElevenLabs. It exposes a single tool to the AI agent, allowing it to dynamically generate and play highly expressive TTS (Text-to-Speech) using Fish Audio's REST API.

## Tech Stack

| Technology                  | Usage                                                 |
| :-------------------------- | :---------------------------------------------------- |
| Go (Golang)                 | Core programming language                             |
| mcp-go                      | MCP SDK for rapid tool generation and stdio transport |
| godotenv                    | Environment variable management                       |
| Fish Audio REST API         | Target API for speech generation                      |
| OS Audio (ffplay or native) | Audio playback execution                              |

## Targets

- Expected to run locally across platforms (Windows/Linux/macOS) as an MCP server.
- Used as a generic text-to-speech module by AI agent platforms.
