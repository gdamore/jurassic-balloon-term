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
	if !visible || frame.Width != 1 || frame.Height != 2 {
		t.Fatalf("strike frame = visible:%t %dx%d, want 1x2", visible, frame.Width, frame.Height)
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

func TestThunderheadCeilingRebuildsOnResize(t *testing.T) {
	m := NewStageModel(80, 24, 2)
	old := m.ceiling
	m.Resize(100, 32)
	if m.ceiling == old || m.ceiling.width != 100 {
		t.Fatalf("resized ceiling = same:%t width:%d", m.ceiling == old, m.ceiling.width)
	}
}
