// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"io"
	"sync"
	"time"
)

// NullBackend advances the stream at approximately real time and discards it.
// It is suitable for tests, headless runs, and graceful fallback when no audio
// device is available.
type NullBackend struct {
	mu        sync.Mutex
	started   bool
	suspended bool
	done      chan struct{}
	wg        sync.WaitGroup
}

func NewNullBackend() *NullBackend {
	return &NullBackend{}
}

func (b *NullBackend) Start(source io.Reader, format Format) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBackendStarted
	}
	if source == nil || !format.valid() {
		return ErrInvalidFormat
	}
	b.started = true
	b.done = make(chan struct{})
	bytesPerTick := format.SampleRate * format.Channels * 4 / 100
	b.wg.Add(1)
	go b.consume(source, bytesPerTick)
	return nil
}

func (b *NullBackend) consume(source io.Reader, bytesPerTick int) {
	defer b.wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	buffer := make([]byte, bytesPerTick)
	for {
		select {
		case <-ticker.C:
			b.mu.Lock()
			suspended := b.suspended
			b.mu.Unlock()
			if !suspended {
				_, _ = io.ReadFull(source, buffer)
			}
		case <-b.done:
			return
		}
	}
}

func (b *NullBackend) Suspend() error {
	b.mu.Lock()
	b.suspended = true
	b.mu.Unlock()
	return nil
}

func (b *NullBackend) Resume() error {
	b.mu.Lock()
	b.suspended = false
	b.mu.Unlock()
	return nil
}

func (b *NullBackend) Close() error {
	b.mu.Lock()
	if !b.started {
		b.mu.Unlock()
		return nil
	}
	b.started = false
	close(b.done)
	b.mu.Unlock()
	b.wg.Wait()
	return nil
}

var _ Backend = (*NullBackend)(nil)
