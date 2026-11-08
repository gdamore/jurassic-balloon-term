// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

// Sound receives semantic presentation events without exposing an audio
// implementation to the game model.
type Sound interface {
	Start()
	SetBurner(bool)
	BalloonPopped()
	PterosaurLunged()
	LightningStruck()
	StageCompleted()
	SetMusicEnabled(bool)
	SetEffectsEnabled(bool)
	Retry()
	SetPaused(bool)
}
