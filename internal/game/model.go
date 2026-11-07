// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import "time"

const moveStep = 2

// Model contains game state without any terminal-specific dependencies.
type Model struct {
	Width    int
	Height   int
	BalloonX int
	BalloonY int
	Paused   bool
	Elapsed  time.Duration
}

// NewModel creates a model centered in the supplied play field.
func NewModel(width, height int) *Model {
	m := &Model{}
	m.Resize(width, height)
	return m
}

// Resize changes the play field and keeps the balloon on screen.
func (m *Model) Resize(width, height int) {
	firstSize := m.Width == 0 || m.Height == 0
	m.Width = max(width, 1)
	m.Height = max(height, 1)
	if firstSize {
		m.BalloonX = m.Width / 2
		m.BalloonY = m.Height / 2
	}
	m.clamp()
}

// Move nudges the balloon and clamps it to the play field.
func (m *Model) Move(dx, dy int) {
	m.BalloonX += dx * moveStep
	m.BalloonY += dy
	m.clamp()
}

// TogglePause pauses or resumes simulation time.
func (m *Model) TogglePause() {
	m.Paused = !m.Paused
}

// Update advances the deterministic simulation by one fixed step.
func (m *Model) Update(step time.Duration) {
	if !m.Paused {
		m.Elapsed += step
	}
}

func (m *Model) clamp() {
	m.BalloonX = min(max(m.BalloonX, 0), m.Width-1)
	m.BalloonY = min(max(m.BalloonY, 0), m.Height-1)
}
