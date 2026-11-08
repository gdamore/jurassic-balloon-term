// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"testing"
	"time"
)

func TestPterosaurFlapsWithDynamicMask(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 6, 0.5)
	glide, _ := pterosaur.sprite.ActiveFrame()
	if glide.Width != 15 || glide.Height != pterosaurHeight {
		t.Fatalf("pterosaur size = %dx%d, want 15x%d", glide.Width, glide.Height, pterosaurHeight)
	}
	if glide.CellAt(10, 1).Rune == 0 {
		t.Fatal("gliding wing cell is transparent")
	}
	if eye := glide.CellAt(2, 1); eye.Rune != '●' || eye.Style != StylePterosaurEye {
		t.Fatalf("eye cell = %+v, want body-embedded eye", eye)
	}
	if beak := glide.CellAt(0, 1); beak.Style != StylePterosaurBeak {
		t.Fatalf("beak style = %d, want %d", beak.Style, StylePterosaurBeak)
	}
	if beakBase := glide.CellAt(1, 1); beakBase.Rune != '▰' || beakBase.Style != StylePterosaurBeak {
		t.Fatalf("beak base = %+v, want thick beak styling", beakBase)
	}
	if crest := glide.CellAt(4, 0); crest.Style != StylePterosaurCrest {
		t.Fatalf("crest style = %d, want %d", crest.Style, StylePterosaurCrest)
	}
	if wing := glide.CellAt(9, 1); wing.Style != StylePterosaurWing {
		t.Fatalf("near wing style = %d, want %d", wing.Style, StylePterosaurWing)
	}
	if tip := glide.CellAt(10, 1); tip.Style != StylePterosaurWingHighlight {
		t.Fatalf("wing-tip style = %d, want %d", tip.Style, StylePterosaurWingHighlight)
	}
	if torso := glide.CellAt(5, 1); torso.Rune != '▀' || torso.Style != StylePterosaur {
		t.Fatalf("torso = %+v, want half-height body", torso)
	}

	pterosaur.sprite.Update(pterosaurGlideTime)
	if got := pterosaur.sprite.FrameName(); got != "up" {
		t.Fatalf("frame after glide = %q, want up", got)
	}
	up, _ := pterosaur.sprite.ActiveFrame()
	if up.CellAt(11, 1).Rune != 0 {
		t.Fatal("gliding wing remains in raised-wing collision mask")
	}
	pterosaur.sprite.Update(pterosaurWingUpTime + pterosaurWingDownTime)
	if got := pterosaur.sprite.FrameName(); got != "glide" {
		t.Fatalf("frame after quick flap = %q, want glide", got)
	}
}

func TestPterosaurPathIsBounded(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 8, 0.5)
	for range 100 {
		pterosaur.updatePath(33 * time.Millisecond)
		_, y := pterosaur.sprite.Position()
		if math.Abs(y-pterosaur.baseY) > pterosaurBobHeight+1e-9 {
			t.Fatalf("pterosaur y = %v outside bob range around %v", y, pterosaur.baseY)
		}
	}
}

func TestPterosaurFlapProducesDecayingSpeedBoost(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 8, 0.5)

	pterosaur.sprite.Update(pterosaurGlideTime)
	pterosaur.updatePath(0)
	pterosaur.sprite.Update(pterosaurWingUpTime)
	pterosaur.updatePath(0)
	pterosaur.sprite.Update(pterosaurWingDownTime)
	pterosaur.updatePath(0)

	boosted, _ := pterosaur.sprite.Velocity()
	if boosted != pterosaurBoostSpeed {
		t.Fatalf("velocity after flap = %v, want %v", boosted, pterosaurBoostSpeed)
	}
	_, lungeY := pterosaur.sprite.Position()
	lungePhase := pterosaur.phase

	pterosaur.updatePath(pterosaurBoostTime / 2)
	decaying, _ := pterosaur.sprite.Velocity()
	if decaying <= boosted || decaying >= pterosaurSpeed {
		t.Fatalf("decaying velocity = %v, want between %v and %v", decaying, boosted, pterosaurSpeed)
	}
	_, y := pterosaur.sprite.Position()
	if y != lungeY || pterosaur.phase != lungePhase {
		t.Fatalf("lunge changed vertical path: y = %v, phase = %v", y, pterosaur.phase)
	}

	pterosaur.updatePath(pterosaurBoostTime / 2)
	cruising, _ := pterosaur.sprite.Velocity()
	if cruising != pterosaurSpeed {
		t.Fatalf("velocity after boost = %v, want %v", cruising, pterosaurSpeed)
	}
	_, y = pterosaur.sprite.Position()
	if y != lungeY || pterosaur.phase != lungePhase {
		t.Fatalf("last lunge step changed vertical path: y = %v, phase = %v", y, pterosaur.phase)
	}

	pterosaur.updatePath(time.Millisecond)
	if pterosaur.phase == lungePhase {
		t.Fatal("vertical path did not resume after lunge")
	}
}

func TestPterosaurCollisionPopsBalloon(t *testing.T) {
	m := NewModel(80, 30)
	balloonX, balloonY := m.balloon.Position()
	m.pterosaur.baseY = balloonY
	m.pterosaur.phase = 0
	m.pterosaur.sprite.SetPosition(balloonX+0.5, balloonY)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("pterosaur contact did not pop balloon")
	}
	if events := m.Events(); len(events) != 1 || events[0] != EventBalloonPopped {
		t.Fatalf("collision events = %v, want balloon pop", events)
	}
}

func TestPterosaurLungeProducesGameEvent(t *testing.T) {
	m := NewModel(80, 30)
	m.pterosaur.lastFrame = "down"
	m.pterosaur.sprite.SetPosition(40, 6)
	if err := m.pterosaur.sprite.SetFrame("glide"); err != nil {
		t.Fatal(err)
	}

	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 1 || events[0] != EventPterosaurLunged {
		t.Fatalf("lunge events = %v, want pterosaur lunge", events)
	}

	m.pterosaur.lastFrame = "down"
	m.pterosaur.sprite.SetPosition(float64(m.Width+1), 6)
	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 0 {
		t.Fatalf("offscreen lunge events = %v, want none", events)
	}
}

func TestDepartedPterosaurIsReplaced(t *testing.T) {
	m := NewModel(80, 24)
	departed := m.pterosaur.sprite
	frame, _ := departed.ActiveFrame()
	departed.SetPosition(float64(-frame.Width-1), 4)

	m.Update(time.Millisecond)
	if m.pterosaur.sprite == departed {
		t.Fatal("departed pterosaur was repositioned instead of replaced")
	}
	x, _ := m.pterosaur.sprite.Position()
	if x <= float64(m.Width) {
		t.Fatalf("replacement pterosaur x = %v, want beyond right edge %d", x, m.Width)
	}
	for _, item := range m.Sprites() {
		if item == departed {
			t.Fatal("departed pterosaur remains in scene")
		}
	}
}
