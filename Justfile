# Build the tts-mcp executable
build:
    go build -o tts-mcp.exe .

# Run the mcp server locally over stdio
run: build
    ./tts-mcp.exe

# Start the mcp inspector to test the server visually
inspect: build
    npx @modelcontextprotocol/inspector ./tts-mcp.exe

# Clean binary and temporary audio files
clean:
    rm -f tts-mcp.exe temp.wav temp.mp3

# Install the dependencies
deps:
    go mod tidy
