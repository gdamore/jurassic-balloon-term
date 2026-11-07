// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const moveStep = 2

// Model contains game state without any terminal-specific dependencies.
type Model struct {
	Width   int
	Height  int
	Paused  bool
	Elapsed time.Duration

	scene   *sprite.Scene
	balloon *sprite.Sprite
}

// NewModel creates a model centered in the supplied play field.
func NewModel(width, height int) *Model {
	m := &Model{scene: sprite.NewScene(), balloon: newBalloon()}
	m.scene.Add(m.balloon)
	m.Resize(width, height)
	return m
}

func newBalloon() *sprite.Sprite {
	balloon := sprite.New("balloon")
	a := sprite.FrameFromStrings(
		" .-. ",
		"(   )",
		" `-' ",
		"  |  ",
	)
	a.Duration, a.Next = 300*time.Millisecond, "sway-right"
	b := sprite.FrameFromStrings(
		" .-. ",
		"(   )",
		" `-' ",
		"   | ",
	)
	b.Duration, b.Next = 300*time.Millisecond, "sway-left"
	if err := balloon.AddFrame("sway-left", a); err != nil {
		panic(err)
	}
	if err := balloon.AddFrame("sway-right", b); err != nil {
		panic(err)
	}
	if err := balloon.SetFrame("sway-left"); err != nil {
		panic(err)
	}
	balloon.SetLayer(10)
	return balloon
}

// Resize changes the play field and keeps the balloon on screen.
func (m *Model) Resize(width, height int) {
	firstSize := m.Width == 0 || m.Height == 0
	m.Width = max(width, 1)
	m.Height = max(height, 1)
	if firstSize {
		frame, _ := m.balloon.ActiveFrame()
		m.balloon.SetPosition(float64(max((m.Width-frame.Width)/2, 0)), float64(max((m.Height-frame.Height)/2, 0)))
	}
	m.clampBalloon()
}

// Move nudges the balloon and clamps it to the play field.
func (m *Model) Move(dx, dy int) {
	position := m.balloon.CellPosition()
	m.balloon.SetPosition(float64(position.X+dx*moveStep), float64(position.Y+dy))
	m.clampBalloon()
}

func (m *Model) BalloonPosition() sprite.Point { return m.balloon.CellPosition() }

// Sprites returns the layer-ordered scene snapshot used by renderers.
func (m *Model) Sprites() []*sprite.Sprite { return m.scene.Sprites() }

// TogglePause pauses or resumes simulation time.
func (m *Model) TogglePause() { m.Paused = !m.Paused }

// Update advances the deterministic simulation by one fixed step and returns
// any collisions detected after movement.
func (m *Model) Update(step time.Duration) []sprite.Collision {
	if m.Paused {
		return nil
	}
	m.Elapsed += step
	return m.scene.Update(step)
}

func (m *Model) clampBalloon() {
	frame, ok := m.balloon.ActiveFrame()
	if !ok {
		return
	}
	position := m.balloon.CellPosition()
	maxX := max(m.Width-frame.Width, 0)
	maxY := max(m.Height-frame.Height, 0)
	m.balloon.SetPosition(float64(min(max(position.X, 0), maxX)), float64(min(max(position.Y, 0), maxY)))
}
