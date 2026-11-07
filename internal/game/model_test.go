// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestMoveStaysInBounds(t *testing.T) {
	m := NewModel(20, 10)

	for range 20 {
		m.Move(-1, -1)
	}
	if position := m.BalloonPosition(); position.X != 0 || position.Y != 0 {
		t.Fatalf("move above bounds = %v, want (0, 0)", position)
	}

	for range 20 {
		m.Move(1, 1)
	}
	if position := m.BalloonPosition(); position.X != 15 || position.Y != 6 {
		t.Fatalf("move below bounds = %v, want (15, 6)", position)
	}
}

func TestUpdateHonorsPause(t *testing.T) {
	m := NewModel(20, 10)
	m.Update(10 * time.Millisecond)
	m.TogglePause()
	m.Update(10 * time.Millisecond)

	if m.Elapsed != 10*time.Millisecond {
		t.Fatalf("elapsed = %v, want 10ms", m.Elapsed)
	}
}

func TestResizeClampsPosition(t *testing.T) {
	m := NewModel(80, 24)
	m.Move(100, 100)
	m.Resize(10, 5)

	if position := m.BalloonPosition(); position.X != 5 || position.Y != 1 {
		t.Fatalf("position after resize = %v, want (5, 1)", position)
	}
}

func TestUpdateAnimatesBalloon(t *testing.T) {
	m := NewModel(20, 10)
	if got := m.Sprites()[0].FrameName(); got != "sway-left" {
		t.Fatalf("initial frame = %q, want sway-left", got)
	}
	m.Update(300 * time.Millisecond)
	if got := m.Sprites()[0].FrameName(); got != "sway-right" {
		t.Fatalf("frame after update = %q, want sway-right", got)
	}
}
