// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

func TestSauropodFramesRaiseHeadAndRetainBody(t *testing.T) {
	longNeck := newSauropod(0, 100)
	if err := longNeck.sprite.SetFrame("low"); err != nil {
		t.Fatal("low sauropod frame is missing")
	}
	low, _ := longNeck.sprite.ActiveFrame()
	if err := longNeck.sprite.SetFrame("high"); err != nil {
		t.Fatal("high sauropod frame is missing")
	}
	high, _ := longNeck.sprite.ActiveFrame()
	if low.Height != sauropodHeight || high.Height != sauropodHeight || low.Width != high.Width {
		t.Fatalf("sauropod frame sizes low=%dx%d high=%dx%d", low.Width, low.Height, high.Width, high.Height)
	}
	if low.CellAt(3, 0).Rune != 0 {
		t.Fatal("low sauropod unexpectedly occupies the upper flight lane")
	}
	if high.CellAt(3, 0).Rune == 0 {
		t.Fatal("raised sauropod head is absent from the upper frame")
	}
	for y := sauropodBodyY; y < sauropodHeight; y++ {
		for x := 0; x < low.Width; x++ {
			if low.CellAt(x, y).Rune != high.CellAt(x, y).Rune {
				t.Fatalf("body changed between low and high poses at %d,%d", x, y)
			}
		}
	}
}

func TestSauropodAnimationRaisesThenLowers(t *testing.T) {
	longNeck := newSauropod(0, 100)
	if err := longNeck.sprite.SetFrame("low"); err != nil {
		t.Fatal(err)
	}
	low, _ := longNeck.sprite.ActiveFrame()
	longNeck.sprite.Update(low.Duration)
	if got := longNeck.sprite.FrameName(); got != "rising" {
		t.Fatalf("frame after low hold = %q, want rising", got)
	}
	longNeck.sprite.Update(320 * time.Millisecond)
	if got := longNeck.sprite.FrameName(); got != "high" {
		t.Fatalf("frame after rise = %q, want high", got)
	}
}

func TestSauropodHerdHasFiveBodyVariants(t *testing.T) {
	var frames [5]sprite.Frame
	silhouettes := make(map[string]bool, len(frames))
	for index := range frames {
		longNeck := newSauropod(index, 100)
		frames[index], _ = longNeck.sprite.ActiveFrame()
		if got, want := frames[index].CellAt(5, 12).Style, sauropodPalettes[index%len(sauropodPalettes)].body; got != want {
			t.Errorf("variant %d body style = %d, want %d", index, got, want)
		}
		silhouettes[sauropodBodySignature(frames[index])] = true
	}
	if len(silhouettes) != len(frames) {
		t.Fatalf("unique body silhouettes = %d, want %d", len(silhouettes), len(frames))
	}
}

func TestSlowSauropodMotionVariants(t *testing.T) {
	swivel := newSauropod(3, 100)
	if swivel.motion != sauropodSwivel {
		t.Fatalf("variant 3 motion = %v, want swivel", swivel.motion)
	}
	if err := swivel.sprite.SetFrame("left"); err != nil {
		t.Fatal(err)
	}
	left, _ := swivel.sprite.ActiveFrame()
	leftX, leftY := frameRunePosition(left, '●')
	swivel.sprite.Update(left.Duration)
	turning, _ := swivel.sprite.ActiveFrame()
	if got := swivel.sprite.FrameName(); got != "turn-right" || turning.Duration < 800*time.Millisecond {
		t.Fatalf("slow swivel transition = %q for %v", got, turning.Duration)
	}
	swivel.sprite.Update(turning.Duration)
	right, _ := swivel.sprite.ActiveFrame()
	rightX, rightY := frameRunePosition(right, '●')
	if rightX <= leftX+3 || rightY != leftY {
		t.Fatalf("swivel eye moved from %d,%d to %d,%d, want a mostly horizontal turn", leftX, leftY, rightX, rightY)
	}
	assertSameSauropodBody(t, left, right)

	diagonal := newSauropod(4, 100)
	if diagonal.motion != sauropodDiagonal {
		t.Fatalf("variant 4 motion = %v, want diagonal", diagonal.motion)
	}
	if err := diagonal.sprite.SetFrame("browse"); err != nil {
		t.Fatal(err)
	}
	browse, _ := diagonal.sprite.ActiveFrame()
	browseX, browseY := frameRunePosition(browse, '●')
	diagonal.sprite.Update(browse.Duration)
	lifting, _ := diagonal.sprite.ActiveFrame()
	if lifting.Duration < 650*time.Millisecond {
		t.Fatalf("diagonal lift duration = %v, want slower than vertical lift", lifting.Duration)
	}
	diagonal.sprite.Update(lifting.Duration)
	reach, _ := diagonal.sprite.ActiveFrame()
	reachX, reachY := frameRunePosition(reach, '●')
	if reachX <= browseX+4 || reachY >= browseY-3 {
		t.Fatalf("diagonal eye moved from %d,%d to %d,%d", browseX, browseY, reachX, reachY)
	}
	assertSameSauropodBody(t, browse, reach)
}

func sauropodBodySignature(frame sprite.Frame) string {
	var signature strings.Builder
	for y := sauropodBodyY; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			signature.WriteRune(frame.CellAt(x, y).Rune)
		}
	}
	return signature.String()
}

func frameRunePosition(frame sprite.Frame, target rune) (int, int) {
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			if frame.CellAt(x, y).Rune == target {
				return x, y
			}
		}
	}
	return -1, -1
}

func assertSameSauropodBody(t *testing.T, first, second sprite.Frame) {
	t.Helper()
	if sauropodBodySignature(first) != sauropodBodySignature(second) {
		t.Fatal("sauropod body moves with its neck animation")
	}
}

func TestLongNeckGrovePopulatesAndAnchorsSauropods(t *testing.T) {
	m := NewStageModel(120, 40, 3)
	if len(m.sauropods) != 11 {
		t.Fatalf("sauropods = %d, want 11", len(m.sauropods))
	}
	if m.stage.Pterosaurs {
		t.Fatal("Long-Neck Grove still spawns routine pterosaurs")
	}
	if span := m.sauropods[2].distance - m.sauropods[0].distance; span >= 50 {
		t.Fatalf("opening herd spans %v cells, want a compact group", span)
	}
	if gap := m.sauropods[3].distance - m.sauropods[2].distance; gap <= 80 {
		t.Fatalf("opening herd and next encounter are only %v cells apart", gap)
	}
	longNeck := m.sauropods[0]
	frame, _ := longNeck.sprite.ActiveFrame()
	x, y := longNeck.sprite.Position()
	ground := m.terrain[len(m.terrain)-1]
	surface, found := 0, false
	for _, footX := range sauropodFootColumns {
		candidate, candidateFound := ground.surfaceYAt(int(math.Floor(x)) + footX)
		if candidateFound && (!found || candidate > surface) {
			surface, found = candidate, true
		}
	}
	if !found {
		t.Fatal("no terrain under first sauropod")
	}
	if bottom := int(y) + frame.Height - 1; bottom != surface {
		t.Fatalf("sauropod bottom = %d, terrain surface = %d", bottom, surface)
	}
}

func TestSauropodContactPopsBalloon(t *testing.T) {
	m := NewStageModel(120, 40, 3)
	makeAirborne(m)
	x, y := m.balloon.Position()
	if err := m.sauropods[0].sprite.SetFrame("high"); err != nil {
		t.Fatal(err)
	}
	m.sauropods[0].sprite.SetPosition(x, y)
	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("sauropod contact did not pop balloon")
	}
}

func TestSauropodCallOccursOnceOnEntry(t *testing.T) {
	m := NewStageModel(120, 40, 3)
	makeAirborne(m)
	m.sauropods[0].sprite.SetPosition(float64(m.Width-1), 10)
	m.Update(time.Millisecond)
	if got := countEvent(m.Events(), EventSauropodCalled); got != 1 {
		t.Fatalf("entry calls = %d, want 1", got)
	}
	m.Update(time.Millisecond)
	if got := countEvent(m.Events(), EventSauropodCalled); got != 0 {
		t.Fatalf("repeated calls = %d, want 0", got)
	}
	m.sauropods[1].sprite.SetPosition(float64(m.Width-1), 10)
	m.Update(time.Millisecond)
	if got := countEvent(m.Events(), EventSauropodCalled); got != 0 {
		t.Fatalf("second animal in herd calls = %d, want 0", got)
	}
}

func TestOutgoingSauropodRemainsHazardous(t *testing.T) {
	previous := NewStageModel(120, 40, 3)
	makeAirborne(previous)
	x, y := previous.balloon.Position()
	hazard := previous.sauropods[0].sprite
	if err := hazard.SetFrame("high"); err != nil {
		t.Fatal(err)
	}
	hazard.SetPosition(x, y)
	previous.completeStage()

	next := previous.ContinueToStage(4)
	if !modelHasSprite(next, hazard) {
		t.Fatal("visible outgoing sauropod was cleared at the stage boundary")
	}
	next.Update(time.Millisecond)
	if !next.Popped() {
		t.Fatal("outgoing sauropod became harmless at the stage boundary")
	}
}

func countEvent(events []Event, target Event) int {
	count := 0
	for _, event := range events {
		if event == target {
			count++
		}
	}
	return count
}
