// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"fmt"
	"io"
	"sync"

	"github.com/ebitengine/oto/v3"
)

// OtoBackend supplies native output on macOS, Windows, and Linux and Web Audio
// output when built for js/wasm. Platform selection remains behind Backend.
type OtoBackend struct {
	mu      sync.Mutex
	started bool
	context *oto.Context
	player  *oto.Player
}

func NewOtoBackend() *OtoBackend {
	return &OtoBackend{}
}

func (b *OtoBackend) Start(source io.Reader, format Format) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBackendStarted
	}
	if source == nil || !format.valid() {
		return ErrInvalidFormat
	}
	b.started = true
	context, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:      format.SampleRate,
		ChannelCount:    format.Channels,
		Format:          oto.FormatFloat32LE,
		BufferSize:      otoOutputBufferSize,
		ApplicationName: "Jurassic Balloon",
	})
	if err != nil {
		return err
	}
	<-ready
	if err := context.Err(); err != nil {
		return fmt.Errorf("initialize oto context: %w", err)
	}
	b.context = context
	b.player = context.NewPlayer(source)
	b.player.SetBufferSize(otoPlayerBufferSize)
	b.player.Play()
	return nil
}

func (b *OtoBackend) Suspend() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.context == nil {
		return nil
	}
	return b.context.Suspend()
}

func (b *OtoBackend) Resume() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.context == nil {
		return nil
	}
	if err := b.context.Resume(); err != nil {
		return err
	}
	if b.player != nil {
		b.player.Play()
	}
	return nil
}

func (b *OtoBackend) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.player != nil {
		b.player.PauseAndStopReading()
		b.player = nil
	}
	if b.context == nil {
		return nil
	}
	err := b.context.Suspend()
	b.context = nil
	return err
}

var _ Backend = (*OtoBackend)(nil)
