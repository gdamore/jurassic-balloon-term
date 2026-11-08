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
	s.SetStage(0)
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

func (s *gameSound) GroundLightningStruck() {
	_ = s.engine.Play(audio.CueGroundLightningStrike)
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

func (s *gameSound) SetStage(stageIndex int) {
	if stageIndex == 2 {
		_ = s.engine.SetEnvironment(audio.EnvironmentStorm)
	} else {
		_ = s.engine.SetEnvironment(audio.EnvironmentPrehistoricDay)
	}
	_ = s.engine.SetMusic(musicForStage(stageIndex))
}

func (s *gameSound) Retry() {
	_ = s.engine.SetContinuous(audio.ContinuousBurner, false)
}

func musicForStage(stageIndex int) audio.Music {
	switch stageIndex {
	case 2:
		return audio.MusicThreat
	default:
		return audio.MusicClearSkies
	}
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
