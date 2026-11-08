// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"fmt"
	"time"
)

// SoundID identifies an engine-owned asset. It is deliberately separate from
// the semantic cue types used by gameplay code.
type SoundID string

// Clip contains interleaved float32 PCM in the Catalog's format. LoopEnd is
// exclusive and expressed in frames. A zero LoopEnd means the end of the clip.
type Clip struct {
	PCM       []float32
	LoopStart int
	LoopEnd   int
}

// Sound defines how one or more interchangeable clips are played.
type Sound struct {
	Variants []Clip
	Bus      Bus
	// Gain defaults to one when zero.
	Gain      float32
	Loop      bool
	MaxVoices int
	Cooldown  time.Duration
}

// Soundscape describes an ambient bed plus sparse randomized details.
type Soundscape struct {
	Bed      SoundID
	Random   []SoundID
	MinDelay time.Duration
	MaxDelay time.Duration
}

// Catalog owns the relationship between semantic events and actual sounds.
// Asset loading and license/provenance records can be layered on top without
// changing Engine's public cue API.
type Catalog struct {
	Format Format
	Sounds map[SoundID]Sound

	Music      map[Music]SoundID
	Continuous map[Continuous]SoundID
	Cues       map[Cue]SoundID
	Ambience   map[Environment]Soundscape
}

// Validate checks structural properties needed by the real-time mixer.
func (c Catalog) Validate() error {
	if !c.Format.valid() {
		return ErrInvalidFormat
	}
	for id, sound := range c.Sounds {
		if len(sound.Variants) == 0 {
			return fmt.Errorf("audio: sound %q has no variants", id)
		}
		if sound.Bus <= BusMaster || sound.Bus >= busCount {
			return fmt.Errorf("audio: sound %q has invalid bus %d", id, sound.Bus)
		}
		for i, clip := range sound.Variants {
			if err := validateClip(c.Format, clip); err != nil {
				return fmt.Errorf("audio: sound %q variant %d: %w", id, i, err)
			}
		}
	}
	for music, id := range c.Music {
		if music != MusicNone && !c.hasSound(id) {
			return fmt.Errorf("audio: music %d refers to missing sound %q", music, id)
		}
	}
	for continuous, id := range c.Continuous {
		if !c.hasSound(id) {
			return fmt.Errorf("audio: continuous cue %d refers to missing sound %q", continuous, id)
		}
	}
	for cue, id := range c.Cues {
		if !c.hasSound(id) {
			return fmt.Errorf("audio: cue %d refers to missing sound %q", cue, id)
		}
	}
	for environment, scene := range c.Ambience {
		if scene.Bed != "" && !c.hasSound(scene.Bed) {
			return fmt.Errorf("audio: environment %d refers to missing bed %q", environment, scene.Bed)
		}
		for _, id := range scene.Random {
			if !c.hasSound(id) {
				return fmt.Errorf("audio: environment %d refers to missing sound %q", environment, id)
			}
		}
		if scene.MinDelay < 0 || scene.MaxDelay < scene.MinDelay || (len(scene.Random) > 0 && scene.MinDelay == 0) {
			return fmt.Errorf("audio: environment %d has invalid ambience delay", environment)
		}
	}
	return nil
}

func validateClip(format Format, clip Clip) error {
	if len(clip.PCM) == 0 || len(clip.PCM)%format.Channels != 0 {
		return fmt.Errorf("PCM length is not a positive number of frames")
	}
	frames := len(clip.PCM) / format.Channels
	end := clip.LoopEnd
	if end == 0 {
		end = frames
	}
	if clip.LoopStart < 0 || clip.LoopStart >= end || end > frames {
		return fmt.Errorf("invalid loop range [%d,%d) for %d frames", clip.LoopStart, end, frames)
	}
	return nil
}

func (c Catalog) hasSound(id SoundID) bool {
	_, ok := c.Sounds[id]
	return ok
}
