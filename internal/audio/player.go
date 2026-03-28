package audio

import (
	"fmt"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

var (
	speakerInitOnce   sync.Once
	speakerInitErr    error
	speakerSampleRate beep.SampleRate
)

// initSpeaker initializes the speaker exactly once per process execution.
// Resampling prevents driver context teardowns.
func initSpeaker(sampleRate beep.SampleRate) error {
	speakerInitOnce.Do(func() {
		speakerSampleRate = sampleRate
		speakerInitErr = speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	})
	return speakerInitErr
}

// resampleToSpeaker matches the input stream rate to the global speaker dynamically lock.
func resampleToSpeaker(streamer beep.Streamer, from beep.SampleRate) beep.Streamer {
	if speakerSampleRate == 0 || from == speakerSampleRate {
		return streamer
	}
	return beep.Resample(4, from, speakerSampleRate, streamer)
}

// WaitAndPlay reads the decoded audio stream dynamically into the hardware speaker blocking until complete.
func WaitAndPlay(stream beep.Streamer, originalRate beep.SampleRate) error {
	if err := initSpeaker(originalRate); err != nil {
		return fmt.Errorf("failed to init speaker driver: %v", err)
	}

	playback := resampleToSpeaker(stream, originalRate)
	done := make(chan bool, 1)

	// Inject sequence callback to block thread via channel exactly like blacktop/mcp-tts
	speaker.Play(beep.Seq(playback, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

// Stop executes immediate explicit silence over the speaker hardware dropping active buffers
func Stop() {
	speaker.Clear()
}
