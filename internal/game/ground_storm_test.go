// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

func TestGroundStormIsLowLargeAndReachesTerrain(t *testing.T) {
	m := NewStageModel(80, 30, 2)
	if len(m.groundStorms) != 2 {
		t.Fatalf("ground storm count = %d, want 2", len(m.groundStorms))
	}
	for index, storm := range m.groundStorms {
		_, y := storm.cloud.Position()
		if y < float64(m.Height)/2 {
			t.Fatalf("ground storm %d y = %v above lower half", index, y)
		}
		storm.setPosition(float64(15+index*30), y)
		m.clipGroundStormsToTerrain()
		x, _ := storm.cloud.Position()
		cloud, _ := storm.cloud.ActiveFrame()
		if cloud.Width != groundStormCloudWidth || cloud.Height != groundStormCloudHeight {
			t.Fatalf("ground storm %d cloud = %dx%d", index, cloud.Width, cloud.Height)
		}
		if err := storm.bolt.SetFrame("strike-bright"); err != nil {
			t.Fatal(err)
		}
		bolt, _ := storm.bolt.ActiveFrame()
		boltX, boltY := storm.bolt.Position()
		surfaceY, found := m.terrain[len(m.terrain)-1].surfaceYAt(int(boltX) + groundStormBoltWidth/2)
		if !found {
			t.Fatalf("ground storm %d has no terrain below it", index)
		}
		lastCellY := -1
		for cellY := range bolt.Height {
			for cellX := range bolt.Width {
				if bolt.CellAt(cellX, cellY).Rune != 0 {
					lastCellY = cellY
				}
			}
		}
		if got := int(boltY) + lastCellY; got != surfaceY {
			t.Fatalf("ground storm %d bolt ends at %d, want terrain surface %d (cloud x %v)", index, got, surfaceY, x)
		}
		_, impactY := storm.impact.Position()
		if impactY != float64(surfaceY-1) {
			t.Fatalf("ground storm %d impact y = %v, want %d", index, impactY, surfaceY-1)
		}
	}
}

func TestGroundStormWarnsBeforeFatalStrike(t *testing.T) {
	storm := newGroundStormHazard(20, 15, 0.5, 30)
	storm.setPosition(12, 15)
	storm.updateForApproach(19)
	if storm.phase != groundStormWarn || !storm.bolt.Visible() || storm.bolt.Collidable() {
		t.Fatalf("warning state = phase:%v visible:%t collidable:%t", storm.phase, storm.bolt.Visible(), storm.bolt.Collidable())
	}
	warning, _ := storm.bolt.ActiveFrame()
	storm.setPosition(10, 15)
	if !storm.updateForApproach(19) || storm.phase != groundStormStrike || !storm.bolt.Collidable() {
		t.Fatalf("strike state = phase:%v collidable:%t", storm.phase, storm.bolt.Collidable())
	}
	if !storm.impact.Visible() {
		t.Fatal("ground impact is not visible at the start of the strike")
	}
	strike, _ := storm.bolt.ActiveFrame()
	for index := range warning.Cells {
		if (warning.Cells[index].Rune == 0) != (strike.Cells[index].Rune == 0) {
			t.Fatalf("warning and strike masks differ at cell %d", index)
		}
	}
	storm.bolt.Update(133 * time.Millisecond)
	if storm.bolt.FrameName() != "strike-fade" {
		t.Fatalf("first decay frame = %q", storm.bolt.FrameName())
	}
	if groundStormFrameHasCollision(storm.bolt) {
		t.Fatal("ground lightning fade still has a collision mask")
	}
	storm.bolt.Update(25 * time.Millisecond)
	if storm.bolt.FrameName() != "strike-dim" {
		t.Fatalf("final decay frame = %q", storm.bolt.FrameName())
	}
	storm.bolt.Update(60 * time.Millisecond)
	if storm.bolt.Visible() {
		t.Fatal("large lightning afterglow did not end promptly")
	}
	storm.setPosition(8, 15)
	storm.updateForApproach(19)
	if storm.phase != groundStormRest || storm.bolt.Visible() || storm.bolt.Collidable() {
		t.Fatalf("finished strike = phase:%v visible:%t collidable:%t", storm.phase, storm.bolt.Visible(), storm.bolt.Collidable())
	}
}

func groundStormFrameHasCollision(bolt *sprite.Sprite) bool {
	frame, visible := bolt.ActiveFrame()
	if !visible {
		return false
	}
	bounds, _ := bolt.Bounds()
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			if frame.CellAt(x-bounds.Min.X, y-bounds.Min.Y).Rune != 0 && bolt.CollidesAt(x, y) {
				return true
			}
		}
	}
	return false
}

func TestGroundStormStrikePopsLowBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 2)
	makeAirborne(m)
	storm := m.groundStorms[0]
	storm.phase = groundStormStrike
	storm.applyPhase()
	boltX, boltY := storm.bolt.Position()
	m.balloon.SetPosition(boltX-3, boltY)
	m.vy = 0

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("ground-reaching lightning did not pop low balloon")
	}
}

func TestGroundStormEmitsDedicatedStrikeEvent(t *testing.T) {
	m := NewStageModel(80, 30, 2)
	makeAirborne(m)
	m.balloon.SetPosition(35, 8)
	storm := m.groundStorms[0]
	storm.phase = groundStormWarn
	storm.applyPhase()
	// Place the bolt at the approach threshold. The balloon is above its top,
	// so this checks the semantic cue without causing a collision.
	storm.setPosition(35, 15)
	m.Update(time.Millisecond)

	events := m.Events()
	found := false
	for _, event := range events {
		found = found || event == EventGroundLightningStruck
	}
	if !found {
		t.Fatalf("ground strike events = %v, want dedicated event", events)
	}
}
