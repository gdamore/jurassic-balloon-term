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
	if m.BalloonX != 0 || m.BalloonY != 0 {
		t.Fatalf("move above bounds = (%d, %d), want (0, 0)", m.BalloonX, m.BalloonY)
	}

	for range 20 {
		m.Move(1, 1)
	}
	if m.BalloonX != 19 || m.BalloonY != 9 {
		t.Fatalf("move below bounds = (%d, %d), want (19, 9)", m.BalloonX, m.BalloonY)
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

	if m.BalloonX != 9 || m.BalloonY != 4 {
		t.Fatalf("position after resize = (%d, %d), want (9, 4)", m.BalloonX, m.BalloonY)
	}
}
