package audio

import (
	"fmt"
	"os/exec"
)

// Play uses ffplay to seamlessly play the generated audio
func Play(filePath string) error {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", filePath)
	
	err := cmd.Run() // Blocks until playback is complete
	if err != nil {
		return fmt.Errorf("failed to play audio with ffplay: %v", err)
	}

	return nil
}
