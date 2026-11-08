// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

// Package audio turns semantic game cues into a mixed PCM stream.
//
// Callers select music, environments, continuous sounds, and one-shot cues.
// They do not select files or manipulate audio players. Device-specific output
// is supplied by a Backend.
package audio

import (
	"errors"
	"io"
	"time"
)

const (
	DefaultSampleRate = 48_000
	DefaultChannels   = 2
)

// Bus identifies an independently adjustable group in the mix.
type Bus uint8

const (
	BusMaster Bus = iota
	BusMusic
	BusAmbience
	BusEffects
	BusUI
	busCount
)

// Music identifies a music program, rather than an audio asset.
type Music uint8

const (
	MusicNone Music = iota
	MusicClearSkies
	MusicThreat
	MusicFailure
	MusicVictory
)

// Environment identifies a self-running ambient soundscape.
type Environment uint8

const (
	EnvironmentNone Environment = iota
	EnvironmentPrehistoricDay
	EnvironmentStorm
	EnvironmentVolcano
)

// Continuous identifies a sound whose lifetime follows game state.
type Continuous uint8

const (
	ContinuousBurner Continuous = iota
	ContinuousStrongWind
	ContinuousRain
	ContinuousVolcanoRumble
)

// Cue identifies a discrete semantic event.
type Cue uint8

const (
	CueBalloonPop Cue = iota
	CuePterosaurCall
	CueDistantDinosaur
	CueDistantSauropod
	CueDistantCarnivore
	CueThunderDistant
	CueThunderNear
	CuePickup
	CueMenuConfirm
	CueBasketCreak
	CueLightningCrack
	CueGroundLightningStrike
)

// Format describes the interleaved PCM stream delivered to a Backend.
type Format struct {
	SampleRate int
	Channels   int
}

func (f Format) valid() bool {
	return f.SampleRate > 0 && (f.Channels == 1 || f.Channels == 2)
}

// Backend delivers an endless, interleaved, little-endian float32 PCM stream
// to an audio device. Start is called exactly once.
type Backend interface {
	Start(io.Reader, Format) error
	Suspend() error
	Resume() error
	Close() error
}

// Options controls engine-wide behavior. Zero values are useful defaults.
type Options struct {
	SampleRate int
	Channels   int

	// MusicGain deliberately defaults below the other buses so that ambience
	// and gameplay feedback remain audible.
	MusicGain    float32
	AmbienceGain float32
	EffectsGain  float32
	UIGain       float32

	MusicFade      time.Duration
	ContinuousFade time.Duration
	RandomSeed     int64
}

func (o Options) normalized() Options {
	if o.SampleRate == 0 {
		o.SampleRate = DefaultSampleRate
	}
	if o.Channels == 0 {
		o.Channels = DefaultChannels
	}
	if o.MusicGain == 0 {
		o.MusicGain = 0.35
	}
	if o.AmbienceGain == 0 {
		o.AmbienceGain = 0.7
	}
	if o.EffectsGain == 0 {
		o.EffectsGain = 1
	}
	if o.UIGain == 0 {
		o.UIGain = 0.8
	}
	if o.MusicFade == 0 {
		o.MusicFade = 750 * time.Millisecond
	}
	if o.ContinuousFade == 0 {
		o.ContinuousFade = 80 * time.Millisecond
	}
	if o.RandomSeed == 0 {
		o.RandomSeed = 1
	}
	return o
}

var (
	ErrClosed         = errors.New("audio: engine is closed")
	ErrInvalidFormat  = errors.New("audio: invalid stream format")
	ErrBackendStarted = errors.New("audio: backend already started")
)
