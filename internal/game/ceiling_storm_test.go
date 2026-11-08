// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestThunderheadCeilingTelegraphsShortDeterministicBolts(t *testing.T) {
	first := newCeilingStorm(80, 17)
	second := newCeilingStorm(80, 17)
	if !first.cloud.Collidable() || first.bolt.Visible() || first.bolt.Collidable() {
		t.Fatalf("initial ceiling state = cloud collidable:%t bolt visible:%t collidable:%t",
			first.cloud.Collidable(), first.bolt.Visible(), first.bolt.Collidable())
	}

	first.update(first.remaining)
	second.update(second.remaining)
	firstX, firstY := first.bolt.Position()
	secondX, secondY := second.bolt.Position()
	if first.phase != ceilingWarn || !first.bolt.Visible() || first.bolt.Collidable() {
		t.Fatalf("warning state = phase:%v visible:%t collidable:%t", first.phase, first.bolt.Visible(), first.bolt.Collidable())
	}
	if firstX != secondX || firstY != secondY || firstY != ceilingCloudRows {
		t.Fatalf("deterministic warning positions = (%v,%v) and (%v,%v)", firstX, firstY, secondX, secondY)
	}

	if !first.update(ceilingWarnTime) || first.phase != ceilingStrike || !first.bolt.Collidable() {
		t.Fatalf("strike state = phase:%v collidable:%t", first.phase, first.bolt.Collidable())
	}
	frame, visible := first.bolt.ActiveFrame()
	if !visible || frame.Width != ceilingBoltWidth || frame.Height != ceilingBoltRows {
		t.Fatalf("strike frame = visible:%t %dx%d, want %dx%d", visible, frame.Width, frame.Height, ceilingBoltWidth, ceilingBoltRows)
	}
	first.bolt.Update(133 * time.Millisecond)
	if first.bolt.FrameName() != "strike-fade" {
		t.Fatalf("ceiling lightning first fade = %q", first.bolt.FrameName())
	}
	if groundStormFrameHasCollision(first.bolt) {
		t.Fatal("ceiling lightning fade still has a collision mask")
	}
	first.bolt.Update(25 * time.Millisecond)
	if first.bolt.FrameName() != "strike-dim" {
		t.Fatalf("ceiling lightning afterglow = %q", first.bolt.FrameName())
	}
}

func TestThunderheadCeilingPopsBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 2)
	makeAirborne(m)
	m.balloon.SetPosition(20, 0)
	m.vy = 0
	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("thunderhead ceiling contact did not pop balloon")
	}
}

func TestThunderheadCeilingHasRoundedVariedCloudProfile(t *testing.T) {
	storm := newCeilingStorm(80, 17)
	frame, _ := storm.cloud.ActiveFrame()
	topRunes := map[rune]bool{}
	styles := map[int]bool{}
	undersideRunes := map[rune]bool{}
	for x := range frame.Width {
		if cell := frame.CellAt(x, 0); cell.Rune == 0 {
			t.Fatalf("ceiling crown has hole at x=%d", x)
		} else {
			topRunes[cell.Rune] = true
			styles[cell.Style] = true
		}
		for y := 1; y < frame.Height; y++ {
			cell := frame.CellAt(x, y)
			if cell.Rune == 0 {
				t.Fatalf("ceiling boundary has gap at (%d,%d)", x, y)
			}
			styles[cell.Style] = true
		}
		undersideRunes[frame.CellAt(x, frame.Height-1).Rune] = true
	}
	if len(topRunes) < 3 || len(styles) < 3 || len(undersideRunes) < 3 {
		t.Fatalf("ceiling profile lacks variation: crown:%d styles:%d underside:%d", len(topRunes), len(styles), len(undersideRunes))
	}
}

func TestThunderheadCeilingTextureScrollsBetweenTerrainRates(t *testing.T) {
	storm := newCeilingStorm(80, 17)
	before, _ := storm.cloud.ActiveFrame()
	storm.update(time.Second)
	after, _ := storm.cloud.ActiveFrame()
	for y := range before.Height {
		for x := 0; x < before.Width-1; x++ {
			if after.CellAt(x, y) != before.CellAt(x+1, y) {
				t.Fatalf("ceiling did not scroll at (%d,%d)", x, y)
			}
		}
	}
	farTerrainSpeed := 0.55 * stormScrollSpeed / worldScrollSpeed
	if ceilingScrollSpeed <= farTerrainSpeed || ceilingScrollSpeed >= stormScrollSpeed {
		t.Fatalf("ceiling speed %v is not between background and foreground rates", ceilingScrollSpeed)
	}
}

func TestThunderheadCeilingRebuildsOnResize(t *testing.T) {
	m := NewStageModel(80, 24, 2)
	old := m.ceiling
	m.Resize(100, 32)
	if m.ceiling == old || m.ceiling.width != 100 {
		t.Fatalf("resized ceiling = same:%t width:%d", m.ceiling == old, m.ceiling.width)
	}
}
