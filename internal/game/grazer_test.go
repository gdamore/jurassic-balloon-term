// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestFirstTwoStagesHaveGrazerPairs(t *testing.T) {
	for stageIndex := range 2 {
		m := NewStageModel(80, 24, stageIndex)
		if len(m.grazers) != 6 {
			t.Fatalf("stage %d grazer count = %d, want 6", stageIndex, len(m.grazers))
		}
		for index, grazer := range m.grazers {
			if !grazer.sprite.Collidable() {
				t.Fatalf("stage %d grazer %d is not collidable", stageIndex, index)
			}
		}
		for index := 0; index < len(m.grazers); index += 2 {
			gap := m.grazers[index+1].distance - m.grazers[index].distance
			if gap <= 0 || gap > 13 {
				t.Fatalf("stage %d pair %d gap = %v, want nearby ordered grazers", stageIndex, index/2, gap)
			}
		}
	}
}

func TestThunderheadPassHasOneGrazerPairAroundMidstageLull(t *testing.T) {
	m := NewStageModel(80, 24, 2)
	if len(m.grazers) != 2 {
		t.Fatalf("Thunderhead Pass grazer count = %d, want one pair", len(m.grazers))
	}
	if first, second := m.grazers[0].distance, m.grazers[1].distance; first < 300 || first > 310 || second < 315 || second > 325 {
		t.Fatalf("Thunderhead Pass grazer distances = %v, %v; want pair around 318", first, second)
	}
}

func TestGrazerArtworkIsCompactAndShaded(t *testing.T) {
	grazer := newGrazer(0, 20)
	frame, _ := grazer.sprite.ActiveFrame()
	if frame.Width != 9 || frame.Height != grazerHeight {
		t.Fatalf("grazer size = %dx%d, want 9x%d", frame.Width, frame.Height, grazerHeight)
	}
	if horn := frame.CellAt(0, 0); horn.Rune != '▲' || horn.Style != StyleGrazerHorn {
		t.Fatalf("upper horn = %+v, want filled up-facing triangle", horn)
	}
	if hornBase := frame.CellAt(1, 0); hornBase.Rune != '▲' || hornBase.Style != StyleGrazerHorn {
		t.Fatalf("upper horn base = %+v, want filled up-facing triangle", hornBase)
	}
	if lowerHorn := frame.CellAt(0, 1); lowerHorn.Rune != '◀' || lowerHorn.Style != StyleGrazerHorn {
		t.Fatalf("lower horn = %+v, want filled left-facing triangle", lowerHorn)
	}
	if face := frame.CellAt(1, 1); face.Rune != '█' || face.Style != StyleGrazerHead {
		t.Fatalf("face cell = %+v, want solid head-colored block", face)
	}
	if head := frame.CellAt(2, 1); head.Style != StyleGrazerHead {
		t.Fatalf("head style = %d, want %d", head.Style, StyleGrazerHead)
	}
	if beak := frame.CellAt(1, 2); beak.Rune != '▼' || beak.Style != StyleGrazerBeak {
		t.Fatalf("beak = %+v, want downward beak styling", beak)
	}
	if frill := frame.CellAt(3, 1); frill.Style != StyleGrazerHighlight {
		t.Fatalf("frill style = %d, want %d", frill.Style, StyleGrazerHighlight)
	}
	if body := frame.CellAt(4, 1); body.Style != StyleGrazerBody {
		t.Fatalf("body style = %d, want %d", body.Style, StyleGrazerBody)
	}
	if belly := frame.CellAt(4, 2); belly.Style != StyleGrazerShadow {
		t.Fatalf("belly style = %d, want %d", belly.Style, StyleGrazerShadow)
	}
	if leg := frame.CellAt(3, 3); leg.Style != StyleGrazerLeg {
		t.Fatalf("leg style = %d, want %d", leg.Style, StyleGrazerLeg)
	}
}

func TestGrazerIsAnchoredToForegroundSurface(t *testing.T) {
	m := NewStageModel(80, 24, 0)
	grazer := m.grazers[0]
	x, y := grazer.sprite.Position()
	frame, _ := grazer.sprite.ActiveFrame()
	surfaceY, found := 0, false
	for _, footX := range grazerFootColumns {
		candidate, candidateFound := m.terrain[len(m.terrain)-1].surfaceYAt(int(x) + footX)
		if candidateFound && (!found || candidate > surfaceY) {
			surfaceY, found = candidate, true
		}
	}
	if !found {
		t.Fatal("no foreground beneath first grazer")
	}
	wantY := max(surfaceY-frame.Height+1, 0)
	if y != float64(wantY) {
		t.Fatalf("grazer y = %v, want surface-embedded %d", y, wantY)
	}
}

func TestGrazerKeepsStableAnchorWhileTerrainScrolls(t *testing.T) {
	m := NewStageModel(120, 40, 3)
	grazer := &m.grazers[1]
	grazer.sprite.SetPosition(40.25, 30)
	grazer.anchored = false
	m.anchorGrazers()
	if !grazer.anchored {
		t.Fatal("visible grazer did not acquire a terrain anchor")
	}
	_, anchoredY := grazer.sprite.Position()

	// Once attached to its world-space patch of terrain, rasterized foot
	// columns must not make the animal hop between adjacent terrain cells.
	grazer.sprite.SetPosition(41.25, anchoredY)
	m.anchorGrazers()
	_, gotY := grazer.sprite.Position()
	if gotY != anchoredY {
		t.Fatalf("grazer anchor moved from %v to %v while scrolling", anchoredY, gotY)
	}
}

func TestGrazerCollisionPopsBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	makeAirborne(m)
	x, y := m.grazers[0].sprite.Position()
	m.balloon.SetPosition(x, y)
	m.vy = 0

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("grazer contact did not pop balloon")
	}
	if events := m.Events(); len(events) != 1 || events[0] != EventBalloonPopped {
		t.Fatalf("collision events = %v, want balloon pop", events)
	}
}
