// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestMixerHonorsLoopPoints(t *testing.T) {
	mixer := newMixer(Format{SampleRate: 10, Channels: 1})
	mixer.play(Clip{
		PCM:       []float32{0.1, 0.2, 0.3, 0.4},
		LoopStart: 1,
		LoopEnd:   3,
	}, BusEffects, 1, true, 0)
	got := make([]float32, 6)
	mixer.render(got)
	wantSamples(t, got, []float32{0.1, 0.2, 0.3, 0.2, 0.3, 0.2})
}

func TestMixerFadesInAndOut(t *testing.T) {
	mixer := newMixer(Format{SampleRate: 10, Channels: 1})
	id := mixer.play(Clip{PCM: []float32{1, 1, 1, 1, 1, 1}}, BusEffects, 1, false, 2)
	got := make([]float32, 3)
	mixer.render(got)
	wantSamples(t, got, []float32{0, 0.5, 1})

	mixer.stop(id, 2)
	mixer.render(got)
	wantSamples(t, got, []float32{1, 0.5, 0})
}

func TestMixerAppliesBusGainAndClips(t *testing.T) {
	mixer := newMixer(Format{SampleRate: 10, Channels: 1})
	mixer.setBusGain(BusEffects, 0.5)
	mixer.play(Clip{PCM: []float32{2}}, BusEffects, 1.5, false, 0)
	got := make([]float32, 1)
	mixer.render(got)
	wantSamples(t, got, []float32{1})
}

func TestPCMReaderPreservesPartialFrames(t *testing.T) {
	mixer := newMixer(Format{SampleRate: 10, Channels: 2})
	mixer.play(Clip{PCM: []float32{0.25, -0.25}}, BusEffects, 1, false, 0)
	reader := &pcmReader{mixer: mixer}
	first := make([]byte, 3)
	second := make([]byte, 5)
	if n, err := reader.Read(first); err != nil || n != len(first) {
		t.Fatalf("first Read() = %d, %v", n, err)
	}
	if n, err := reader.Read(second); err != nil || n != len(second) {
		t.Fatalf("second Read() = %d, %v", n, err)
	}
	data := append(first, second...)
	left := math.Float32frombits(binary.LittleEndian.Uint32(data[:4]))
	right := math.Float32frombits(binary.LittleEndian.Uint32(data[4:]))
	wantSamples(t, []float32{left, right}, []float32{0.25, -0.25})
}
