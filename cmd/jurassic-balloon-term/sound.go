// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package main

import (
	"github.com/gdamore/jurassic-balloon-term/internal/audio"
	"github.com/gdamore/jurassic-balloon-term/internal/sound"
)

type gameSound struct {
	engine *audio.Engine
}

func newGameSound() (*gameSound, error) {
	engine, err := audio.New(audio.NewOtoBackend(), sound.Catalog(), audio.Options{})
	if err != nil {
		return nil, err
	}
	return &gameSound{engine: engine}, nil
}

func (s *gameSound) Start() {
	_ = s.engine.Resume()
	_ = s.engine.SetEnvironment(audio.EnvironmentPrehistoricDay)
	_ = s.engine.SetMusic(audio.MusicClearSkies)
}

func (s *gameSound) SetBurner(enabled bool) {
	_ = s.engine.SetContinuous(audio.ContinuousBurner, enabled)
}

func (s *gameSound) BalloonPopped() {
	_ = s.engine.SetContinuous(audio.ContinuousBurner, false)
	_ = s.engine.Play(audio.CueBalloonPop)
	_ = s.engine.SetMusic(audio.MusicFailure)
}

func (s *gameSound) PterosaurLunged() {
	_ = s.engine.Play(audio.CuePterosaurCall)
}

func (s *gameSound) LightningStruck() {
	_ = s.engine.Play(audio.CueLightningCrack)
}

func (s *gameSound) StageCompleted() {
	_ = s.engine.SetContinuous(audio.ContinuousBurner, false)
	_ = s.engine.SetMusic(audio.MusicVictory)
}

func (s *gameSound) SetMusicEnabled(enabled bool) {
	gain := float32(0)
	if enabled {
		gain = 0.35
	}
	_ = s.engine.SetBusGain(audio.BusMusic, gain)
}

func (s *gameSound) SetEffectsEnabled(enabled bool) {
	effects, ambience, ui := float32(0), float32(0), float32(0)
	if enabled {
		effects, ambience, ui = 1, 0.7, 0.8
	}
	_ = s.engine.SetBusGain(audio.BusEffects, effects)
	_ = s.engine.SetBusGain(audio.BusAmbience, ambience)
	_ = s.engine.SetBusGain(audio.BusUI, ui)
}

func (s *gameSound) Retry() {
	_ = s.engine.SetContinuous(audio.ContinuousBurner, false)
	_ = s.engine.SetEnvironment(audio.EnvironmentPrehistoricDay)
	_ = s.engine.SetMusic(audio.MusicClearSkies)
}

func (s *gameSound) SetPaused(paused bool) {
	if paused {
		_ = s.engine.Suspend()
		return
	}
	_ = s.engine.Resume()
}

func (s *gameSound) Close() {
	_ = s.engine.Close()
}
