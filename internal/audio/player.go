package audio

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"
)

// AudioEngine encapsulates isolated audio stream context, eliminating
// the old package-level global deadlocks preventing multi-device routing scaling.
type AudioEngine struct {
	mu         sync.Mutex
	initOnce   sync.Once
	initErr    error
	sampleRate beep.SampleRate
	activeCtx  context.Context
}

// NewEngine safely invokes a completely separated driver tracking state.
func NewEngine() *AudioEngine {
	return &AudioEngine{}
}

// initSpeaker initializes the localized beep loop mapping to standard default speaker.
func (e *AudioEngine) initSpeaker(sampleRate beep.SampleRate) error {
	e.initOnce.Do(func() {
		e.sampleRate = sampleRate
		e.initErr = speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	})
	return e.initErr
}

// resampleToSpeaker buffers the playback sample mapping correctly to localized frequency execution paths.
func (e *AudioEngine) resampleToSpeaker(streamer beep.Streamer, from beep.SampleRate) beep.Streamer {
	if e.sampleRate == 0 || from == e.sampleRate {
		return streamer
	}
	return beep.Resample(4, from, e.sampleRate, streamer)
}

// WaitAndPlay blocks actively listening on localized mutex pointer until speaker hardware explicitly succeeds sequence execution.
func (e *AudioEngine) WaitAndPlay(ctx context.Context, stream beep.StreamSeeker, originalRate beep.SampleRate, volumeMultiplier float64, reporter func(pos int, total int, message string)) error {
	e.mu.Lock()
	e.activeCtx = ctx
	defer func() {
		e.activeCtx = nil
		e.mu.Unlock()
	}()

	if err := e.initSpeaker(originalRate); err != nil {
		return fmt.Errorf("failed to init speaker driver: %v", err)
	}

	var playback beep.Streamer = stream
	if volumeMultiplier <= 0.01 {
		playback = &effects.Volume{Streamer: stream, Base: 2, Volume: 0, Silent: true}
	} else if volumeMultiplier != 1.0 {
		playback = &effects.Volume{Streamer: stream, Base: 2, Volume: math.Log2(volumeMultiplier)}
	}

	playback = e.resampleToSpeaker(playback, originalRate)
	done := make(chan bool, 1)

	if reporter != nil {
		go func() {
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					total := stream.Len()
					pos := stream.Position()
					msg := fmt.Sprintf("Playing: %.1fs", float64(pos)/float64(originalRate))
					if total > 0 {
						msg = fmt.Sprintf("Playing: %.1fs / %.1fs", float64(pos)/float64(originalRate), float64(total)/float64(originalRate))
					}
					reporter(pos, total, msg)
				}
			}
		}()
	}

	speaker.Play(beep.Seq(playback, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

// Stop sends instantaneous silence bytes to speaker buffer directly wiping hardware stream locking buffers
// Only triggers if the context trying to stop it is currently the active playback.
func (e *AudioEngine) Stop(ctx context.Context) {
	if e.activeCtx != nil && e.activeCtx == ctx {
		speaker.Clear()
	}
}
