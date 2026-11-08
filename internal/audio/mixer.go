// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"encoding/binary"
	"math"
	"sync"
)

type voiceID uint64

type voice struct {
	id       voiceID
	clip     Clip
	bus      Bus
	loop     bool
	position int

	gain          float32
	targetGain    float32
	gainStep      float32
	fadeRemaining int
	stopAfterFade bool
}

type mixer struct {
	mu     sync.Mutex
	format Format
	buses  [busCount]float32
	voices []*voice
	nextID voiceID
}

func newMixer(format Format) *mixer {
	m := &mixer{format: format}
	for i := range m.buses {
		m.buses[i] = 1
	}
	return m
}

func (m *mixer) play(clip Clip, bus Bus, gain float32, loop bool, fadeFrames int) voiceID {
	if gain == 0 {
		gain = 1
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	v := &voice{
		id:         m.nextID,
		clip:       clip,
		bus:        bus,
		loop:       loop,
		gain:       gain,
		targetGain: gain,
	}
	if fadeFrames > 0 {
		v.gain = 0
		v.fadeRemaining = fadeFrames
		v.gainStep = gain / float32(fadeFrames)
	}
	m.voices = append(m.voices, v)
	return v.id
}

func (m *mixer) stop(id voiceID, fadeFrames int) {
	if id == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.voices {
		if v.id != id {
			continue
		}
		if fadeFrames <= 0 {
			v.gain = 0
			v.stopAfterFade = true
			v.fadeRemaining = 0
			return
		}
		v.targetGain = 0
		v.fadeRemaining = fadeFrames
		v.gainStep = -v.gain / float32(fadeFrames)
		v.stopAfterFade = true
		return
	}
}

func (m *mixer) setBusGain(bus Bus, gain float32) {
	if bus >= busCount {
		return
	}
	if gain < 0 || math.IsNaN(float64(gain)) {
		gain = 0
	}
	m.mu.Lock()
	m.buses[bus] = gain
	m.mu.Unlock()
}

func (m *mixer) active(id voiceID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.voices {
		if v.id == id && !(v.stopAfterFade && v.fadeRemaining == 0) {
			return true
		}
	}
	return false
}

func (m *mixer) render(dst []float32) {
	clear(dst)
	frames := len(dst) / m.format.Channels
	if frames == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	kept := m.voices[:0]
	for _, v := range m.voices {
		alive := m.renderVoice(dst, frames, v)
		if alive {
			kept = append(kept, v)
		}
	}
	m.voices = kept

	master := m.buses[BusMaster]
	for i := range dst {
		sample := dst[i] * master
		if sample > 1 {
			sample = 1
		} else if sample < -1 {
			sample = -1
		}
		dst[i] = sample
	}
}

func (m *mixer) renderVoice(dst []float32, frames int, v *voice) bool {
	channels := m.format.Channels
	clipFrames := len(v.clip.PCM) / channels
	loopEnd := v.clip.LoopEnd
	if loopEnd == 0 {
		loopEnd = clipFrames
	}

	for frame := 0; frame < frames; frame++ {
		if v.stopAfterFade && v.fadeRemaining == 0 {
			return false
		}
		if v.position >= loopEnd {
			if !v.loop {
				return false
			}
			v.position = v.clip.LoopStart
		}

		gain := v.gain * m.buses[v.bus]
		source := v.position * channels
		target := frame * channels
		for channel := 0; channel < channels; channel++ {
			dst[target+channel] += v.clip.PCM[source+channel] * gain
		}
		v.position++

		if v.fadeRemaining > 0 {
			v.gain += v.gainStep
			v.fadeRemaining--
			if v.fadeRemaining == 0 {
				v.gain = v.targetGain
			}
		}
	}
	return true
}

// pcmReader adapts the mixer's frame-oriented float stream to Oto's byte
// reader contract. It never returns EOF; an idle engine produces silence.
type pcmReader struct {
	mixer   *mixer
	samples []float32
	bytes   []byte
	pending []byte
}

func (r *pcmReader) Read(p []byte) (int, error) {
	written := copy(p, r.pending)
	r.pending = r.pending[written:]
	if written == len(p) {
		return written, nil
	}

	needBytes := len(p) - written
	frameBytes := r.mixer.format.Channels * 4
	frames := (needBytes + frameBytes - 1) / frameBytes
	sampleCount := frames * r.mixer.format.Channels
	if cap(r.samples) < sampleCount {
		r.samples = make([]float32, sampleCount)
	} else {
		r.samples = r.samples[:sampleCount]
	}
	if cap(r.bytes) < sampleCount*4 {
		r.bytes = make([]byte, sampleCount*4)
	} else {
		r.bytes = r.bytes[:sampleCount*4]
	}

	r.mixer.render(r.samples)
	for i, sample := range r.samples {
		binary.LittleEndian.PutUint32(r.bytes[i*4:], math.Float32bits(sample))
	}
	copied := copy(p[written:], r.bytes)
	written += copied
	r.pending = r.bytes[copied:]
	return written, nil
}
