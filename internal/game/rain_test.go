// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestThunderheadPassHasAnimatedHarmlessRain(t *testing.T) {
	for stageIndex := range StageCount() {
		m := NewStageModel(80, 24, stageIndex)
		if stageIndex == 2 {
			if m.rain == nil {
				t.Fatal("Thunderhead Pass has no rain field")
			}
			if m.rain.Collidable() {
				t.Fatal("rain field is collidable")
			}
			frame, _ := m.rain.ActiveFrame()
			if frame.Width != 80 || frame.Height != 24 {
				t.Fatalf("rain field size = %dx%d, want 80x24", frame.Width, frame.Height)
			}
			m.rain.Update(120 * time.Millisecond)
			if m.rain.FrameName() != "rain-2" {
				t.Fatalf("rain frame = %q, want rain-2", m.rain.FrameName())
			}
		} else if m.rain != nil {
			t.Fatalf("stage %d unexpectedly has persistent rain", stageIndex)
		}
	}
}

func TestRainFieldRebuildsOnResize(t *testing.T) {
	m := NewStageModel(80, 24, 2)
	old := m.rain
	m.Resize(100, 32)
	if m.rain == old {
		t.Fatal("resize retained old rain field")
	}
	frame, _ := m.rain.ActiveFrame()
	if frame.Width != 100 || frame.Height != 32 {
		t.Fatalf("resized rain field = %dx%d, want 100x32", frame.Width, frame.Height)
	}
}
