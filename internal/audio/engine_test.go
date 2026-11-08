// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"encoding/binary"
	"io"
	"math"
	"sync"
	"testing"
	"time"
)

type captureBackend struct {
	mu        sync.Mutex
	source    io.Reader
	format    Format
	started   bool
	suspended bool
	closed    bool
}

func (b *captureBackend) Start(source io.Reader, format Format) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return ErrBackendStarted
	}
	b.started = true
	b.source = source
	b.format = format
	return nil
}

func (b *captureBackend) Suspend() error {
	b.mu.Lock()
	b.suspended = true
	b.mu.Unlock()
	return nil
}

func (b *captureBackend) Resume() error {
	b.mu.Lock()
	b.suspended = false
	b.mu.Unlock()
	return nil
}

func (b *captureBackend) Close() error {
	b.mu.Lock()
	b.closed = true
	b.mu.Unlock()
	return nil
}

func (b *captureBackend) pull(t *testing.T, sampleCount int) []float32 {
	t.Helper()
	b.mu.Lock()
	source := b.source
	b.mu.Unlock()
	data := make([]byte, sampleCount*4)
	if _, err := io.ReadFull(source, data); err != nil {
		t.Fatalf("read mixed PCM: %v", err)
	}
	result := make([]float32, sampleCount)
	for i := range result {
		result[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:]))
	}
	return result
}

func TestEngineRoutesSemanticCues(t *testing.T) {
	backend := &captureBackend{}
	engine, err := New(backend, testCatalog(), testOptions())
	if err != nil {
		t.Fatalf("New() returned %v", err)
	}
	defer engine.Close()

	if err := engine.Play(CueBalloonPop); err != nil {
		t.Fatalf("Play() returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0.75, 0.75})

	// MaxVoices prevents an accidental stack of identical pop sounds.
	if err := engine.Play(CueBalloonPop); err != nil {
		t.Fatalf("Play() returned %v", err)
	}
	if err := engine.Play(CueBalloonPop); err != nil {
		t.Fatalf("Play() returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0.75, 0.75})
}

func TestEngineMusicLoopsAndUsesMusicBus(t *testing.T) {
	backend := &captureBackend{}
	engine, err := New(backend, testCatalog(), testOptions())
	if err != nil {
		t.Fatalf("New() returned %v", err)
	}
	defer engine.Close()

	if err := engine.SetBusGain(BusMusic, 0.5); err != nil {
		t.Fatalf("SetBusGain() returned %v", err)
	}
	if err := engine.SetMusic(MusicClearSkies); err != nil {
		t.Fatalf("SetMusic() returned %v", err)
	}
	wantSamples(t, backend.pull(t, 6), []float32{0.1, 0.1, 0.2, 0.2, 0.1, 0.1})
	if err := engine.SetMusic(MusicNone); err != nil {
		t.Fatalf("SetMusic(MusicNone) returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0, 0})
}

func TestEngineContinuousIsIdempotent(t *testing.T) {
	backend := &captureBackend{}
	engine, err := New(backend, testCatalog(), testOptions())
	if err != nil {
		t.Fatalf("New() returned %v", err)
	}
	defer engine.Close()

	if err := engine.SetContinuous(ContinuousBurner, true); err != nil {
		t.Fatalf("SetContinuous() returned %v", err)
	}
	if err := engine.SetContinuous(ContinuousBurner, true); err != nil {
		t.Fatalf("repeated SetContinuous() returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0.4, 0.4})
	if err := engine.SetContinuous(ContinuousBurner, false); err != nil {
		t.Fatalf("SetContinuous(false) returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0, 0})
}

func TestEngineEnvironmentStartsBed(t *testing.T) {
	backend := &captureBackend{}
	engine, err := New(backend, testCatalog(), testOptions())
	if err != nil {
		t.Fatalf("New() returned %v", err)
	}
	defer engine.Close()

	if err := engine.SetEnvironment(EnvironmentPrehistoricDay); err != nil {
		t.Fatalf("SetEnvironment() returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0.1, 0.1})
	if err := engine.SetEnvironment(EnvironmentNone); err != nil {
		t.Fatalf("SetEnvironment(EnvironmentNone) returned %v", err)
	}
	wantSamples(t, backend.pull(t, 2), []float32{0, 0})
}

func TestEngineBackendLifecycleAndClosedErrors(t *testing.T) {
	backend := &captureBackend{}
	engine, err := New(backend, testCatalog(), testOptions())
	if err != nil {
		t.Fatalf("New() returned %v", err)
	}
	if err := engine.Suspend(); err != nil {
		t.Fatalf("Suspend() returned %v", err)
	}
	if err := engine.Resume(); err != nil {
		t.Fatalf("Resume() returned %v", err)
	}
	if err := engine.Close(); err != nil {
		t.Fatalf("Close() returned %v", err)
	}
	if err := engine.Close(); err != nil {
		t.Fatalf("second Close() returned %v", err)
	}
	if err := engine.Play(CueBalloonPop); err != ErrClosed {
		t.Fatalf("Play() after Close = %v, want ErrClosed", err)
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	if !backend.started || backend.suspended || !backend.closed {
		t.Fatalf("backend state = started:%t suspended:%t closed:%t", backend.started, backend.suspended, backend.closed)
	}
}

func testCatalog() Catalog {
	const (
		pop    SoundID = "pop"
		music  SoundID = "clear-music"
		burner SoundID = "burner"
		wind   SoundID = "wind"
	)
	return Catalog{
		Format: Format{SampleRate: DefaultSampleRate, Channels: DefaultChannels},
		Sounds: map[SoundID]Sound{
			pop: {
				Variants:  []Clip{{PCM: stereo(0.75, 0.75)}},
				Bus:       BusEffects,
				MaxVoices: 1,
			},
			music: {
				Variants: []Clip{{PCM: stereo(0.2, 0.4)}},
				Bus:      BusMusic,
				Loop:     true,
			},
			burner: {
				Variants: []Clip{{PCM: stereo(0.4)}},
				Bus:      BusEffects,
				Loop:     true,
			},
			wind: {
				Variants: []Clip{{PCM: stereo(0.1)}},
				Bus:      BusAmbience,
				Loop:     true,
			},
		},
		Music:      map[Music]SoundID{MusicClearSkies: music},
		Continuous: map[Continuous]SoundID{ContinuousBurner: burner},
		Cues:       map[Cue]SoundID{CueBalloonPop: pop},
		Ambience: map[Environment]Soundscape{
			EnvironmentPrehistoricDay: {Bed: wind},
		},
	}
}

func testOptions() Options {
	return Options{
		MusicGain:      1,
		AmbienceGain:   1,
		EffectsGain:    1,
		UIGain:         1,
		MusicFade:      time.Nanosecond,
		ContinuousFade: time.Nanosecond,
	}
}

func stereo(values ...float32) []float32 {
	result := make([]float32, 0, len(values)*2)
	for _, value := range values {
		result = append(result, value, value)
	}
	return result
}

func wantSamples(t *testing.T, got, want []float32) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d samples, want %d", len(got), len(want))
	}
	for i := range got {
		if math.Abs(float64(got[i]-want[i])) > 0.0001 {
			t.Fatalf("sample %d = %v, want %v (all got: %v)", i, got[i], want[i], got)
		}
	}
}
