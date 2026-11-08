// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestFernValleyBeginsWithHarmlessFerns(t *testing.T) {
	m := NewStageModel(80, 24, 0)
	if len(m.ferns) != 5 {
		t.Fatalf("fern count = %d, want 5", len(m.ferns))
	}
	for index, fern := range m.ferns {
		if fern.sprite.Collidable() {
			t.Fatalf("fern %d is collidable", index)
		}
		if fern.distance >= m.stage.GoalDistance/2 {
			t.Fatalf("fern %d distance = %v, want first half", index, fern.distance)
		}
	}
}

func TestRedEyeRunHasHarmlessFernScenery(t *testing.T) {
	m := NewStageModel(80, 24, 1)
	if len(m.ferns) != 5 {
		t.Fatalf("fern count = %d, want 5", len(m.ferns))
	}
	for index, fern := range m.ferns {
		if fern.sprite.Collidable() {
			t.Fatalf("fern %d is collidable", index)
		}
	}
}

func TestFernArtworkSwaysAndUsesShading(t *testing.T) {
	fern := newFern(0, 20)
	frame, _ := fern.sprite.ActiveFrame()
	if frame.Width != 7 || frame.Height != fernHeight {
		t.Fatalf("fern size = %dx%d, want 7x%d", frame.Width, frame.Height, fernHeight)
	}
	if highlight := frame.CellAt(3, 0); highlight.Style != StyleFernHighlight {
		t.Fatalf("fern highlight style = %d, want %d", highlight.Style, StyleFernHighlight)
	}
	if leaf := frame.CellAt(3, 1); leaf.Style != StyleFernLeaf {
		t.Fatalf("fern leaf style = %d, want %d", leaf.Style, StyleFernLeaf)
	}
	if stem := frame.CellAt(3, 3); stem.Style != StyleFernStem {
		t.Fatalf("fern stem style = %d, want %d", stem.Style, StyleFernStem)
	}
	fern.sprite.Update(750 * time.Millisecond)
	if got := fern.sprite.FrameName(); got != "right" {
		t.Fatalf("fern frame after sway = %q, want right", got)
	}
}

func TestFernsRemainGroundedOnResize(t *testing.T) {
	m := NewStageModel(80, 24, 0)
	m.Resize(100, 32)
	balloonX, _ := m.balloon.Position()
	for _, fern := range m.ferns {
		x, y := fern.sprite.Position()
		if x != balloonX+fern.distance-m.distance {
			t.Fatalf("fern at distance %v has x %v", fern.distance, x)
		}
		if y != 32-fernHeight-1 {
			t.Fatalf("fern at distance %v has y %v", fern.distance, y)
		}
	}
}
