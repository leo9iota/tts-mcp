package main

import (
	"fmt"
	"os/exec"
)

// PlayAudio uses ffplay to seamlessly play the generated audio
func PlayAudio() error {
	// Execute ffplay synchronously and silently so it doesn't pollute MCP stdio transport
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", "temp.wav")
	
	err := cmd.Run() // Blocks until playback is complete
	if err != nil {
		return fmt.Errorf("failed to play audio with ffplay: %v", err)
	}

	return nil
}
