package gofast

import (
	"testing"
)

func TestSequenceIDGenerator_Uniqueness(t *testing.T) {
	gen := NewSequenceIDGenerator()
	seen := make(map[string]bool, 1000)

	for range 1000 {
		id := gen.Next()
		if seen[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestSequenceIDGenerator_ConcurrentAccess(t *testing.T) {
	gen := NewSequenceIDGenerator()
	goroutines := 100
	n := 100

	results := make(chan string, goroutines*n)
	done := make(chan struct{})

	// Start concurrent goroutines generating IDs
	for range goroutines {
		go func() {
			for range n {
				results <- gen.Next()
			}
			done <- struct{}{}
		}()
	}

	// Wait for all goroutines to complete
	for range goroutines {
		<-done
	}
	close(results)

	// Check for duplicates
	seen := make(map[string]bool)
	for id := range results {
		if seen[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func BenchmarkSequenceIDGenerator_Next(b *testing.B) {
	generator := NewSequenceIDGenerator()

	for b.Loop() {
		_ = generator.Next()
	}
}
