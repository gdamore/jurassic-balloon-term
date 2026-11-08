// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestFernValleyTreesOccupySecondHalf(t *testing.T) {
	stage := StageAt(0)
	if len(stage.TreeDistances) != 5 {
		t.Fatalf("Fern Valley tree count = %d, want 5", len(stage.TreeDistances))
	}
	previous := stage.GoalDistance / 2
	for index, distance := range stage.TreeDistances {
		if distance <= previous || distance >= stage.GoalDistance {
			t.Fatalf("tree %d distance = %v after %v with goal %v", index, distance, previous, stage.GoalDistance)
		}
		previous = distance
	}
	for index := 1; index < StageCount(); index++ {
		if trees := StageAt(index).TreeDistances; len(trees) != 0 {
			t.Fatalf("stage %d unexpectedly schedules trees at %v", index, trees)
		}
	}
}

func TestTreeArtworkIsAnimatedAndColored(t *testing.T) {
	tree := newTree(0, 200)
	frame, _ := tree.sprite.ActiveFrame()
	if frame.Width != 7 || frame.Height != treeHeight {
		t.Fatalf("tree size = %dx%d, want 7x%d", frame.Width, frame.Height, treeHeight)
	}
	if canopy := frame.CellAt(3, 2); canopy.Style != StyleTreeCanopy {
		t.Fatalf("canopy style = %d, want %d", canopy.Style, StyleTreeCanopy)
	}
	if highlight := frame.CellAt(2, 2); highlight.Style != StyleTreeHighlight {
		t.Fatalf("highlight style = %d, want %d", highlight.Style, StyleTreeHighlight)
	}
	if trunk := frame.CellAt(2, 6); trunk.Style != StyleTreeTrunk {
		t.Fatalf("trunk style = %d, want %d", trunk.Style, StyleTreeTrunk)
	}

	tree.sprite.Update(650 * time.Millisecond)
	if got := tree.sprite.FrameName(); got != "right" {
		t.Fatalf("frame after sway = %q, want right", got)
	}
}

func TestTreeCollisionPopsBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	balloonX, balloonY := m.balloon.Position()
	m.trees[0].sprite.SetPosition(balloonX, balloonY)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("tree contact did not pop balloon")
	}
}

func TestTreesRemainAnchoredOnResize(t *testing.T) {
	m := NewStageModel(80, 24, 0)
	m.Resize(100, 32)
	balloonX, _ := m.balloon.Position()
	for _, tree := range m.trees {
		x, y := tree.sprite.Position()
		if x != balloonX+tree.distance-m.distance {
			t.Fatalf("tree at distance %v has x %v", tree.distance, x)
		}
		if y != 32-treeHeight-1 {
			t.Fatalf("tree at distance %v has y %v", tree.distance, y)
		}
	}
}
