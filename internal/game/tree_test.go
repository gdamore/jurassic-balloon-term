// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestAuthoredStagesHaveTreeLines(t *testing.T) {
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
	redEyeRun := StageAt(1)
	if len(redEyeRun.TreeDistances) != 5 {
		t.Fatalf("Red-Eye Run tree count = %d, want 5", len(redEyeRun.TreeDistances))
	}
	for index, distance := range redEyeRun.TreeDistances {
		if distance <= 0 || distance >= redEyeRun.GoalDistance {
			t.Fatalf("Red-Eye Run tree %d distance = %v outside stage", index, distance)
		}
	}
	thunderhead := StageAt(2)
	if len(thunderhead.TreeDistances) != 6 {
		t.Fatalf("Thunderhead Pass tree count = %d, want 6", len(thunderhead.TreeDistances))
	}
	longNeckGrove := StageAt(3)
	if len(longNeckGrove.TreeDistances) != 7 {
		t.Fatalf("Long-Neck Grove tree count = %d, want 7", len(longNeckGrove.TreeDistances))
	}
	for index := 4; index < StageCount(); index++ {
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

func TestGroveCyclesBroadleafPalmAndCycad(t *testing.T) {
	want := []treeKind{treeBroadleaf, treeBroadleaf, treePalm, treeBroadleaf, treeCycad}
	for index, kind := range want {
		tree := newTree(index, float64(index*40))
		if tree.kind != kind {
			t.Errorf("tree %d kind = %v, want %v", index, tree.kind, kind)
		}
		frame, _ := tree.sprite.ActiveFrame()
		if frame.Height != treeHeight {
			t.Errorf("tree %d height = %d, want %d", index, frame.Height, treeHeight)
		}
	}
}

func TestPalmSwaysFrondsButKeepsLargeCrownAndFixedTrunk(t *testing.T) {
	palm := newTree(2, 200)
	left, _ := palm.sprite.ActiveFrame()
	if left.Width != 9 {
		t.Fatalf("palm width = %d, want 9", left.Width)
	}
	palm.sprite.Update(650 * time.Millisecond)
	right, _ := palm.sprite.ActiveFrame()
	for y := 2; y <= 8; y++ {
		if left.CellAt(4, y) != right.CellAt(4, y) || left.CellAt(4, y).Rune != '█' {
			t.Fatalf("palm trunk changed at row %d: left=%+v right=%+v", y, left.CellAt(4, y), right.CellAt(4, y))
		}
	}
}

func TestTreeCollisionPopsBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	makeAirborne(m)
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
