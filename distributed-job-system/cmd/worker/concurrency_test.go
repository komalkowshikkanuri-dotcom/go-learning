package main

import (
	"sync"
	"sync/atomic"
	"testing"
	//"time"
)

func TestUpdatePeak(t *testing.T) {
	const expected int32 = 7

	var peak int32 = 2
	var current int32 = 7

	updatePeak(&peak, current)

	got := atomic.LoadInt32(&peak)
	if got != expected {
		t.Errorf("expected peak to become %d, got %d", expected, got)
	}
}

func TestUpdatePeakConcurrent(t *testing.T) {
	var peak int32 = 2
	values := []int32{3, 5, 4}

	var wg sync.WaitGroup

	wg.Add(len(values))

	for _, value := range values {
		go func(value int32) {
			defer wg.Done()
			updatePeak(&peak, value)
		}(value)
	}

	wg.Wait()

	const expected int32 = 5

	got := atomic.LoadInt32(&peak)

	if got != expected {
		t.Errorf("expected: %d, got: %d", expected, got)
	}
}

func TestSemaphoreLimit(t *testing.T) {
	const maxConcurrency = 3
	semaphore := make(chan struct{}, maxConcurrency)

	values := []int32{1, 2, 3, 4, 5, 6}
	started := make(chan struct{}, len(values))
	release := make(chan struct{})

	var current int32
	var peak int32

	var wg sync.WaitGroup
	wg.Add(len(values))

	go func() {
		for range values {
			semaphore <- struct{}{}
			go func() {
				defer wg.Done()
				defer atomic.AddInt32(&current, -1)
				defer func() {
					<-semaphore
				}()

				now := atomic.AddInt32(&current, 1)
				updatePeak(&peak, now)

				started <- struct{}{}

				<-release
			}()
		}
	}()

	for i := 0; i < maxConcurrency; i++ {
		<-started
	}

	select {
	case <-started:
		t.Fatal("a fourth worker started while all semaphore slots were occupied")
	default:
		// Correct: no fourth worker has started.
	}

	release <- struct{}{}
	<-started

	for i := 0; i < len(values)-1; i++ {
		release <- struct{}{}
	}

	wg.Wait()

	got := atomic.LoadInt32(&peak)

	if got != int32(maxConcurrency) {
		t.Errorf("expected peak concurrency %d, got %d", maxConcurrency, got)
	}
}
