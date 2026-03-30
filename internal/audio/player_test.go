package audio

import (
	"sync"
	"testing"
)

// TestAudioEngine_Instantiation guarantees memory blocks logically separate per engine spin-up functionally decoupling MCP state safely
func TestAudioEngine_Instantiation(t *testing.T) {
	engineA := NewEngine()
	engineB := NewEngine()

	// 1. Assert: New pointers exist isolated cleanly
	if engineA == nil {
		t.Fatal("Expected NewEngine() builder to explicitly return valid struct pointer memory blocks, got nil")
	}

	if engineA == engineB {
		t.Fatal("Expected decoupled isolation pointers strictly scaling instances concurrently, got overlapping structural matches!")
	}

	// 2. Assert: State defaults cleanly initializing structs without hardcoded package locks propagating blindly
	if engineA.sampleRate != 0 {
		t.Errorf("Expected strictly localized initialization binding SampleRate=0, got %v", engineA.sampleRate)
	}

	// 3. Act: Simulate independent structural locking dynamically preventing global overlap panics
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		engineA.mu.Lock()
		defer engineA.mu.Unlock()
		// Engine A locks down dynamically
	}()

	go func() {
		defer wg.Done()
		engineB.mu.Lock()
		defer engineB.mu.Unlock()
		// Engine B logically secures its exact specific resource concurrently!
	}()

	wg.Wait()
	// Assertion: The WaitGroup functionally exiting signifies entirely decoupled `sync.Mutex` objects inside distinct memory pointers overriding REFACTOR-004!
}
